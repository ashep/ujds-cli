package command

import (
	ujdscli "github.com/ashep/ujds/sdk/client"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/ashep/ujds-cli/internal/export"
)

func New(cli *ujdscli.Client, l zerolog.Logger) *cobra.Command {
	root := &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	root.AddCommand(newExport(cli, l))

	return root
}

func newExport(cli *ujdscli.Client, l zerolog.Logger) *cobra.Command {
	var (
		indices   *[]string
		filename  *string
		format    *string
		overwrite *bool
	)

	c := &cobra.Command{
		Use:   "export",
		Short: "export data",
		RunE: func(cmd *cobra.Command, args []string) error {
			return export.New(cli, l).Export(cmd.Context(), *indices, *filename, *format, *overwrite)
		},
	}

	indices = c.Flags().StringSliceP("index", "i", nil, "index name patterns to scan")
	filename = c.Flags().StringP("out", "o", "out.{FORMAT}", "output path")
	format = c.Flags().StringP("format", "f", "csv", "output format")
	overwrite = c.Flags().Bool("overwrite", false, "overwrite exiting file")

	_ = c.MarkFlagRequired("index")

	return c
}
