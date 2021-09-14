package lib

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/fatih/structs"
)

func OnlyFields(commaFields string, structure interface{}) (map[string]interface{}, []string, error) {
	unparsedFields := strings.Split(commaFields, ",")
	jsonStructure, _ := json.MarshalIndent(structure, "", "    ")
	intermediateMap := make(map[string]interface{})
	returnMap := make(map[string]interface{})
	json.Unmarshal(jsonStructure, &intermediateMap)
	orderedKeys := jsonTags(structure)
	var fields []string
	for _, key := range unparsedFields {
		fields = append(fields, strings.ToLower(key))
	}

	if len(fields) > 0 && fields[0] != "" {
		orderedKeys = fields
		for _, key := range fields {
			_, ok := intermediateMap[key]
			if ok {
				returnMap[key] = intermediateMap[key]
			} else {
				if hasField(structure, key) {
					continue
				}

				return returnMap, orderedKeys, errors.New("field: `" + key + "` is not valid.")
			}
		}
	} else {
		for _, key := range orderedKeys {
			returnMap[key] = intermediateMap[key]
		}
	}

	return returnMap, orderedKeys, nil
}

func jsonTags(structure interface{}) []string {
	var tags []string
	for _, field := range structs.New(structure).Fields() {
		tags = append(tags, strings.Split(field.Tag("json"), ",")[0])
	}
	return tags
}

func hasField(structure interface{}, key string) bool {
	for _, field := range structs.New(structure).Fields() {
		if strings.Split(field.Tag("json"), ",")[0] == key {
			return true
		}
	}

	return false
}
