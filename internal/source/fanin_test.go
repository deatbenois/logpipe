package source_test

import (
	"context"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/source"
)

func TestFanIn_MergesMultipleSources(t *testing.T) {
	a := source.New("a", strings.NewReader("line-a1\nline-a2\n"))
	b := source.New("b", strings.NewReader("line-b1\n"))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var texts []string
	for line := range source.TailAll(ctx, []*source.Reader{a, b}) {
		texts = append(texts, line.Text)
	}

	if len(texts) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(texts), texts)
	}

	sort.Strings(texts)
	expected := []string{"line-a1", "line-a2", "line-b1"}
	for i, want := range expected {
		if texts[i] != want {
			t.Errorf("line[%d]: expected %q, got %q", i, want, texts[i])
		}
	}
}

func TestFanIn_EmptySources(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var lines []source.Line
	for line := range source.TailAll(ctx, []*source.Reader{}) {
		lines = append(lines, line)
	}
	if len(lines) != 0 {
		t.Errorf("expected no lines from empty sources, got %d", len(lines))
	}
}

func TestFanIn_SourceLabelsPreserved(t *testing.T) {
	a := source.New("app.log", strings.NewReader("{\"msg\":\"hello\"}\n"))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for line := range source.TailAll(ctx, []*source.Reader{a}) {
		if line.Source != "app.log" {
			t.Errorf("expected source 'app.log', got %q", line.Source)
		}
	}
}
