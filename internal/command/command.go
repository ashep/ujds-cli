package command

import (
	ujdscli "github.com/ashep/ujds/sdk/client"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func New(cli *ujdscli.Client, l zerolog.Logger) *cobra.Command {
	root := &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.AddCommand(newExport(cli, l))
	root.AddCommand(newIndex(cli, l))

	return root
}
