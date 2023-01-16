package lib

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strings"

	"github.com/fatih/structs"
)

func OnlyFields(unparsedFields []string, structure interface{}) (map[string]interface{}, []string, error) {
	if reflect.ValueOf(structure).Kind() == reflect.Map {
		m, _ := structure.(map[string]interface{})
		var ordered []string
		for k := range m {
			ordered = append(ordered, k)
		}
		sort.Slice(ordered, func(i, j int) bool {
			return ordered[i] < ordered[j]
		})
		return structure.(map[string]interface{}), ordered, nil
	}
	jsonStructure, err := json.MarshalIndent(structure, "", "    ")
	if err != nil {
		return structure.(map[string]interface{}), []string{}, err
	}
	intermediateMap := make(map[string]interface{})
	returnMap := make(map[string]interface{})
	json.Unmarshal(jsonStructure, &intermediateMap)
	orderedKeys := jsonTags(structure)
	var fields []string
	for _, key := range unparsedFields {
		fields = append(fields, strings.ToLower(strings.Replace(key, "-", "_", -1)))
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
		tag := strings.Split(field.Tag("json"), ",")[0]
		if tag == "-" {
			continue
		}
		tags = append(tags, tag)
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
