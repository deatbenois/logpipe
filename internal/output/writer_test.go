package output

import (
	"bytes"
	"io"
	"sync"
	"testing"
)

func TestWrite_Basic(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	_, err := w.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "hello" {
		t.Errorf("expected 'hello', got %q", buf.String())
	}
}

func TestWriteLine_AppendsNewline(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	w.WriteLine([]byte("line"))
	if buf.String() != "line\n" {
		t.Errorf("expected 'line\\n', got %q", buf.String())
	}
}

func TestLinesWritten_Counts(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	w.Write([]byte("a"))
	w.Write([]byte("b"))
	w.Write([]byte("c"))
	if w.LinesWritten() != 3 {
		t.Errorf("expected 3 lines written, got %d", w.LinesWritten())
	}
}

func TestClose_RejectsWrites(t *testing.T) {
	var buf bytes.Buffer
	w := New(&buf)
	w.Close()
	_, err := w.Write([]byte("x"))
	if err != io.ErrClosedPipe {
		t.Errorf("expected ErrClosedPipe, got %v", err)
	}
}

func TestWrite_Concurrent(t *testing.T) {
	var buf safeBuffer
	w := New(&buf)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.Write([]byte("x"))
		}()
	}
	wg.Wait()
	if w.LinesWritten() != 50 {
		t.Errorf("expected 50 writes, got %d", w.LinesWritten())
	}
}

// safeBuffer is a thread-safe bytes.Buffer for testing.
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *safeBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}
