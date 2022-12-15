package lib

import (
	"encoding/json"
)

func StructToMap(s interface{}) (map[string]interface{}, error) {
	var mapParams map[string]interface{}
	j, err := json.Marshal(s)
	if err != nil {
		return mapParams, err
	}
	err = json.Unmarshal(j, &mapParams)
	return mapParams, err
}
