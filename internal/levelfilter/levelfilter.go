// Package levelfilter provides filtering of JSON log lines by log level.
package levelfilter

import (
	"encoding/json"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var levelNames = map[string]Level{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

// Filter discards log lines whose level is below the configured minimum.
type Filter struct {
	min      Level
	disabled bool
	fields   []string
}

// New returns a Filter that passes lines at or above minLevel.
// If minLevel is empty, the filter is disabled and all lines pass.
func New(minLevel string, fields []string) (*Filter, error) {
	if minLevel == "" {
		return &Filter{disabled: true}, nil
	}
	lvl, ok := levelNames[strings.ToLower(minLevel)]
	if !ok {
		return nil, &UnknownLevelError{Name: minLevel}
	}
	f := []string{"level", "lvl", "severity"}
	if len(fields) > 0 {
		f = fields
	}
	return &Filter{min: lvl, fields: f}, nil
}

// Allow returns true if the line should be passed through.
func (f *Filter) Allow(line string) bool {
	if f.disabled {
		return true
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		// Non-JSON lines always pass through.
		return true
	}
	for _, field := range f.fields {
		val, ok := obj[field]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		lvl, ok := levelNames[strings.ToLower(s)]
		if !ok {
			return true
		}
		return lvl >= f.min
	}
	// No recognised level field found — pass through.
	return true
}

// UnknownLevelError is returned when an unrecognised level name is supplied.
type UnknownLevelError struct {
	Name string
}

func (e *UnknownLevelError) Error() string {
	return "unknown log level: " + e.Name
}
