package history

import (
	"context"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	ujdscli "github.com/ashep/ujds/sdk/client"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/bufbuild/connect-go"
)

type History struct {
	cli *ujdscli.Client
}

type historyRecord struct {
	ID      string
	Time    time.Time
	TimeStr string
	Data    string
}

func New(cli *ujdscli.Client) *History {
	return &History{
		cli: cli,
	}
}

func (f *History) History(ctx context.Context, index, id, format string, since int64, limit uint32, cursor uint64, out io.Writer) error {
	var err error

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	tpl := template.New("find")
	if tpl, err = tpl.Parse(format); err != nil {
		return fmt.Errorf("parse format: %w", err)
	}

	res, err := f.cli.R.History(ctx, connect.NewRequest(&recordproto.HistoryRequest{
		Index:  index,
		Id:     id,
		Since:  since,
		Limit:  limit,
		Cursor: cursor,
	}))

	if err != nil {
		return fmt.Errorf("ujds response: %w", err)
	}

	for _, rec := range res.Msg.Records {
		hRec := historyRecord{
			ID:      rec.Id,
			Time:    time.Unix(rec.GetCreatedAt(), 0),
			TimeStr: time.Unix(rec.GetCreatedAt(), 0).Format(time.DateTime),
			Data:    rec.Data,
		}

		if err := tpl.Execute(out, hRec); err != nil {
			return fmt.Errorf("template execute: %w", err)
		}
	}

	return nil
}
