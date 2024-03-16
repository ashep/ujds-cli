package jsontree

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrNotFound = errors.New("not found")

type Tree struct {
	a []any
	o map[string]any
}

func FromBytes(b []byte) (*Tree, error) {
	t := &Tree{}

	if err := json.Unmarshal(b, t); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return t, nil
}

func (t *Tree) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &t.o); err != nil {
		if err := json.Unmarshal(b, &t.a); err != nil {
			return errors.New("unmarshal failed")
		}
	}

	return nil
}

func (t *Tree) Keys() []string {
	if t.a == nil && t.o == nil {
		return nil
	}

	var res []string

	if t.a != nil {
		res = keys(t.a, "", nil)
	} else {
		res = keys(t.o, "", nil)
	}

	return res
}

func keys(node any, prefix string, out []string) []string {
	switch nodeT := node.(type) {
	case []any:
		for k, v := range nodeT {
			pp := ""
			ks := strconv.Itoa(k)
			if prefix != "" {
				pp += prefix + "." + ks
			} else {
				pp += ks
			}
			out = keys(v, pp, out)
		}
	case map[string]any:
		for k, v := range nodeT {
			pp := ""
			if prefix != "" {
				pp += prefix + "." + k
			} else {
				pp += k
			}
			out = keys(v, pp, out)
		}
	default:
		if nodeT != nil {
			out = append(out, prefix)
		}
	}

	return out
}

func (t *Tree) Get(q string) (any, error) {
	if t.a == nil && t.o == nil {
		return nil, ErrNotFound
	}

	var node any
	if t.a != nil {
		node = t.a
	} else {
		node = t.o
	}

	for _, k := range strings.Split(q, ".") {
		switch nodeT := node.(type) {
		case []any:
			ki, err := strconv.Atoi(k)
			if err != nil || ki >= len(nodeT) || ki < 0 {
				return nil, ErrNotFound
			}
			node = nodeT[ki]
		case map[string]any:
			var ok bool
			node, ok = nodeT[k]
			if !ok {
				return nil, ErrNotFound
			}
		default:
			return nil, ErrNotFound
		}
	}

	return node, nil
}
