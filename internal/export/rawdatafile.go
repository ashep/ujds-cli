package export

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

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
				// id, created_at, updated_at, touched_at, data
				out := fmt.Sprintf("%s\x1e%d\x1e%d\x1e%d\x1e%s",
					rec.GetId(), rec.GetCreatedAt(), rec.GetCreatedAt(), rec.GetTouchedAt(), rec.GetData())

				if _, err := fd.WriteString(out + "\n"); err != nil {
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

		bs := bytes.Split(b, []byte{0x1e})
		if len(bs) != 5 {
			return nil, fmt.Errorf("wrong format at line: %d", i)
		}

		tree, err := jsontree.FromBytes(bs[4])
		if err != nil {
			return nil, fmt.Errorf("parse data at line: %d", i)
		}

		if err := tree.Set("@id", string(bs[0])); err != nil {
			return nil, fmt.Errorf("set metadata at line %d: %w", i, err)
		}

		if err := tree.Set("@created_at", strTimestampToDatetime(bs[1])); err != nil {
			return nil, fmt.Errorf("set metadata at line %d: %w", i, err)
		}

		if err := tree.Set("@updated_at", strTimestampToDatetime(bs[2])); err != nil {
			return nil, fmt.Errorf("set metadata at line %d: %w", i, err)
		}

		if err := tree.Set("@touched_at", strTimestampToDatetime(bs[3])); err != nil {
			return nil, fmt.Errorf("set metadata at line %d: %w", i, err)
		}

		res = append(res, tree)
	}

	return res, nil
}

func strTimestampToDatetime(ts []byte) string {
	tsi, err := strconv.Atoi(string(ts))
	if err != nil {
		tsi = 0
	}

	return time.Unix(int64(tsi), 0).Format(time.DateTime)
}
