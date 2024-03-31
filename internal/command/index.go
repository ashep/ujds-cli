package command

import (
	ujdscli "github.com/ashep/ujds/sdk/client"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/ashep/ujds-cli/internal/index/list"
)

func newIndex(cli *ujdscli.Client, l zerolog.Logger) *cobra.Command {

	c := &cobra.Command{
		Use:   "index",
		Short: "Index operations",
	}

	c.AddCommand(newIndexList(cli, l))

	return c
}

func newIndexList(cli *ujdscli.Client, l zerolog.Logger) *cobra.Command {
	var (
		names  *[]string
		format *string
	)

	c := &cobra.Command{
		Use:   "list",
		Short: "List indices",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.New(cli, l).List(cmd.Context(), *names, *format, cmd.Println)
		},
	}

	names = c.Flags().StringSliceP("names", "n", nil, "Index name patterns to list")
	format = c.Flags().StringP("format", "f", "{name}", "Output format; allowed variables: {name}, {title}")

	return c
}
