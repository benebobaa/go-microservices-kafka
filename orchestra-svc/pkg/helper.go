package pkg

import "encoding/json"

func MergeJSON(dst interface{}, sources ...interface{}) error {
	mergedMap := make(map[string]interface{})

	for _, src := range sources {
		jsonData, err := json.Marshal(src)
		if err != nil {
			return err
		}

		tempMap := make(map[string]interface{})
		if err := json.Unmarshal(jsonData, &tempMap); err != nil {
			return err
		}

		for key, value := range tempMap {
			mergedMap[key] = value
		}
	}

	finalJSON, err := json.Marshal(mergedMap)
	if err != nil {
		return err
	}

	return json.Unmarshal(finalJSON, dst)
}
