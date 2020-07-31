package lib

import (
	"encoding/json"
	"errors"
	"strings"
)

func OnlyFields(commaFields string, structure interface{}) (map[string]interface{}, error) {
	fields := strings.Split(commaFields, ",")
	jsonStructure, _ := json.MarshalIndent(structure, "", "    ")
	intermediateMap := make(map[string]interface{})
	returnMap := make(map[string]interface{})
	json.Unmarshal(jsonStructure, &intermediateMap)
	if len(fields) > 0 && fields[0] != "" {
		for _, key := range fields {
			if intermediateMap[key] != nil {
				returnMap[key] = intermediateMap[key]
			} else {
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
