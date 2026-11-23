package output

import (
	"bytes"
	"io"
	"regexp"
	"sync"
)

// maxPendingBytes bounds the buffer of a line without newline, e.g. progress
// output, so it cannot grow without limit.
const maxPendingBytes = 64 * 1024

// githubTokenPattern matches GitHub personal, OAuth, user, server and refresh
// tokens as well as fine-grained PATs. The prefix stays readable in the log.
var githubTokenPattern = regexp.MustCompile(`\b(gh[pousr]_|github_pat_)[A-Za-z0-9_]{20,}`)

func redactSecrets(p []byte) []byte {
	return githubTokenPattern.ReplaceAll(p, []byte("${1}<REDACTED>"))
}

// redactingWriter redacts secrets before writing to dst. It writes complete
// lines only, so a token split across two writes is still caught.
type redactingWriter struct {
	mu      sync.Mutex
	dst     io.Writer
	pending []byte
}

func (w *redactingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.pending = append(w.pending, p...)

	end := bytes.LastIndexByte(w.pending, '\n') + 1
	if len(w.pending) > maxPendingBytes {
		end = forcedFlushEnd(w.pending)
	}

	if err := w.writeRedacted(end); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Flush writes a pending line without trailing newline.
func (w *redactingWriter) Flush() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.writeRedacted(len(w.pending))
}

// forcedFlushEnd cuts an oversized line before its last run of word
// characters. Tokens consist of word characters only and start at a word
// boundary, so the cut never splits one. A buffer that is a single word run is
// flushed whole: a token in it can only start at its beginning, and the greedy
// pattern then redacts the run.
func forcedFlushEnd(p []byte) int {
	if runStart := bytes.LastIndexFunc(p, isNotWordChar) + 1; runStart > 0 {
		return runStart
	}
	return len(p)
}

func isNotWordChar(r rune) bool {
	return r != '_' && !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') && !('0' <= r && r <= '9')
}

func (w *redactingWriter) writeRedacted(end int) error {
	if end == 0 {
		return nil
	}
	_, err := w.dst.Write(redactSecrets(w.pending[:end]))
	w.pending = append(w.pending[:0], w.pending[end:]...)
	return err
}
