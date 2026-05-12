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

var statusMap = map[string]string{
	"running":   "1",
	"waiting":   "0,4,5",
	"failed":    "-1,2",
	"completed": "3",
}

var typeMap = map[string]string{
	"config-apply":        "1",
	"interactive-command": "2",
	"fetch-config":        "4",
	"import-firmware":     "6",
	"vpn-channel":         "12",
	"vpn-link-order":      "13",
	"token-cleanup":       "15",
	"traffic-stats":       "18",
	"idle-notice":         "20",
	"remote-web":          "23",
}

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
	Status   string
	TaskType string
	ObjectID string
}

func newCmdTaskList(f *factory.Factory) *cobra.Command {
	opts := &TaskListOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tasks",
		Aliases: []string{"ls"},
		Example: `  # List running tasks
  devicemanager task list --status running

  # List waiting tasks
  devicemanager task list --status waiting

  # List failed tasks
  devicemanager task list --status failed

  # List completed tasks
  devicemanager task list --status completed

  # Filter by type
  devicemanager task list --type config-apply
  devicemanager task list --type interactive-command`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.APIClient()
			if err != nil {
				return err
			}

			q := url.Values{}
			opts.ApplyTo(q)

			if opts.Status != "" {
				states, ok := statusMap[opts.Status]
				if !ok {
					return fmt.Errorf("invalid status %q, valid values: running, waiting, failed, completed", opts.Status)
				}
				q.Set("states", states)
			}

			if opts.TaskType != "" {
				typeVal, ok := typeMap[opts.TaskType]
				if !ok {
					return fmt.Errorf("invalid type %q, valid values: config-apply, interactive-command, fetch-config, import-firmware, vpn-channel, vpn-link-order, token-cleanup, traffic-stats, idle-notice, remote-web", opts.TaskType)
				}
				q.Set("types", typeVal)
			}

			cmdutil.SetQueryParam(q, "object_id", opts.ObjectID)

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
	cmd.Flags().StringVar(&opts.Status, "status", "", "Filter by status: running, waiting, failed, completed")
	cmd.Flags().StringVar(&opts.TaskType, "type", "", "Filter by type: config-apply, interactive-command, fetch-config, import-firmware, vpn-channel, vpn-link-order, token-cleanup, traffic-stats, idle-notice, remote-web")
	cmd.Flags().StringVar(&opts.ObjectID, "object-id", "", "Filter by device ID")

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
