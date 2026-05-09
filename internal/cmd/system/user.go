package system

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

func NewCmdUser(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}

	cmd.AddCommand(newCmdUserList(f))
	cmd.AddCommand(newCmdUserGet(f))
	cmd.AddCommand(newCmdUserCreate(f))
	cmd.AddCommand(newCmdUserUpdate(f))
	cmd.AddCommand(newCmdUserDelete(f))

	return cmd
}

func newCmdUserList(f *factory.Factory) *cobra.Command {
	var flags cmdutil.ListFlags

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List users in the organization",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			flags.ApplyTo(q)
			q.Set("verbose", "100")

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api2/users", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "email", "userType", "phone", "status", "roleName"))
		},
	}

	flags.Register(cmd)

	return cmd
}

func newCmdUserGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <user-id>",
		Short: "Get user details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			q.Set("verbose", "100")

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api2/users/%s", args[0]), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}
}

func newCmdUserCreate(f *factory.Factory) *cobra.Command {
	var (
		name     string
		email    string
		roleID   string
		external bool
		lang     string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a user",
		Long:  "Create a user and send an invitation email. The user sets their password via the email link.",
		Example: `  # Create an internal user (invitation email sent automatically)
  devicemanager system user create --name "test" --email "test@example.com"

  # Create with a specific role ID
  devicemanager system user create --name "test" --email "test@example.com" --role-id <role-id>

  # Create an external user
  devicemanager system user create --email "ext@example.com" --external`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"email":    email,
				"external": external,
				"lang":     1, // English; overridden below for Chinese
			}
			if name != "" {
				body["name"] = name
			}
			if roleID != "" {
				body["roleId"] = roleID
			}
			if lang != "" {
				switch lang {
				case "zh", "cn":
					body["lang"] = 2
				default:
					body["lang"] = 1
				}
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api2/users", body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "User name")
	cmd.Flags().StringVar(&email, "email", "", "User email (required)")
	cmd.Flags().StringVar(&roleID, "role-id", "", "Role ID (required, use 'system role list' to find IDs)")
	cmd.Flags().BoolVar(&external, "external", false, "Create as external user")
	cmd.Flags().StringVar(&lang, "lang", "", `Invitation email language: "en" (default) or "zh"`)
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("role-id")

	return cmd
}

func newCmdUserUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <user-id>",
		Short: "Update a user",
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
			if cmd.Flags().Changed("role-id") {
				v, _ := cmd.Flags().GetString("role-id")
				body["roleId"] = v
			}

			if len(body) == 0 {
				return fmt.Errorf("at least one flag is required")
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api2/users/%s?verbose=100", args[0]), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("name", "", "New user name")
	cmd.Flags().String("role-id", "", "New role ID (use 'system role list' to find IDs)")

	return cmd
}

func newCmdUserDelete(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <user-id>",
		Short: "Delete a user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ui.Confirm(f.IO, fmt.Sprintf("Delete user %s?", args[0]), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api2/users/%s", args[0]))
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}
