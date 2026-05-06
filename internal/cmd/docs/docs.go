package docs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

const (
	defaultRepo   = "inhandnet/model-reference"
	defaultBranch = "main"
)

func repoBase() string {
	if v := os.Getenv("DEVICEMANAGER_DOCS_REPO"); v != "" {
		return v
	}
	return defaultRepo
}

func rawURL(path string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repoBase(), defaultBranch, path)
}

func fetchRaw(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

func NewCmdDocs(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "docs",
		Short: "Browse device model reference documentation",
		Long: `Browse device-side documentation for configuration and troubleshooting.

Documents are organized by device model in a public GitHub repository.
Each model has its own index.md describing available documents.

Reading flow: index.md (root) → <model>/index.md → specific document.

Override the default repo with DEVICEMANAGER_DOCS_REPO environment variable.`,
	}

	cmd.AddCommand(newCmdDocsList(f))
	cmd.AddCommand(newCmdDocsGet(f))
	cmd.AddCommand(newCmdDocsSearch(f))

	return cmd
}

func newCmdDocsList(f *factory.Factory) *cobra.Command {
	var model string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List available models or model-specific documents",
		Aliases: []string{"ls"},
		Example: `  # List all available models
  devicemanager docs list

  # List documents for a specific model
  devicemanager docs list --model ER805`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "index.md"
			if model != "" {
				path = model + "/index.md"
			}

			data, err := fetchRaw(rawURL(path))
			if err != nil {
				if model != "" {
					return fmt.Errorf("model %q not found or has no index", model)
				}
				return fmt.Errorf("fetching document index: %w", err)
			}

			fmt.Fprint(f.IO.Out, string(data))
			return nil
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Show index for a specific model (e.g. ER805, IR615)")

	return cmd
}

func newCmdDocsGet(f *factory.Factory) *cobra.Command {
	var model string

	cmd := &cobra.Command{
		Use:   "get <path>",
		Short: "Fetch a specific document by path",
		Args:  cobra.ExactArgs(1),
		Example: `  # Fetch with full path
  devicemanager docs get ER805/troubleshooting/vpn.md

  # Fetch with --model prefix (equivalent to above)
  devicemanager docs get troubleshooting/vpn.md --model ER805`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			if model != "" && !strings.Contains(path, "/") || (model != "" && !strings.HasPrefix(path, model)) {
				path = model + "/" + path
			}

			data, err := fetchRaw(rawURL(path))
			if err != nil {
				return err
			}

			fmt.Fprint(f.IO.Out, string(data))
			return nil
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Prepend model directory to path (e.g. ER805)")

	return cmd
}

type docEntry struct {
	Title string
	Path  string
	Desc  string
	Line  string
}

// parseIndexEntries extracts doc entries from markdown index content.
// Supports formats:
//   - [Title](path) — description
//   - [Title](path)
//   - plain text lines containing the keyword
func parseIndexEntries(data []byte, keyword string) []docEntry {
	var results []docEntry
	for _, line := range strings.Split(string(data), "\n") {
		lower := strings.ToLower(line)
		if !strings.Contains(lower, keyword) {
			continue
		}

		entry := docEntry{Line: strings.TrimSpace(line)}

		// Try to parse markdown link: [Title](path)
		if start := strings.Index(line, "["); start >= 0 {
			if end := strings.Index(line[start:], "]("); end >= 0 {
				entry.Title = line[start+1 : start+end]
				rest := line[start+end+2:]
				if pEnd := strings.Index(rest, ")"); pEnd >= 0 {
					entry.Path = rest[:pEnd]
					after := strings.TrimSpace(rest[pEnd+1:])
					after = strings.TrimLeft(after, "—–-")
					entry.Desc = strings.TrimSpace(after)
				}
			}
		}

		results = append(results, entry)
	}
	return results
}

func newCmdDocsSearch(f *factory.Factory) *cobra.Command {
	var model string

	cmd := &cobra.Command{
		Use:   "search <keyword>",
		Short: "Search documents by keyword in model index",
		Args:  cobra.ExactArgs(1),
		Example: `  # Search across a model's documents
  devicemanager docs search VPN --model ER805

  # Search root index (lists matching models)
  devicemanager docs search ER805

  # Search multiple keywords (quote them)
  devicemanager docs search "port forwarding" --model IR615`,
		RunE: func(cmd *cobra.Command, args []string) error {
			keyword := strings.ToLower(args[0])

			type indexSource struct {
				label  string
				path   string
				prefix string // prepend to relative paths in results
			}

			sources := []indexSource{
				{label: "root", path: "index.md"},
			}
			if model != "" {
				sources = append(sources, indexSource{
					label:  model,
					path:   model + "/index.md",
					prefix: model + "/",
				})
			}

			out := f.IO.Out
			var allEntries []docEntry

			for _, src := range sources {
				data, err := fetchRaw(rawURL(src.path))
				if err != nil {
					continue
				}
				entries := parseIndexEntries(data, keyword)
				for i := range entries {
					// Prefix relative paths with model directory
					if src.prefix != "" && entries[i].Path != "" && !strings.Contains(entries[i].Path, "/") {
						entries[i].Path = src.prefix + entries[i].Path
					}
					if entries[i].Desc == "" {
						entries[i].Desc = "[" + src.label + "]"
					} else {
						entries[i].Desc = "[" + src.label + "] " + entries[i].Desc
					}
				}
				allEntries = append(allEntries, entries...)
			}

			if len(allEntries) == 0 {
				fmt.Fprintf(out, "No documents found matching %q\n", args[0])
				if model == "" {
					fmt.Fprintln(out, "Tip: use --model <model> to search within a model's documents")
				}
				return nil
			}

			fmt.Fprintf(out, "%s Found %d result(s) for %q:\n\n",
				iostreams.Green("✓"), len(allEntries), args[0])

			for i, e := range allEntries {
				if e.Path != "" {
					fmt.Fprintf(out, "  %d. %s\n     Path: %s\n",
						i+1, e.Title, e.Path)
					if e.Desc != "" {
						fmt.Fprintf(out, "     %s\n", e.Desc)
					}
				} else {
					fmt.Fprintf(out, "  %d. %s\n", i+1, e.Line)
				}
				fmt.Fprintln(out)
			}

			fmt.Fprintln(out, "Use 'devicemanager docs get <path>' to read a document.")
			return nil
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Search within a model's index (e.g. ER805, IR615)")

	return cmd
}
