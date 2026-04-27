package iostreams

import (
	"bytes"
	"encoding/json"

	"github.com/tidwall/pretty"
)

var jsonStyle = &pretty.Style{
	Key:    [2]string{"\x1b[1m", "\x1b[0m"},
	String: [2]string{"\x1b[32m", "\x1b[0m"},
	Number: [2]string{"\x1b[33m", "\x1b[0m"},
	True:   [2]string{"\x1b[31m", "\x1b[0m"},
	False:  [2]string{"\x1b[31m", "\x1b[0m"},
	Null:   [2]string{"\x1b[31m", "\x1b[0m"},
	Escape: [2]string{"\x1b[35m", "\x1b[0m"},
}

func FormatJSON(data []byte, io *IOStreams, outputMode string) string {
	if !io.IsStdoutTTY() {
		var buf bytes.Buffer
		if err := json.Compact(&buf, data); err != nil {
			return string(data)
		}
		return buf.String()
	}
	if outputMode == "json" {
		var buf bytes.Buffer
		if err := json.Indent(&buf, data, "", "  "); err != nil {
			return string(data)
		}
		return buf.String()
	}
	if !json.Valid(data) {
		return string(data)
	}
	result := pretty.Pretty(data)
	return string(pretty.Color(result, jsonStyle))
}
