package output

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"

	symbolSuccess = "✓"
	symbolFailure = "✗"
	symbolSkipped = "⊘"
	symbolRunning = "→"
)

// Context manages logging and output formatting for the provisioning process
type Context struct {
	logFile *os.File
	logPath string
}

// New creates a new output context with a log file
func New() (*Context, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	logDir := filepath.Join(homeDir, ".machine", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	logPath := filepath.Join(logDir, fmt.Sprintf("provision-%s.log", timestamp))

	logFile, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	ctx := &Context{
		logFile: logFile,
		logPath: logPath,
	}

	ctx.writeLog("=== Machine Provisioning Log ===\n")
	ctx.writeLog("Started at: %s\n\n", time.Now().Format(time.RFC3339))

	return ctx, nil
}

// Close closes the log file
func (c *Context) Close() error {
	if c.logFile != nil {
		c.writeLog("\n=== Provisioning Complete ===\n")
		c.writeLog("Finished at: %s\n", time.Now().Format(time.RFC3339))
		return c.logFile.Close()
	}
	return nil
}

// LogPath returns the path to the log file
func (c *Context) LogPath() string {
	return c.logPath
}

// writeLog writes a message to the log file only
func (c *Context) writeLog(format string, args ...interface{}) {
	if c.logFile != nil {
		fmt.Fprintf(c.logFile, format, args...)
	}
}

// StartModule prints the module start message and logs it
func (c *Context) StartModule(name string) {
	msg := fmt.Sprintf("%s %sRunning module:%s %s", symbolRunning, colorBlue, colorReset, name)
	fmt.Println(msg)
	c.writeLog("\n--- Module: %s ---\n", name)
}

// Success prints a success message
func (c *Context) Success(message string) {
	fmt.Printf("  %s%s%s %s\n", colorGreen, symbolSuccess, colorReset, message)
	c.writeLog("[SUCCESS] %s\n", message)
}

// Failure prints a failure message
func (c *Context) Failure(message string) {
	fmt.Printf("  %s%s%s %s\n", colorRed, symbolFailure, colorReset, message)
	c.writeLog("[FAILURE] %s\n", message)
}

// Skipped prints a skipped message
func (c *Context) Skipped(message string) {
	fmt.Printf("  %s%s%s %s\n", colorYellow, symbolSkipped, colorReset, message)
	c.writeLog("[SKIPPED] %s\n", message)
}

// Info prints an info message
func (c *Context) Info(message string) {
	fmt.Printf("  %s%s%s\n", colorGray, message, colorReset)
	c.writeLog("[INFO] %s\n", message)
}

// Step prints a step being executed
func (c *Context) Step(message string) {
	fmt.Printf("  %s→%s %s\n", colorBlue, colorReset, message)
	c.writeLog("[STEP] %s\n", message)
}

// RunCommand executes a command and returns the error if any
// All command output is written to the log file only
func (c *Context) RunCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)

	// Create a multi-writer that writes to both log file and captures output
	logWriter := &logWriter{ctx: c, prefix: ""}

	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	c.writeLog("\n$ %s %s\n", name, strings.Join(arg, " "))

	err := cmd.Run()

	if err != nil {
		c.writeLog("[Command failed with error: %v]\n", err)
	}

	return err
}

// RunCommandWithMessage executes a command with a step message
func (c *Context) RunCommandWithMessage(message, name string, arg ...string) error {
	c.Step(message)
	return c.RunCommand(name, arg...)
}

// CheckAndRun checks if a marker file exists, and if not, runs the provided function
// If successful, it creates the marker file
func (c *Context) CheckAndRun(markerPath string, skipMessage string, fn func() error) error {
	// Check if already done
	if _, err := os.Stat(markerPath); err == nil {
		c.Skipped(skipMessage)
		return nil
	}

	// Run the function
	if err := fn(); err != nil {
		return err
	}

	// Create marker file
	if err := os.MkdirAll(filepath.Dir(markerPath), 0755); err != nil {
		return fmt.Errorf("failed to create marker directory: %w", err)
	}

	if err := os.WriteFile(markerPath, []byte(""), 0644); err != nil {
		return fmt.Errorf("failed to create marker file: %w", err)
	}

	return nil
}

// logWriter implements io.Writer to write command output to the log file
type logWriter struct {
	ctx    *Context
	prefix string
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	if w.ctx.logFile != nil {
		return w.ctx.logFile.Write(p)
	}
	return len(p), nil
}

// MultiWriter creates an io.Writer that writes to both the log and another writer
// If w is nil, returns a writer that only writes to the log file
func (c *Context) MultiWriter(w io.Writer) io.Writer {
	if w == nil {
		return c.logFile
	}
	return io.MultiWriter(c.logFile, w)
}

// PrintSummary prints a final summary with the log file location
func (c *Context) PrintSummary() {
	fmt.Printf("\n%sProvisioning complete!%s\n", colorGreen, colorReset)
	fmt.Printf("%sFull logs available at:%s %s\n", colorGray, colorReset, c.logPath)
}

// PrintError prints an error summary
func (c *Context) PrintError(err error) {
	fmt.Printf("\n%sProvisioning failed:%s %v\n", colorRed, colorReset, err)
	fmt.Printf("%sFull logs available at:%s %s\n", colorGray, colorReset, c.logPath)
}
