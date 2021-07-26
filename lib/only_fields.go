package lib

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/fatih/structs"
)

func OnlyFields(commaFields string, structure interface{}) (map[string]interface{}, error) {
	fields := strings.Split(commaFields, ",")
	jsonStructure, _ := json.MarshalIndent(structure, "", "    ")
	intermediateMap := make(map[string]interface{})
	returnMap := make(map[string]interface{})
	json.Unmarshal(jsonStructure, &intermediateMap)
	if len(fields) > 0 && fields[0] != "" {
		for _, key := range fields {
			_, ok := intermediateMap[key]
			if ok {
				returnMap[key] = intermediateMap[key]
			} else {
				if hasField(structure, key) {
					continue
				}

				return returnMap, errors.New("field: `" + key + "` is not valid.")
			}
		}
	} else {
		for key, value := range intermediateMap {
			returnMap[key] = value
		}
	}

	return returnMap, nil
}

func hasField(structure interface{}, key string) bool {
	for _, field := range structs.New(structure).Fields() {
		if strings.Split(field.Tag("json"), ",")[0] == key {
			return true
		}
	}

	return false
}
