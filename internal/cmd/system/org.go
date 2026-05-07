package system

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
)

func NewCmdOrg(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "View and manage organization info",
	}

	cmd.AddCommand(newCmdOrgList(f))
	cmd.AddCommand(newCmdOrgGet(f))
	cmd.AddCommand(newCmdOrgUpdate(f))

	return cmd
}

func newCmdOrgList(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List organizations",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			q.Set("verbose", "100")

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api2/organizations", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "email"))
		},
	}
}

func newCmdOrgGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get current organization info",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			q.Set("verbose", "100")

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api2/organizations/this", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}
}

func newCmdOrgUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <org-id>",
		Short: "Update organization info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{}
			if cmd.Flags().Changed("name") {
				v, _ := cmd.Flags().GetString("name")
				body["name"] = v
			}
			if cmd.Flags().Changed("address") {
				v, _ := cmd.Flags().GetString("address")
				body["address"] = v
			}
			if cmd.Flags().Changed("contact") {
				v, _ := cmd.Flags().GetString("contact")
				body["contact"] = v
			}
			if cmd.Flags().Changed("email") {
				v, _ := cmd.Flags().GetString("email")
				body["email"] = v
			}
			if cmd.Flags().Changed("country") {
				v, _ := cmd.Flags().GetString("country")
				body["country"] = v
			}
			if cmd.Flags().Changed("biz-category") {
				v, _ := cmd.Flags().GetString("biz-category")
				body["bizCategory"] = v
			}

			if len(body) == 0 {
				return fmt.Errorf("at least one flag is required")
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api2/organizations/%s", args[0]), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("name", "", "Organization name")
	cmd.Flags().String("address", "", "Organization address")
	cmd.Flags().String("contact", "", "Contact information")
	cmd.Flags().String("email", "", "Organization email")
	cmd.Flags().String("country", "", "Country")
	cmd.Flags().String("biz-category", "", "Business category")

	return cmd
}
