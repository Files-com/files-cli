package lib

import (
	"context"

	"github.com/Files-com/files-cli/lib/errcheck"
)

func NoneMarshalIter(ctx context.Context, it Iter) error {
	for it.Next() {
		current := it.Current()
		if err := errcheck.CheckEmbeddedErrors(current); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return nil
		}
		if it.Err() != nil {
			return it.Err()
		}
	}

	if it.Err() != nil {
		return it.Err()
	}

	return nil
}
