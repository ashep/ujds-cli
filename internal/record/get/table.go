package get

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/ashep/ujds-cli/pkg/jsontree"
)

// indent is prepended to every line of the data section output.
const indent = "    "

// dataTable renders a JSON document as an indented two-column key/value table, where keys are flattened dotted
// paths. If the data is not a JSON object or array, it is returned as is, indented the same way.
func dataTable(data string) string {
	tree, err := jsontree.FromBytes([]byte(data))
	if err != nil {
		return indentLines(data)
	}

	keys := tree.Keys()
	if len(keys) == 0 {
		return indentLines(data)
	}

	slices.SortFunc(keys, compareKeys)

	sb := &strings.Builder{}
	tw := tabwriter.NewWriter(sb, 0, 0, 2, ' ', 0)

	for _, k := range keys {
		v, err := tree.Get(k)
		if errors.Is(err, jsontree.ErrNotFound) {
			continue
		} else if err != nil {
			return indentLines(data)
		}

		_, _ = fmt.Fprintf(tw, "%s%s\t%s\n", indent, k, formatValue(v))
	}

	if err := tw.Flush(); err != nil {
		return indentLines(data)
	}

	return sb.String()
}

// indentLines prepends indent to every non-empty line of s.
func indentLines(s string) string {
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}

	return strings.Join(lines, "\n")
}

// compareKeys orders dotted key paths segment by segment, comparing numeric segments as numbers so that array
// indices don't end up in lexicographic order, e.g. "genres.2" before "genres.10".
func compareKeys(a, b string) int {
	aSeg, bSeg := strings.Split(a, "."), strings.Split(b, ".")

	for i := 0; i < len(aSeg) && i < len(bSeg); i++ {
		if aSeg[i] == bSeg[i] {
			continue
		}

		aNum, aErr := strconv.Atoi(aSeg[i])
		bNum, bErr := strconv.Atoi(bSeg[i])

		if aErr == nil && bErr == nil {
			return aNum - bNum
		}

		return strings.Compare(aSeg[i], bSeg[i])
	}

	return len(aSeg) - len(bSeg)
}

// formatValue renders a JSON value for a table cell, keeping numbers in their decimal form and making control
// characters visible so that they don't break the table layout.
func formatValue(v any) string {
	switch vT := v.(type) {
	case float64:
		return strconv.FormatFloat(vT, 'f', -1, 64)
	case string:
		if strings.ContainsAny(vT, "\n\r\t") {
			return strconv.Quote(vT)
		}
		return vT
	default:
		return fmt.Sprintf("%v", vT)
	}
}
