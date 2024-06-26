package export

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	ujdscli "github.com/ashep/ujds/sdk/client"
	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
)

const (
	formatCSV = "csv"
)

type Export struct {
	cli *ujdscli.Client
	l   zerolog.Logger
}

func New(cli *ujdscli.Client, l zerolog.Logger) *Export {
	return &Export{
		cli: cli,
		l:   l,
	}
}

func (e *Export) Export(ctx context.Context, idxPatterns []string, filename string, overwrite bool) error {
	di := strings.Index(filename, ".")
	if di < 0 || di == 0 && strings.Count(filename, ".") == 1 {
		return errors.New("output filename must have an extension")
	}

	ss := strings.Split(filename, ".")
	ext := ss[len(ss)-1]
	if strings.ToLower(ext) != formatCSV {
		return fmt.Errorf("unsupported format: %s", ext)
	}

	if _, err := os.Stat(filename); err == nil && !overwrite {
		return fmt.Errorf("file already exists: %s; use the '--overwrite' flag to replace it", filename)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	indices, err := e.getIndexList(ctx, idxPatterns)
	if err != nil {
		return fmt.Errorf("get index list: %w", err)
	}

	rawFilepath := filename + ".tmp"
	if err := e.fetchRecordsToFile(ctx, indices, rawFilepath); err != nil {
		return fmt.Errorf("get records: %w", err)
	}

	records, err := e.readRecordsFromFile(rawFilepath)
	if err != nil {
		return fmt.Errorf("read records from raw file: %w", err)
	}

	if err := e.recordsToCSV(records, filename); err != nil {
		return fmt.Errorf("build csv: %w", err)
	}

	if err := os.Remove(rawFilepath); err != nil {
		return fmt.Errorf("delete tmp data file: %w", err)
	}

	return nil
}

func (e *Export) getIndexList(ctx context.Context, patterns []string) ([]string, error) {
	res := make([]string, 0)

	cRes, err := e.cli.I.List(ctx, connect.NewRequest(&indexproto.ListRequest{
		Filter: &indexproto.ListRequestFilter{
			Names: patterns,
		},
	}))

	if err != nil {
		return nil, fmt.Errorf("ujds response: %w", err)
	}

	for _, cr := range cRes.Msg.GetIndices() {
		res = append(res, cr.GetName())
	}

	return res, nil
}
