package iostreams

import (
	"encoding/json"
	"fmt"
)

type FormatOption func(*formatOptions)

type formatOptions struct {
	columns []string
}

func WithColumns(cols ...string) FormatOption {
	return func(o *formatOptions) {
		o.columns = cols
	}
}

// FormatOutput renders body according to the output mode (table/yaml/json/jq).
func FormatOutput(body []byte, io *IOStreams, output string, opts ...FormatOption) error {
	var o formatOptions
	for _, opt := range opts {
		opt(&o)
	}

	// --jq overrides output mode
	if io.JQExpr != "" {
		result, err := ApplyJQ(unwrapResult(body), io.JQExpr)
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
		return FormatTable(body, io, o.columns)
	case "yaml":
		s, err := FormatYAML(unwrapResult(body))
		if err != nil {
			return err
		}
		fmt.Fprintln(io.Out, s)
	default:
		if json.Valid(body) {
			fmt.Fprintln(io.Out, FormatJSON(unwrapResult(body), io, output))
		} else {
			fmt.Fprintln(io.Out, string(body))
		}
	}
	return nil
}

// unwrapResult strips the envelope when the JSON object has "result" as its only key.
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
