package export

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/bufbuild/connect-go"

	"github.com/ashep/ujds-cli/pkg/jsontree"
)

func (e *Export) fetchRecordsToFile(ctx context.Context, indices []string, filename string) error {
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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
					return fmt.Errorf("write to %s: %w", filename, err)
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

func (e *Export) readRecordsFromFile(filename string) ([]*jsontree.Tree, error) {
	res := make([]*jsontree.Tree, 0)

	fd, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			e.l.Error().Err(err).Msg("close file failed")
		}
	}()

	rdr := bufio.NewReader(fd)
	for i := 1; ; i++ {
		b, err := rdr.ReadBytes('\n')

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, fmt.Errorf("read line %d: %w", i, err)
		}

		tree, err := jsontree.FromBytes(b)
		if err != nil {
			e.l.Warn().Err(err).Msgf("parse data at line %d", i)
			continue
		}

		res = append(res, tree)
	}

	return res, nil
}
