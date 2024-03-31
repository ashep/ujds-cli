package command

import (
	ujdscli "github.com/ashep/ujds/sdk/client"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/ashep/ujds-cli/internal/export"
)

func newExport(cli *ujdscli.Client, l zerolog.Logger) *cobra.Command {
	var (
		indices   *[]string
		filename  *string
		overwrite *bool
	)

	c := &cobra.Command{
		Use:   "export",
		Short: "Export records",
		RunE: func(cmd *cobra.Command, args []string) error {
			return export.New(cli, l).Export(cmd.Context(), *indices, *filename, *overwrite)
		},
	}

	indices = c.Flags().StringSliceP("index", "i", []string{"*"}, "Index name patterns to scan")
	filename = c.Flags().StringP("out", "o", "out.csv", "Output file path")
	overwrite = c.Flags().Bool("overwrite", false, "Overwrite exiting file")

	return c
}
