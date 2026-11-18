package acloud

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/modell-aachen/machine/internal/output"
	"github.com/modell-aachen/machine/internal/platform"
	"github.com/modell-aachen/machine/internal/provision/completions"
)

const (
	binaryName        = "acloud"
	defaultBaseURL    = "https://nexus.modac.cloud/repository/acloud/latest/"
	nexusUser         = "bot-ro"
	tokenEnv          = "NEXUS_BOT_TOKEN"
	downloadTimeout   = 5 * time.Minute
	maxVersionBytes   = 64
	maxChecksumsBytes = 64 << 10
	maxArchiveBytes   = 256 << 20
	completionCommand = "completion"
)

var errNotPublished = errors.New("not published")

type installer struct {
	client           *http.Client
	baseURL          string
	token            string
	binDir           string
	installedVersion func(binPath string) string
}

func Run(out *output.Context, plat platform.Platform) error {
	_ = plat

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to determine home directory: %w", err)
	}

	inst := installer{
		client:           &http.Client{Timeout: downloadTimeout},
		baseURL:          defaultBaseURL,
		token:            os.Getenv(tokenEnv),
		binDir:           filepath.Join(home, ".local", "bin"),
		installedVersion: probeVersion,
	}
	version, err := inst.install(out)
	if err != nil || version == "" {
		return err
	}

	installCompletions(out, home, version)
	return nil
}

func (i installer) install(out *output.Context) (string, error) {
	latest, err := i.fetch("VERSION", maxVersionBytes)
	if errors.Is(err, errNotPublished) {
		out.Skipped("no acloud release published in Nexus yet; skipping install")
		return "", nil
	}
	if err != nil {
		return "", err
	}
	version := strings.TrimSpace(string(latest))
	if version == "" {
		return "", fmt.Errorf("%sVERSION is empty", i.baseURL)
	}

	binPath := filepath.Join(i.binDir, binaryName)
	if i.installedVersion(binPath) == version {
		out.Skipped(fmt.Sprintf("acloud %s already installed", version))
		return version, nil
	}

	out.Step(fmt.Sprintf("Installing acloud %s", version))
	archive, err := i.downloadVerified()
	if err != nil {
		return "", err
	}
	if err := installBinary(archive, binPath); err != nil {
		return "", err
	}
	return version, nil
}

func (i installer) downloadVerified() ([]byte, error) {
	checksums, err := i.fetch("checksums.txt", maxChecksumsBytes)
	if err != nil {
		return nil, err
	}
	expected, err := checksumFor(checksums, archiveName())
	if err != nil {
		return nil, err
	}
	archive, err := i.fetch(archiveName(), maxArchiveBytes)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(archive)
	if actual := hex.EncodeToString(sum[:]); actual != expected {
		return nil, fmt.Errorf("checksum mismatch for %s: expected %s, got %s", archiveName(), expected, actual)
	}
	return archive, nil
}

func (i installer) fetch(name string, limit int64) ([]byte, error) {
	url := i.baseURL + name
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request for %s: %w", url, err)
	}
	if i.token != "" {
		req.SetBasicAuth(nexusUser, i.token)
	}

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("%s: %w", url, errNotPublished)
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, fmt.Errorf("access to %s denied (%s); check that %s is set and valid", url, resp.Status, tokenEnv)
	default:
		return nil, fmt.Errorf("failed to download %s: %s", url, resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", url, err)
	}
	if int64(len(body)) > limit {
		return nil, fmt.Errorf("%s exceeds %d bytes", url, limit)
	}
	return body, nil
}

func archiveName() string {
	return fmt.Sprintf("%s_%s_%s.tar.gz", binaryName, runtime.GOOS, runtime.GOARCH)
}

func checksumFor(checksums []byte, name string) (string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(checksums))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[1] == name {
			return strings.ToLower(fields[0]), nil
		}
	}
	return "", fmt.Errorf("no checksum for %s in checksums.txt", name)
}

func installBinary(archive []byte, binPath string) error {
	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}
	tmp, err := os.CreateTemp(filepath.Dir(binPath), "."+binaryName+"-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	if err := extractBinary(archive, tmp); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to write %s: %w", binaryName, err)
	}
	if err := os.Chmod(tmp.Name(), 0o755); err != nil {
		return fmt.Errorf("failed to make %s executable: %w", binaryName, err)
	}
	if err := os.Rename(tmp.Name(), binPath); err != nil {
		return fmt.Errorf("failed to install %s: %w", binPath, err)
	}
	return nil
}

func extractBinary(archive []byte, dst io.Writer) error {
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", archiveName(), err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return fmt.Errorf("%s contains no %s binary", archiveName(), binaryName)
		}
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", archiveName(), err)
		}
		if header.Typeflag != tar.TypeReg || path.Base(header.Name) != binaryName {
			continue
		}
		if _, err := io.Copy(dst, tr); err != nil {
			return fmt.Errorf("failed to extract %s: %w", binaryName, err)
		}
		return nil
	}
}

func probeVersion(binPath string) string {
	out, err := exec.Command(binPath, "--version").Output()
	if err != nil {
		return ""
	}
	return parseVersion(string(out))
}

func parseVersion(output string) string {
	firstLine, _, _ := strings.Cut(output, "\n")
	fields := strings.Fields(firstLine)
	if len(fields) < 3 || fields[0] != binaryName || fields[1] != "version" {
		return ""
	}
	return fields[2]
}

func completionVersion(version string) string {
	return strings.ReplaceAll(version, ".", "_")
}

func installCompletions(out *output.Context, home, version string) {
	for _, shell := range []string{"bash", "zsh"} {
		if err := completions.Install(out, home, shell, binaryName, completionVersion(version), completionCommand); err != nil {
			out.Info(fmt.Sprintf("acloud completion for %s not installed: %v", shell, err))
		}
	}
}
