package jsn

import "fmt"

func GetObj(path string) (map[string]interface{}, error) {
	j, err := Get(path)
	if err != nil {
		return nil, err
	}

	jo, ok := j.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%v is not a json object", j)
	}

	return jo, nil
}
