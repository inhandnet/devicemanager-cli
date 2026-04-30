package task

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
	"github.com/inhandnet/devicemanager-cli/internal/iostreams"
	"github.com/inhandnet/devicemanager-cli/internal/ui"
)

func NewCmdTask(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks (DRC, firmware upgrade, etc.)",
	}

	cmd.AddCommand(newCmdTaskList(f))
	cmd.AddCommand(newCmdTaskCancel(f))
	cmd.AddCommand(newCmdTaskRestart(f))

	return cmd
}

type TaskListOptions struct {
	cmdutil.ListFlags
	Status     string
	TaskType   string
	DeviceName string
}

func newCmdTaskList(f *factory.Factory) *cobra.Command {
	opts := &TaskListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tasks",
		Aliases: []string{"ls"},
		Example: `  # List all running tasks
  devicemanager task list --status running

  # List failed tasks
  devicemanager task list --status failed

  # Filter by type
  devicemanager task list --type firmware_upgrade`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ApplyTo(q)
			q.Set("verbose", "50")
			cmdutil.SetQueryParam(q, "status", opts.Status)
			cmdutil.SetQueryParam(q, "type", opts.TaskType)
			cmdutil.SetQueryParam(q, "device_name", opts.DeviceName)

			output, _ := cmd.Flags().GetString("output")

			body, err := client.Get("/api2/tasks", q)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(body, f.IO, output,
				iostreams.WithColumns("_id", "deviceName", "type", "status",
					"progress", "creator", "createdAt", "startedAt", "updatedAt"))
		},
	}

	opts.Register(cmd)
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by status (running/waiting/failed/completed)")
	cmd.Flags().StringVar(&opts.TaskType, "type", "", "Filter by task type")
	cmd.Flags().StringVar(&opts.DeviceName, "device-name", "", "Filter by device name")

	return cmd
}

func newCmdTaskCancel(f *factory.Factory) *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "cancel <task-id>",
		Short: "Cancel a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !ui.Confirm(f.IO, fmt.Sprintf("Cancel task %s?", args[0]), yes) {
				return nil
			}

			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Delete(fmt.Sprintf("/api2/tasks/%s", args[0]))
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Skip confirmation prompt")

	return cmd
}

func newCmdTaskRestart(f *factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "restart <task-id>",
		Short: "Restart a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			output, _ := cmd.Flags().GetString("output")

			resp, err := client.Post(fmt.Sprintf("/api2/tasks/%s/restart", args[0]), nil)
			if err != nil {
				return err
			}

			return iostreams.FormatOutput(resp, f.IO, output)
		},
	}
}
