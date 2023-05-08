package lib

import (
	"context"
	"io"
)

func TableMarshalV2(_ string, result interface{}, fields []string, _ bool, out io.Writer, _ string) error {
	model := &tableModel{fields: fields, out: out}
	model.Init()
	model.tableLoader = (&tableResource{}).Init(context.Background(), result, "")
	program, err := model.LoadAndBuild(context.Background())
	if err != nil {
		return err
	}
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}

func TableMarshalV2Iter(parentCtx context.Context, _ string, it Iter, fields []string, _ bool, out io.Writer, filterIter FilterIter) error {
	model := &tableModel{fields: fields, out: out}
	model.Init()
	model.tableLoader, _ = (&tableLoaderIter{}).Init(parentCtx, "", func(ctx context.Context) (Iter, error) {
		return it, nil
	})
	program, err := model.LoadAndBuild(parentCtx)
	if err != nil {
		return err
	}
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
