package highlight

import (
	"fmt"
	"strings"
)

// Color ANSI escape codes.
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Bold    = "\033[1m"
)

// LevelColors maps log level strings to ANSI color codes.
var LevelColors = map[string]string{
	"debug": Cyan,
	"info":  Green,
	"warn":  Yellow,
	"error": Red,
	"fatal": Magenta,
}

// Highlighter applies ANSI color codes to log output.
type Highlighter struct {
	enabled bool
}

// New returns a new Highlighter. If enabled is false, all methods
// return the input unchanged (useful for non-TTY outputs).
func New(enabled bool) *Highlighter {
	return &Highlighter{enabled: enabled}
}

// Level colorizes a log level string based on its severity.
func (h *Highlighter) Level(level string) string {
	if !h.enabled {
		return level
	}
	color, ok := LevelColors[strings.ToLower(level)]
	if !ok {
		return level
	}
	return fmt.Sprintf("%s%s%s%s", Bold, color, level, Reset)
}

// Field colorizes a key=value pair, bolding the key.
func (h *Highlighter) Field(key, value string) string {
	if !h.enabled {
		return fmt.Sprintf("%s=%s", key, value)
	}
	return fmt.Sprintf("%s%s%s=%s", Bold, key, Reset, value)
}

// Source colorizes a source label for fan-in output.
func (h *Highlighter) Source(label string) string {
	if !h.enabled {
		return label
	}
	return fmt.Sprintf("%s%s%s", Blue, label, Reset)
}

// Message returns the message string, bolded when highlighting is enabled.
func (h *Highlighter) Message(msg string) string {
	if !h.enabled {
		return msg
	}
	return fmt.Sprintf("%s%s%s", Bold, msg, Reset)
}
