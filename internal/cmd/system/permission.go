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

func NewCmdPermission(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "permission",
		Short:   "Manage device permissions",
		Aliases: []string{"perm"},
	}

	cmd.AddCommand(newCmdPermissionList(f))
	cmd.AddCommand(newCmdPermissionGet(f))
	cmd.AddCommand(newCmdPermissionCreate(f))
	cmd.AddCommand(newCmdPermissionUpdate(f))
	cmd.AddCommand(newCmdPermissionDelete(f))
	cmd.AddCommand(newCmdPermissionUsers(f))
	cmd.AddCommand(newCmdPermissionDevices(f))

	return cmd
}

func newCmdPermissionList(f *factory.Factory) *cobra.Command {
	var flags cmdutil.ListFlags

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List device permission groups",
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

			body, err := client.Get("/api/groups", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "description", "deviceCount"))
		},
	}

	flags.Register(cmd)

	return cmd
}

func newCmdPermissionGet(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <permission-id>",
		Short: "Get device permission group details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/groups/%s", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output)
		},
	}
}

func newCmdPermissionCreate(f *factory.Factory) *cobra.Command {
	var (
		name string
		desc string
	)

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a device permission group",
		Example: `  devicemanager system permission create --name "office-devices" --desc "Office routers"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"name": name,
			}
			if desc != "" {
				body["description"] = desc
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post("/api/groups", body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Permission group name (required)")
	cmd.Flags().StringVar(&desc, "desc", "", "Description")
	_ = cmd.MarkFlagRequired("name")

	return cmd
}

func newCmdPermissionDelete(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <permission-id>",
		Short: "Delete a device permission group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ui.Confirm(f.IO, fmt.Sprintf("Delete permission group %s?", args[0]), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api/groups/%s", args[0]))
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}

func newCmdPermissionUpdate(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <group-id>",
		Short: "Update a device permission group",
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

			if len(body) == 0 {
				return fmt.Errorf("at least one flag is required")
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Put(fmt.Sprintf("/api/groups/%s", args[0]), body)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().String("name", "", "New group name")

	return cmd
}

func newCmdPermissionUsers(f *factory.Factory) *cobra.Command {
	var flags cmdutil.ListFlags

	cmd := &cobra.Command{
		Use:   "users <group-id>",
		Short: "List users in a permission group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			flags.ApplyTo(q)
			q.Set("verbose", "100")

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/groups/%s/users", args[0]), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "email", "roleName"))
		},
	}

	flags.Register(cmd)

	return cmd
}

func newCmdPermissionDevices(f *factory.Factory) *cobra.Command {
	var flags cmdutil.ListFlags

	cmd := &cobra.Command{
		Use:   "devices <group-id>",
		Short: "List devices in a permission group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			flags.ApplyTo(q)
			q.Set("verbose", "100")

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get(fmt.Sprintf("/api/groups/%s/devices", args[0]), q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "name", "serialNumber", "model", "online"))
		},
	}

	flags.Register(cmd)

	return cmd
}
