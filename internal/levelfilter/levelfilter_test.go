package levelfilter_test

import (
	"testing"

	"github.com/yourorg/logpipe/internal/levelfilter"
)

func TestNew_EmptyLevel_Disabled(t *testing.T) {
	f, err := levelfilter.New("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(`{"level":"debug","msg":"hi"}`) {
		t.Error("disabled filter should allow everything")
	}
}

func TestNew_UnknownLevel_ReturnsError(t *testing.T) {
	_, err := levelfilter.New("verbose", nil)
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestAllow_BelowMin_Rejected(t *testing.T) {
	f, _ := levelfilter.New("warn", nil)
	if f.Allow(`{"level":"debug","msg":"low"}`) {
		t.Error("debug should be rejected when min is warn")
	}
	if f.Allow(`{"level":"info","msg":"mid"}`) {
		t.Error("info should be rejected when min is warn")
	}
}

func TestAllow_AtMin_Passes(t *testing.T) {
	f, _ := levelfilter.New("warn", nil)
	if !f.Allow(`{"level":"warn","msg":"at"}`) {
		t.Error("warn should pass when min is warn")
	}
}

func TestAllow_AboveMin_Passes(t *testing.T) {
	f, _ := levelfilter.New("warn", nil)
	if !f.Allow(`{"level":"error","msg":"high"}`) {
		t.Error("error should pass when min is warn")
	}
	if !f.Allow(`{"level":"fatal","msg":"crit"}`) {
		t.Error("fatal should pass when min is warn")
	}
}

func TestAllow_NonJSON_Passes(t *testing.T) {
	f, _ := levelfilter.New("error", nil)
	if !f.Allow("plain text log line") {
		t.Error("non-JSON lines should always pass")
	}
}

func TestAllow_NoLevelField_Passes(t *testing.T) {
	f, _ := levelfilter.New("error", nil)
	if !f.Allow(`{"msg":"no level here"}`) {
		t.Error("lines without a level field should pass")
	}
}

func TestAllow_CustomFields(t *testing.T) {
	f, _ := levelfilter.New("error", []string{"severity"})
	if f.Allow(`{"severity":"info","msg":"low"}`) {
		t.Error("info severity should be rejected when min is error")
	}
	if !f.Allow(`{"severity":"fatal","msg":"crit"}`) {
		t.Error("fatal severity should pass when min is error")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f, _ := levelfilter.New("INFO", nil)
	if !f.Allow(`{"level":"INFO","msg":"upper"}`) {
		t.Error("uppercase level should be recognised")
	}
	if f.Allow(`{"level":"DEBUG","msg":"low"}`) {
		t.Error("DEBUG should be rejected when min is INFO")
	}
}
