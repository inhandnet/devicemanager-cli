package drc

import (
	"net/url"

	"github.com/spf13/cobra"

	"github.com/inhandnet/devicemanager-cli/internal/cmdutil"
	"github.com/inhandnet/devicemanager-cli/internal/factory"
)

// oidQuery returns a url.Values with the org ID set if impersonating.
func oidQuery(f *factory.Factory) url.Values {
	q := url.Values{}
	cmdutil.SetQueryParam(q, "oid", f.OrgID())
	return q
}

func NewCmdDRC(f *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drc",
		Short: "Manage DRC configuration templates",
	}

	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdDevices(f))

	return cmd
}
