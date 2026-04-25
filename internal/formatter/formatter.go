package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// Format controls the output style.
type Format string

const (
	FormatPretty Format = "pretty"
	FormatJSON   Format = "json"
	FormatRaw    Format = "raw"
)

// Formatter writes log entries to an output writer.
type Formatter struct {
	format Format
	w      io.Writer
}

// New creates a new Formatter writing to w.
func New(w io.Writer, format Format) *Formatter {
	return &Formatter{w: w, format: format}
}

// Write formats and writes a raw JSON log line, optionally annotated with a source label.
func (f *Formatter) Write(label, line string) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(label, line)
	case FormatPretty:
		return f.writePretty(label, line)
	default:
		return f.writeRaw(label, line)
	}
}

func (f *Formatter) writeRaw(label, line string) error {
	if label != "" {
		_, err := fmt.Fprintf(f.w, "[%s] %s\n", label, line)
		return err
	}
	_, err := fmt.Fprintln(f.w, line)
	return err
}

func (f *Formatter) writeJSON(label, line string) error {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return f.writeRaw(label, line)
	}
	if label != "" {
		obj["_source"] = label
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f.w, "%s\n", b)
	return err
}

func (f *Formatter) writePretty(label, line string) error {
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return f.writeRaw(label, line)
	}

	var sb strings.Builder
	if label != "" {
		sb.WriteString(fmt.Sprintf("\033[36m[%s]\033[0m ", label))
	}
	if ts, ok := obj["time"]; ok {
		if s, ok := ts.(string); ok {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				sb.WriteString(fmt.Sprintf("\033[90m%s\033[0m ", t.Format("15:04:05")))
			}
		}
	}
	if lvl, ok := obj["level"]; ok {
		sb.WriteString(levelColor(fmt.Sprintf("%v", lvl)))
		sb.WriteString(" ")
	}
	if msg, ok := obj["msg"]; ok {
		sb.WriteString(fmt.Sprintf("\033[1m%v\033[0m", msg))
	}

	keys := make([]string, 0, len(obj))
	for k := range obj {
		if k == "time" || k == "level" || k == "msg" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf(" \033[33m%s\033[0m=%v", k, obj[k]))
	}
	_, err := fmt.Fprintln(f.w, sb.String())
	return err
}

func levelColor(level string) string {
	switch strings.ToLower(level) {
	case "error", "err", "fatal":
		return fmt.Sprintf("\033[31m%-5s\033[0m", strings.ToUpper(level))
	case "warn", "warning":
		return fmt.Sprintf("\033[33m%-5s\033[0m", strings.ToUpper(level))
	case "debug":
		return fmt.Sprintf("\033[35m%-5s\033[0m", strings.ToUpper(level))
	default:
		return fmt.Sprintf("\033[32m%-5s\033[0m", strings.ToUpper(level))
	}
}

// ParseFormat converts a string to a Format value. It returns FormatRaw and an
// error if the string does not match a known format.
func ParseFormat(s string) (Format, error) {
	switch Format(strings.ToLower(s)) {
	case FormatPretty:
		return FormatPretty, nil
	case FormatJSON:
		return FormatJSON, nil
	case FormatRaw:
		return FormatRaw, nil
	default:
		return FormatRaw, fmt.Errorf("unknown format %q: must be one of pretty, json, raw", s)
	}
}
