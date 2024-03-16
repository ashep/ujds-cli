package export

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ashep/ujds-cli/pkg/jsontree"
)

func (e *Export) recordsToCSV(records []*jsontree.Tree, outFilename string) error {
	fd, err := os.OpenFile(outFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		if err := fd.Close(); err != nil {
			e.l.Error().Err(err).Msg("failed to close output file")
		}
	}()

	keys := e.allKeys(records)

	if _, err := fd.Write(e.csvHeader(keys)); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	for l, item := range records {
		values := make([]string, 0)
		for _, k := range keys {
			v, err := item.Get(k)
			if errors.Is(err, jsontree.ErrNotFound) {
				values = append(values, "")
			} else if err != nil {
				e.l.Warn().Err(err).Msgf("line %d encode failed", l)
			} else {
				values = append(values, fmt.Sprintf("%v", v))
			}
		}
		if _, err := fd.Write(e.csvLine(values)); err != nil {
			e.l.Warn().Err(err).Msgf("line %d write failed", l)
		}
	}

	return nil
}

// allKeys returns all possible keys found in records; each key can be found at least at one record.
func (e *Export) allKeys(in []*jsontree.Tree) []string {
	keys := make(map[string]struct{})

	for _, t := range in {
		for _, k := range t.Keys() {
			keys[k] = struct{}{}
		}
	}

	res := make([]string, 0, len(keys))
	for k := range keys {
		res = append(res, k)
	}

	// Sort alphabetically
	slices.SortStableFunc(res, func(a, b string) int {
		for i := 0; ; i++ {
			if i == len(a) {
				return -1
			}

			if i == len(b) {
				return 1
			}

			if a[i] < b[i] {
				return -1
			}

			if a[i] > b[i] {
				return 1
			}
		}
	})

	// Sort by string length
	slices.SortStableFunc(res, func(a, b string) int {
		return len(a) - len(b)
	})

	// Sort by delimiter count
	slices.SortStableFunc(res, func(a, b string) int {
		return strings.Count(a, ".") - strings.Count(b, ".")
	})

	return res
}

func (e *Export) csvHeader(items []string) []byte {
	itemCopy := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.ReplaceAll(item, " ", "_")
		item = strings.ReplaceAll(item, " ", "/")
		itemCopy = append(itemCopy, item)
	}

	return []byte(strings.Join(itemCopy, ",") + "\n")
}

func (e *Export) csvLine(items []string) []byte {
	itemCopy := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.ReplaceAll(item, `"`, `""`)
		itemCopy = append(itemCopy, `"`+item+`"`)
	}

	return []byte(strings.Join(itemCopy, ",") + "\n")
}
