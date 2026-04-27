package iostreams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/tidwall/gjson"
)

var borderlessStyle = table.Style{
	Name: "Borderless",
	Box: table.BoxStyle{
		PaddingLeft:  "",
		PaddingRight: "  ",
	},
	Format: table.FormatOptions{
		Header: text.FormatDefault,
		Row:    text.FormatDefault,
	},
	Options: table.Options{
		DrawBorder:      false,
		SeparateColumns: false,
		SeparateHeader:  false,
		SeparateRows:    false,
		SeparateFooter:  false,
	},
}

type TablePrinter struct {
	out   io.Writer
	isTTY bool
	rows  [][]string
}

func NewTablePrinter(out io.Writer, isTTY bool) *TablePrinter {
	return &TablePrinter{out: out, isTTY: isTTY}
}

func (t *TablePrinter) AddRow(cols ...string) {
	t.rows = append(t.rows, cols)
}

func (t *TablePrinter) Render() error {
	if len(t.rows) == 0 {
		return nil
	}
	if !t.isTTY {
		return t.renderTSV()
	}
	return t.renderTable()
}

func (t *TablePrinter) renderTSV() error {
	for _, row := range t.rows {
		_, err := fmt.Fprintln(t.out, strings.Join(row, "\t"))
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TablePrinter) renderTable() error {
	tw := table.NewWriter()
	tw.SetStyle(borderlessStyle)
	for _, row := range t.rows {
		tableRow := make(table.Row, len(row))
		for i, col := range row {
			tableRow[i] = col
		}
		tw.AppendRow(tableRow)
	}
	rendered := tw.Render()
	lines := strings.Split(rendered, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	_, err := fmt.Fprintln(t.out, strings.Join(lines, "\n"))
	return err
}

// FormatTable renders JSON data as a table.
// It handles both {"result": [...]} envelopes and bare arrays/objects.
// priorityCols controls the order of displayed columns.
func FormatTable(data []byte, ios *IOStreams, priorityCols []string) error {
	tp := NewTablePrinter(ios.Out, ios.IsStdoutTTY())

	// Try to extract result array from envelope
	result := gjson.GetBytes(data, "result")
	if !result.Exists() {
		result = gjson.ParseBytes(data)
	}

	if result.IsArray() {
		items := result.Array()
		if len(items) == 0 {
			return nil
		}
		// Collect all keys from first item
		cols := collectKeys(items[0], priorityCols)
		// Header row
		header := make([]string, len(cols))
		for i, c := range cols {
			header[i] = strings.ToUpper(c)
		}
		tp.AddRow(header...)
		// Data rows
		for _, item := range items {
			row := make([]string, len(cols))
			for i, c := range cols {
				row[i] = formatValue(item.Get(c))
			}
			tp.AddRow(row...)
		}
	} else if result.IsObject() {
		// Single object: key-value pairs
		result.ForEach(func(key, value gjson.Result) bool {
			tp.AddRow(key.String(), formatValue(value))
			return true
		})
	}

	return tp.Render()
}

func collectKeys(item gjson.Result, priority []string) []string {
	seen := map[string]bool{}
	var cols []string
	// Priority columns first
	for _, c := range priority {
		if item.Get(c).Exists() {
			cols = append(cols, c)
			seen[c] = true
		}
	}
	// Remaining columns alphabetically
	var remaining []string
	item.ForEach(func(key, _ gjson.Result) bool {
		k := key.String()
		if !seen[k] {
			remaining = append(remaining, k)
		}
		return true
	})
	sort.Strings(remaining)
	return append(cols, remaining...)
}

func formatValue(v gjson.Result) string {
	switch v.Type {
	case gjson.Null:
		return ""
	case gjson.JSON:
		var buf bytes.Buffer
		if err := json.Compact(&buf, []byte(v.Raw)); err != nil {
			return v.Raw
		}
		return buf.String()
	default:
		return v.String()
	}
}
