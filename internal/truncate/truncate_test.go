package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/truncate"
)

func TestApply_ShortLine_Unchanged(t *testing.T) {
	tr := truncate.New(80)
	line := "short line"
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_ExactLimit_Unchanged(t *testing.T) {
	tr := truncate.New(10)
	line := "1234567890"
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_LongLine_Truncated(t *testing.T) {
	tr := truncate.New(20)
	line := strings.Repeat("a", 50)
	got := tr.Apply(line)
	if len(got) > 20 {
		t.Fatalf("expected len <= 20, got %d", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestApply_Disabled_NoTruncation(t *testing.T) {
	tr := truncate.New(0)
	line := strings.Repeat("x", 10000)
	if got := tr.Apply(line); got != line {
		t.Fatal("expected line unchanged when disabled")
	}
}

func TestApply_CustomSuffix(t *testing.T) {
	tr := truncate.NewWithSuffix(15, " [cut]")
	line := strings.Repeat("b", 50)
	got := tr.Apply(line)
	if !strings.HasSuffix(got, " [cut]") {
		t.Fatalf("expected custom suffix, got %q", got)
	}
	if len(got) > 15 {
		t.Fatalf("expected len <= 15, got %d", len(got))
	}
}

func TestApply_UTF8_NotSplit(t *testing.T) {
	tr := truncate.New(5)
	// 'é' is 2 bytes; "aé" is 3 bytes, "aaé" is 4 bytes
	line := "aaébb"
	got := tr.Apply(line)
	if len(got) > 5 {
		t.Fatalf("expected len <= 5, got %d", len(got))
	}
}

func TestEnabled(t *testing.T) {
	if truncate.New(0).Enabled() {
		t.Fatal("expected disabled for maxLen=0")
	}
	if !truncate.New(100).Enabled() {
		t.Fatal("expected enabled for maxLen=100")
	}
}

func TestMaxLen(t *testing.T) {
	tr := truncate.New(42)
	if tr.MaxLen() != 42 {
		t.Fatalf("expected 42, got %d", tr.MaxLen())
	}
}
