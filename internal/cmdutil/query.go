package cmdutil

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// NewQuery builds url.Values from well-known list flags.
// DM platform uses cursor/limit for pagination (cursor is a skip offset).
func NewQuery(cmd *cobra.Command) url.Values {
	q := url.Values{}
	flags := cmd.Flags()

	// cursor: skip offset
	if cursor, err := flags.GetInt("cursor"); err == nil && cursor > 0 {
		q.Set("cursor", strconv.Itoa(cursor))
	}

	// limit
	if limit, err := flags.GetInt("limit"); err == nil {
		q.Set("limit", strconv.Itoa(limit))
	}

	// verbose
	if verbose, err := flags.GetInt("verbose"); err == nil && verbose > 0 {
		q.Set("verbose", strconv.Itoa(verbose))
	}

	return q
}

// ListFlags holds common pagination flags for list commands.
type ListFlags struct {
	Cursor  int
	Limit   int
	Verbose int
}

func (lf *ListFlags) Register(cmd *cobra.Command) {
	cmd.Flags().IntVar(&lf.Cursor, "cursor", 0, "Skip N items (pagination offset)")
	cmd.Flags().IntVar(&lf.Limit, "limit", 20, "Number of items per page")
	cmd.Flags().IntVar(&lf.Verbose, "verbose", 10, "Detail level (1-100, higher = more fields)")

	// AI-friendly aliases for --limit
	for _, alias := range []string{"page-size", "per-page"} {
		cmd.Flags().IntVar(&lf.Limit, alias, 20, "Alias for --limit")
		_ = cmd.Flags().MarkHidden(alias)
	}
}

// ApplyTo writes list flag values into a url.Values.
func (lf *ListFlags) ApplyTo(q url.Values) {
	if lf.Cursor > 0 {
		q.Set("cursor", strconv.Itoa(lf.Cursor))
	}
	if lf.Limit > 0 {
		q.Set("limit", strconv.Itoa(lf.Limit))
	}
	if lf.Verbose > 0 {
		q.Set("verbose", strconv.Itoa(lf.Verbose))
	}
}

// SetQueryParam sets a query param only if the value is non-empty.
func SetQueryParam(q url.Values, key, value string) {
	if strings.TrimSpace(value) != "" {
		q.Set(key, value)
	}
}
