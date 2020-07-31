package lib

import (
	"encoding/json"
	"fmt"
)

type Iter interface {
	Next() bool
	Current() interface{}
	Err() error
}

func JsonMarshalIter(it Iter, fields string) {
	firstObject := true
	for it.Next() {
		recordMap, err := OnlyFields(fields, it.Current())
		if err != nil {
			fmt.Println(err)
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
		fmt.Println(it.Err())
	}
}

func JsonMarshal(i interface{}, fields string) {
	recordMap, err := OnlyFields(fields, i)
	if err != nil {
		fmt.Println(err)
	}
	prettyJSON, err := json.MarshalIndent(recordMap, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(prettyJSON))
}
