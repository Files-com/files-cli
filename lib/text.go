package lib

import (
	"context"
	"fmt"
	"io"
)

func TextMarshalIter(_ context.Context, it Iter, _usePager bool, out io.Writer, filterIter FilterIter) error {
	for it.Next() {
		if it.Err() != nil {
			return it.Err()
		}

		current := it.Current()
		filter := true
		if filterIter != nil {
			var err error
			current, filter, err = filterIter(current)
			if err != nil {
				return err
			}
		}

		if filter {
			fmt.Fprintf(out, "%v\n", current)
		}
	}

	if it.Err() != nil {
		return it.Err()
	}

	return nil
}
