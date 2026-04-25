package output

import (
	"io"
	"os"
	"sync"
)

// Writer wraps an io.Writer with thread-safe writes and optional line counting.
type Writer struct {
	mu      sync.Mutex
	w       io.Writer
	lines   int64
	closed  bool
}

// New returns a Writer backed by the given io.Writer.
func New(w io.Writer) *Writer {
	return &Writer{w: w}
}

// NewStdout returns a Writer backed by os.Stdout.
func NewStdout() *Writer {
	return New(os.Stdout)
}

// NewStderr returns a Writer backed by os.Stderr.
func NewStderr() *Writer {
	return New(os.Stderr)
}

// Write writes p to the underlying writer in a thread-safe manner.
func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return 0, io.ErrClosedPipe
	}
	n, err := w.w.Write(p)
	if err == nil {
		w.lines++
	}
	return n, err
}

// WriteLine writes p followed by a newline character.
func (w *Writer) WriteLine(p []byte) (int, error) {
	return w.Write(append(p, '\n'))
}

// LinesWritten returns the number of successful Write calls made.
func (w *Writer) LinesWritten() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lines
}

// Close marks the writer as closed; subsequent writes will return an error.
func (w *Writer) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.closed = true
}
