package highlight

import (
	"strings"
	"testing"
)

func TestLevel_Disabled(t *testing.T) {
	h := New(false)
	got := h.Level("error")
	if got != "error" {
		t.Errorf("expected plain level, got %q", got)
	}
}

func TestLevel_Enabled_KnownLevel(t *testing.T) {
	h := New(true)
	for level, color := range LevelColors {
		got := h.Level(level)
		if !strings.Contains(got, color) {
			t.Errorf("level %q: expected color %q in output %q", level, color, got)
		}
		if !strings.Contains(got, Reset) {
			t.Errorf("level %q: expected reset code in output %q", level, got)
		}
	}
}

func TestLevel_Enabled_UnknownLevel(t *testing.T) {
	h := New(true)
	got := h.Level("trace")
	if got != "trace" {
		t.Errorf("unknown level should be returned unchanged, got %q", got)
	}
}

func TestField_Disabled(t *testing.T) {
	h := New(false)
	got := h.Field("key", "value")
	if got != "key=value" {
		t.Errorf("expected plain key=value, got %q", got)
	}
}

func TestField_Enabled(t *testing.T) {
	h := New(true)
	got := h.Field("key", "value")
	if !strings.Contains(got, "key") || !strings.Contains(got, "value") {
		t.Errorf("field output missing key or value: %q", got)
	}
	if !strings.Contains(got, Bold) {
		t.Errorf("expected bold code in field output: %q", got)
	}
}

func TestSource_Disabled(t *testing.T) {
	h := New(false)
	got := h.Source("app.log")
	if got != "app.log" {
		t.Errorf("expected plain label, got %q", got)
	}
}

func TestSource_Enabled(t *testing.T) {
	h := New(true)
	got := h.Source("app.log")
	if !strings.Contains(got, Blue) {
		t.Errorf("expected blue color in source output: %q", got)
	}
}

func TestMessage_Disabled(t *testing.T) {
	h := New(false)
	got := h.Message("hello world")
	if got != "hello world" {
		t.Errorf("expected plain message, got %q", got)
	}
}

func TestMessage_Enabled(t *testing.T) {
	h := New(true)
	got := h.Message("hello world")
	if !strings.Contains(got, Bold) {
		t.Errorf("expected bold in message output: %q", got)
	}
}
