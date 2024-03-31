package list

import (
	"context"
	"fmt"
	"strings"

	ujdscli "github.com/ashep/ujds/sdk/client"
	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
)

type List struct {
	cli *ujdscli.Client
	l   zerolog.Logger
}

type printer func(s ...any)

func New(cli *ujdscli.Client, l zerolog.Logger) *List {
	return &List{
		cli: cli,
		l:   l,
	}
}

func (l *List) List(ctx context.Context, names []string, format string, p printer) error {
	res, err := l.cli.I.List(ctx, connect.NewRequest(&indexproto.ListRequest{
		Filter: &indexproto.ListRequestFilter{
			Names: names,
		},
	}))

	if err != nil {
		return fmt.Errorf("ujds response: %w", err)
	}

	for _, idx := range res.Msg.GetIndices() {
		out := format
		out = strings.ReplaceAll(out, "{name}", idx.GetName())
		out = strings.ReplaceAll(out, "{title}", idx.GetTitle())

		p(out)
	}

	return nil
}
