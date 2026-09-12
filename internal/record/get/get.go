package get

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

type Get struct {
	cli *ujdscli.Client
}

type getRecord struct {
	Id            string //nolint:revive // template field name mirrors the proto one
	Rev           uint64
	Index         string
	CreatedAt     string
	UpdatedAt     string
	TouchedAt     string
	CreatedAtUnix int64
	UpdatedAtUnix int64
	TouchedAtUnix int64
	Data          string
	DataTable     string
}

// utcTime formats a unix timestamp as a human readable UTC date and time.
func utcTime(ts int64) string {
	return time.Unix(ts, 0).UTC().Format(time.DateTime) + " UTC"
}

func New(cli *ujdscli.Client) *Get {
	return &Get{
		cli: cli,
	}
}

func (g *Get) Get(ctx context.Context, index, id, format string, out io.Writer) error {
	var err error

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	tpl := template.New("get")
	if tpl, err = tpl.Parse(format); err != nil {
		return fmt.Errorf("parse format: %w", err)
	}

	res, err := g.cli.R.Get(ctx, connect.NewRequest(&recordproto.GetRequest{
		Index: index,
		Id:    id,
	}))

	if err != nil {
		return fmt.Errorf("ujds response: %w", err)
	}

	rec := res.Msg.GetRecord()

	gRec := getRecord{
		Id:            rec.GetId(),
		Rev:           rec.GetRev(),
		Index:         rec.GetIndex(),
		CreatedAt:     utcTime(rec.GetCreatedAt()),
		UpdatedAt:     utcTime(rec.GetUpdatedAt()),
		TouchedAt:     utcTime(rec.GetTouchedAt()),
		CreatedAtUnix: rec.GetCreatedAt(),
		UpdatedAtUnix: rec.GetUpdatedAt(),
		TouchedAtUnix: rec.GetTouchedAt(),
		Data:          rec.GetData(),
		DataTable:     dataTable(rec.GetData()),
	}

	if err := tpl.Execute(out, gRec); err != nil {
		return fmt.Errorf("template execute: %w", err)
	}

	return nil
}
