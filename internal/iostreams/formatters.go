package iostreams

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// ColumnFormatter transforms a string cell value for display.
type ColumnFormatter func(string) string

// ColumnFormatters maps column names (gjson paths) to formatters.
type ColumnFormatters map[string]ColumnFormatter

// FormatBytes converts a numeric string of bytes to human-readable form (e.g. "20 GiB").
func FormatBytes(s string) string {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	if f < 1 {
		return "0 B"
	}
	i := int(math.Log2(f) / 10)
	if i >= len(units) {
		i = len(units) - 1
	}
	val := f / math.Pow(1024, float64(i))
	if i == 0 {
		return fmt.Sprintf("%.0f B", val)
	}
	return fmt.Sprintf("%.1f %s", val, units[i])
}

// FormatDuration converts a numeric string of seconds to human-readable duration (e.g. "1d 2h 3m").
func FormatDuration(s string) string {
	secs, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	d := time.Duration(secs * float64(time.Second))
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}

	var parts []string
	if days := int(d.Hours()) / 24; days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
		d -= time.Duration(days) * 24 * time.Hour
	}
	if hours := int(d.Hours()); hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
		d -= time.Duration(hours) * time.Hour
	}
	if mins := int(d.Minutes()); mins > 0 {
		parts = append(parts, fmt.Sprintf("%dm", mins))
		d -= time.Duration(mins) * time.Minute
	}
	if secs := int(d.Seconds()); secs > 0 {
		parts = append(parts, fmt.Sprintf("%ds", secs))
	}
	return strings.Join(parts, " ")
}

// FormatRelativeTime converts an ISO 8601 timestamp to a relative time string (e.g. "2 hours ago").
func FormatRelativeTime(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		// Try Unix timestamp
		if unix, err2 := strconv.ParseInt(s, 10, 64); err2 == nil {
			t = time.Unix(unix, 0)
		} else {
			return s
		}
	}
	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		m := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", m)
	case diff < 24*time.Hour:
		h := int(diff.Hours())
		return fmt.Sprintf("%dh ago", h)
	case diff < 30*24*time.Hour:
		d := int(diff.Hours()) / 24
		return fmt.Sprintf("%dd ago", d)
	default:
		return t.Format("2006-01-02")
	}
}

// FormatPercent converts a decimal string to percentage (e.g. "0.452" → "45.2%").
func FormatPercent(s string) string {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	return fmt.Sprintf("%.1f%%", f*100)
}

// TruncateRunes truncates a string to maxLen runes, appending "…" if truncated.
func TruncateRunes(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen-1]) + "…"
}
