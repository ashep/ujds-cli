package command

import (
	"github.com/ashep/ujds-cli/internal/record/find"
	"github.com/ashep/ujds-cli/internal/record/history"
	ujdscli "github.com/ashep/ujds/sdk/client"
	"github.com/spf13/cobra"
)

func newRecord(cli *ujdscli.Client) *cobra.Command {

	c := &cobra.Command{
		Use:   "record",
		Short: "Record operations",
	}

	c.AddCommand(newRecordFind(cli))
	c.AddCommand(newRecordHistory(cli))

	return c
}

func newRecordFind(cli *ujdscli.Client) *cobra.Command {
	var (
		index  *string
		query  *string
		format *string
		since  *int64
		limit  *uint32
		cursor *uint64
	)

	c := &cobra.Command{
		Use:   "find",
		Short: "Find records",
		RunE: func(cmd *cobra.Command, args []string) error {
			return find.New(cli).Find(cmd.Context(), *index, *query, *format, *since, *limit, *cursor, cmd.OutOrStdout())
		},
	}

	index = c.Flags().StringP("index", "i", "", "Index name")
	query = c.Flags().StringP("query", "q", "", "Search query")
	since = c.Flags().Int64P("since", "s", 0, "Since")
	limit = c.Flags().Uint32P("limit", "l", 100, "Limit")
	cursor = c.Flags().Uint64P("cursor", "c", 0, "Cursor")
	format = c.Flags().StringP("format", "f", "ID: {{.Id}}\nCreated: {{.CreatedAt}}\nUpdated: {{.UpdatedAt}}\nTouched: {{.TouchedAt}}\nData: {{.Data}}\n\n", "Output format")

	_ = c.MarkFlagRequired("index")

	return c
}

func newRecordHistory(cli *ujdscli.Client) *cobra.Command {
	var (
		index  *string
		id     *string
		format *string
		since  *int64
		limit  *uint32
		cursor *uint64
	)

	c := &cobra.Command{
		Use:   "history",
		Short: "Get record history",
		RunE: func(cmd *cobra.Command, args []string) error {
			return history.New(cli).History(cmd.Context(), *index, *id, *format, *since, *limit, *cursor, cmd.OutOrStdout())
		},
	}

	index = c.Flags().StringP("index", "i", "", "Index name")
	id = c.Flags().String("id", "", "Record ID")
	since = c.Flags().Int64P("since", "s", 0, "Since")
	limit = c.Flags().Uint32P("limit", "l", 100, "Limit")
	cursor = c.Flags().Uint64P("cursor", "c", 0, "Cursor")
	format = c.Flags().StringP("format", "f", "ID: {{.ID}}\nTime: {{.TimeStr}}\nData: {{.Data}}\n\n", "Output format")

	_ = c.MarkFlagRequired("index")

	return c
}
