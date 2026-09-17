package lib

import (
	"context"
	"fmt"
	"io"

	"github.com/Files-com/files-cli/lib/errcheck"
)

func TextMarshalIter(_ context.Context, it Iter, _usePager bool, out io.Writer, filterIter FilterIter) error {
	escape := terminalEscaper(out)
	for it.Next() {
		if it.Err() != nil {
			return it.Err()
		}

		current := it.Current()

		if err := errcheck.CheckEmbeddedErrors(current); err != nil {
			return err
		}

		filter := true
		if filterIter != nil {
			var err error
			current, filter, err = filterIter(current)
			if err != nil {
				return err
			}
		}

		if filter {
			fmt.Fprintf(out, "%s\n", escape(fmt.Sprintf("%v", current)))
		}
	}

	if it.Err() != nil {
		return it.Err()
	}

	return nil
}
