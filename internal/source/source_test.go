package source_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/source"
)

func TestTail_ReadsAllLines(t *testing.T) {
	input := "{\"level\":\"info\",\"msg\":\"start\"}\n{\"level\":\"error\",\"msg\":\"fail\"}\n"
	r := source.New("test", strings.NewReader(input))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var lines []source.Line
	for line := range r.Tail(ctx) {
		lines = append(lines, line)
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0].Source != "test" {
		t.Errorf("expected source 'test', got %q", lines[0].Source)
	}
	if lines[1].Text != "{\"level\":\"error\",\"msg\":\"fail\"}" {
		t.Errorf("unexpected line text: %q", lines[1].Text)
	}
}

func TestTail_CancelContext(t *testing.T) {
	// Infinite reader simulation via pipe
	pr, pw := strings.NewReader(""), strings.NewReader("")
	_ = pr
	_ = pw

	// Use a reader that blocks — simulate with a large repeated input
	big := strings.Repeat("{\"x\":\"y\"}\n", 10000)
	r := source.New("big", strings.NewReader(big))

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	count := 0
	for range r.Tail(ctx) {
		count++
	}
	// Should read very few or zero lines after immediate cancel
	if count > 100 {
		t.Errorf("expected context cancellation to stop reading early, got %d lines", count)
	}
}

func TestNew_Name(t *testing.T) {
	r := source.New("myfile.log", strings.NewReader(""))
	if r.Name() != "myfile.log" {
		t.Errorf("expected name 'myfile.log', got %q", r.Name())
	}
}
