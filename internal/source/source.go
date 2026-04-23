package source

import (
	"bufio"
	"context"
	"io"
	"os"
)

// Line represents a single log line with its origin.
type Line struct {
	Text   string
	Source string
}

// Reader reads lines from an io.Reader and sends them to a channel.
type Reader struct {
	name   string
	reader io.Reader
}

// New creates a new Reader with the given name and underlying io.Reader.
func New(name string, r io.Reader) *Reader {
	return &Reader{name: name, reader: r}
}

// NewFromFile opens a file and returns a Reader for it.
func NewFromFile(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return New(path, f), nil
}

// NewFromStdin returns a Reader reading from os.Stdin.
func NewFromStdin() *Reader {
	return New("stdin", os.Stdin)
}

// Tail reads lines from the source and sends each to the returned channel.
// The channel is closed when the context is cancelled or EOF is reached.
func (r *Reader) Tail(ctx context.Context) <-chan Line {
	ch := make(chan Line)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(r.reader)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- Line{Text: scanner.Text(), Source: r.name}:
			}
		}
	}()
	return ch
}

// Name returns the name/label of this source.
func (r *Reader) Name() string {
	return r.name
}
