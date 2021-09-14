package lib

import (
	"encoding/json"
	"fmt"
)

type Iter interface {
	Next() bool
	Current() interface{}
	Err() error
	EOFPage() bool
}

func JsonMarshalIter(it Iter, fields string) error {
	firstObject := true
	for it.Next() {
		recordMap, _, err := OnlyFields(fields, it.Current())
		if err != nil {
			return err
		}
		prettyJSON, err := json.MarshalIndent(recordMap, "", "    ")
		if err != nil {
			panic(err)
		}
		if firstObject {
			fmt.Printf("[%s", string(prettyJSON))
		} else {
			fmt.Printf(",\n%s", string(prettyJSON))
		}

		firstObject = false
	}
	if firstObject {
		fmt.Print("[")
	}
	fmt.Println("]")
	if it.Err() != nil {
		return it.Err()
	}
	return nil
}

func JsonMarshal(i interface{}, fields string) error {
	recordMap, _, err := OnlyFields(fields, i)
	if err != nil {
		return err
	}
	prettyJSON, err := json.MarshalIndent(recordMap, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(prettyJSON))
	return err
}
