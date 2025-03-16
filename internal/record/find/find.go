package find

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/template"

	ujdscli "github.com/ashep/ujds/sdk/client"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/bufbuild/connect-go"
)

type Find struct {
	cli *ujdscli.Client
}

func New(cli *ujdscli.Client) *Find {
	return &Find{
		cli: cli,
	}
}

func (f *Find) Find(ctx context.Context, index, query, format string, since int64, limit uint32, cursor uint64, out io.Writer) error {
	var err error

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	tpl := template.New("find")
	if tpl, err = tpl.Parse(format); err != nil {
		return fmt.Errorf("parse format: %w", err)
	}

	res, err := f.cli.R.Find(ctx, connect.NewRequest(&recordproto.FindRequest{
		Index:  index,
		Search: query,
		Since:  since,
		Limit:  limit,
		Cursor: cursor,
	}))

	if err != nil {
		return fmt.Errorf("ujds response: %w", err)
	}

	for _, rec := range res.Msg.Records {
		if err := tpl.Execute(out, rec); err != nil {
			return fmt.Errorf("template execute: %w", err)
		}
	}

	return nil
}
