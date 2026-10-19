package output

import (
	"fmt"
	"io"

	"github.com/grafana/gcx/internal/format"
)

// TableAdapter extracts rows only for human output, leaving the complete payload
// available to JSON/YAML codecs. P is the command's payload and T is one table row.
type TableAdapter[P, T any] struct {
	Table Table[T]
	Rows  func(P) ([]T, error)
}

func (a TableAdapter[P, T]) Codec(name string) format.Codec { //nolint:ireturn // codec registration requires the interface
	return &adaptedTableCodec[P, T]{Codec: a.Table.Codec(name), rows: a.Rows}
}

type adaptedTableCodec[P, T any] struct {
	format.Codec

	rows func(P) ([]T, error)
}

func (c *adaptedTableCodec[P, T]) Encode(w io.Writer, v any) error {
	payload, ok := v.(P)
	if !ok {
		return fmt.Errorf("invalid data type for %s codec: expected %T, got %T", c.Format(), payload, v)
	}
	rows, err := c.rows(payload)
	if err != nil {
		return err
	}
	return c.Codec.Encode(w, rows)
}
