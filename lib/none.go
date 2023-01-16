package lib

import (
	"context"
)

func NoneMarshalIter(ctx context.Context, it Iter) error {
	for it.Next() {
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
