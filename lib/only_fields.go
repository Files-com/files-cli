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
	var subtractFields []string
	for _, key := range unparsedFields {
		if strings.HasPrefix(key, "-") {
			subtractFields = append(subtractFields, strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(key, "-"), "-", "_")))
		} else {
			fields = append(fields, strings.ToLower(strings.ReplaceAll(key, "-", "_")))
		}
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

	for _, key := range subtractFields {
		delete(returnMap, key)
		for i, orderedKey := range orderedKeys {
			if orderedKey == key {
				orderedKeys = remove(orderedKeys, i)
			}
		}
	}

	return returnMap, orderedKeys, nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func jsonTags(structure interface{}) []string {
	var tags []string
	return jsonTagsByFields(structs.New(structure).Fields(), tags)
}

func jsonTagsByFields(fields []*structs.Field, tags []string) []string {
	for _, field := range fields {
		if field.IsEmbedded() && field.Kind() == reflect.Struct {
			tags = jsonTagsByFields(field.Fields(), tags)
		}
		tag := strings.Split(field.Tag("json"), ",")[0]
		if tag == "-" || tag == "" {
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
