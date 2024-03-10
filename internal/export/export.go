package export

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	ujdscli "github.com/ashep/ujds/sdk/client"
	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
)

const (
	formatJSON = "json"
	formatCSV  = "csv"
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

func (e *Export) Export(ctx context.Context, idxPatterns []string, outFilepath, format string, overwrite bool) error {
	if format != formatCSV {
		return fmt.Errorf("unsupported format: %s", format)
	}

	outFilepath = strings.ReplaceAll(outFilepath, "{FORMAT}", format)
	if _, err := os.Stat(outFilepath); err == nil && !overwrite {
		return fmt.Errorf("file already exists: %s", outFilepath)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// indices, err := e.getIndexList(ctx, idxPatterns)
	// if err != nil {
	// 	return fmt.Errorf("get index list: %w", err)
	// }

	rawFilepath := outFilepath + ".tmp"
	// if err := e.dumpRecords(ctx, indices, rawFilepath); err != nil {
	// 	return fmt.Errorf("get records: %w", err)
	// }

	if err := e.rawToCSV(rawFilepath, outFilepath); err != nil {
		return fmt.Errorf("build csv: %w", err)
	}

	return nil
}

func (e *Export) getIndexList(ctx context.Context, patterns []string) ([]string, error) {
	rePatterns := make([]*regexp.Regexp, len(patterns))
	for i, pat := range patterns {
		pat = strings.ReplaceAll(pat, ".", "\\.")
		pat = strings.ReplaceAll(pat, "*", ".*")

		re, err := regexp.Compile("^" + pat + "$")
		if err != nil {
			return nil, fmt.Errorf("invalid index name pattern: %s; %w", err)
		}
		rePatterns[i] = re
	}

	res := make([]string, 0)

	cRes, err := e.cli.I.List(ctx, connect.NewRequest(&indexproto.ListRequest{}))
	if err != nil {
		return nil, fmt.Errorf("ujds: %w", err)
	}

	for _, cr := range cRes.Msg.GetIndices() {
		name := cr.GetName()

		for _, re := range rePatterns {
			if re.MatchString(name) {
				res = append(res, name)
				break
			}
		}
	}

	return res, nil
}

func (e *Export) dumpRecords(ctx context.Context, indices []string, filepath string) error {
	fd, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			e.l.Error().Err(err).Msg("failed to close file")
		}
	}()

	for _, idx := range indices {
		count := 0
		cursor := uint64(0)

		for {
			res, err := e.cli.R.Find(ctx, connect.NewRequest(&recordproto.FindRequest{
				Index:  idx,
				Cursor: cursor,
			}))

			if err != nil {
				return fmt.Errorf("ujds response: %w", err)
			}

			for _, rec := range res.Msg.GetRecords() {
				if _, err := fd.WriteString(rec.GetData() + "\n"); err != nil {
					return fmt.Errorf("write to %s: %w", filepath, err)
				}
				count++
			}

			newCursor := res.Msg.Cursor
			if newCursor == 0 {
				e.l.Debug().Int("count", count).Str("index", idx).Msg("records read")
				break
			}

			cursor = newCursor
		}
	}

	return nil
}

func (e *Export) rawToCSV(inFilename, outFilename string) error {
	inFD, err := os.OpenFile(inFilename, os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}
	defer func() {
		if err := inFD.Close(); err != nil {
			e.l.Error().Err(err).Msg("failed to close input file")
		}
	}()

	outFD, err := os.OpenFile(outFilename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("open output file: %w", err)
	}
	defer func() {
		if err := outFD.Close(); err != nil {
			e.l.Error().Err(err).Msg("failed to close output file")
		}
	}()

	inRdr := bufio.NewReader(inFD)

	for i := 1; ; i++ {
		str, err := inRdr.ReadString('\n')

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("read line %d: %w", i, err)
		}

		fmt.Println(str)
	}

	return nil
}
