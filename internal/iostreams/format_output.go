package iostreams

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// TransformFunc transforms raw JSON before table rendering.
type TransformFunc func([]byte) ([]byte, error)

// FormatOption configures FormatOutput behavior.
type FormatOption func(*formatOptions)

type formatOptions struct {
	columns    []string
	transform  TransformFunc
	formatters ColumnFormatters
}

func WithColumns(cols ...string) FormatOption {
	return func(o *formatOptions) {
		o.columns = cols
	}
}

func WithTransform(fn TransformFunc) FormatOption {
	return func(o *formatOptions) {
		o.transform = fn
	}
}

func WithFormatters(fmts ColumnFormatters) FormatOption {
	return func(o *formatOptions) {
		o.formatters = fmts
	}
}

// ChainTransforms composes multiple TransformFunc into one.
func ChainTransforms(fns ...TransformFunc) TransformFunc {
	return func(data []byte) ([]byte, error) {
		var err error
		for _, fn := range fns {
			data, err = fn(data)
			if err != nil {
				return nil, err
			}
		}
		return data, nil
	}
}

// FormatOutput renders body according to the output mode (table/yaml/json/jq).
func FormatOutput(body []byte, io *IOStreams, output string, opts ...FormatOption) error {
	var o formatOptions
	for _, opt := range opts {
		opt(&o)
	}

	// Check for empty results
	if isEmptyResult(body) {
		fmt.Fprintln(io.ErrOut, "No results.")
		return nil
	}

	// --jq overrides output mode
	if io.JQExpr != "" {
		result, err := ApplyJQ(unwrapResult(normalizePage(body)), io.JQExpr)
		if err != nil {
			return err
		}
		if result != "" {
			fmt.Fprintln(io.Out, result)
		}
		return nil
	}

	switch output {
	case "table":
		data := body
		if o.transform != nil {
			var err error
			data, err = o.transform(data)
			if err != nil {
				return err
			}
		}
		if o.formatters != nil {
			data = applyFormatters(data, o.formatters)
		}
		return FormatTable(data, io, o.columns)
	case "yaml":
		s, err := FormatYAML(unwrapResult(normalizePage(body)))
		if err != nil {
			return err
		}
		fmt.Fprintln(io.Out, s)
	default:
		normalized := normalizePage(body)
		if json.Valid(normalized) {
			fmt.Fprintln(io.Out, FormatJSON(unwrapResult(normalized), io, output))
		} else {
			fmt.Fprintln(io.Out, string(normalized))
		}
	}
	return nil
}

// normalizePage converts 0-based page numbers to 1-based in paginated responses.
func normalizePage(data []byte) []byte {
	parsed := gjson.ParseBytes(data)
	if !parsed.IsObject() {
		return data
	}
	pageVal := parsed.Get("page")
	if !pageVal.Exists() || pageVal.Type != gjson.Number {
		return data
	}
	// Only convert if there's also a "total" or "result" field (pagination envelope)
	if !parsed.Get("total").Exists() && !parsed.Get("result").Exists() {
		return data
	}
	page := pageVal.Int()
	result, err := sjson.SetBytes(data, "page", page+1)
	if err != nil {
		return data
	}
	return result
}

// unwrapResult strips the envelope when the JSON object has "result" as its only key.
// isEmptyResult checks if the response contains no data.
// Matches: [], {"result":[],...}, or {"result":null,...}
func isEmptyResult(data []byte) bool {
	parsed := gjson.ParseBytes(data)
	if parsed.IsArray() && len(parsed.Array()) == 0 {
		return true
	}
	if parsed.IsObject() {
		result := parsed.Get("result")
		if result.Exists() && result.IsArray() && len(result.Array()) == 0 {
			return true
		}
	}
	return false
}

func unwrapResult(data []byte) []byte {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return data
	}
	if len(raw) == 1 {
		if inner, ok := raw["result"]; ok {
			return inner
		}
	}
	return data
}

// applyFormatters applies column formatters to result data.
func applyFormatters(data []byte, fmts ColumnFormatters) []byte {
	parsed := gjson.ParseBytes(data)

	// Handle envelope: {"result": [...]}
	var items gjson.Result
	if parsed.IsObject() && parsed.Get("result").Exists() {
		items = parsed.Get("result")
	} else {
		items = parsed
	}

	if !items.IsArray() {
		// Single object
		return applyFormattersToObject(data, fmts, "")
	}

	// Array of objects
	result := data
	items.ForEach(func(key, value gjson.Result) bool {
		prefix := "result." + strconv.Itoa(int(key.Int()))
		if !parsed.IsObject() || !parsed.Get("result").Exists() {
			prefix = strconv.Itoa(int(key.Int()))
		}
		result = applyFormattersToObject(result, fmts, prefix)
		return true
	})
	return result
}

func applyFormattersToObject(data []byte, fmts ColumnFormatters, prefix string) []byte {
	for col, fn := range fmts {
		path := col
		if prefix != "" {
			path = prefix + "." + col
		}
		val := gjson.GetBytes(data, path)
		if !val.Exists() {
			continue
		}
		formatted := fn(val.String())
		var err error
		data, err = sjson.SetBytes(data, path, formatted)
		if err != nil {
			continue
		}
	}
	return data
}
