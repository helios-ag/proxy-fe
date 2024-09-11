package serializer

import "encoding/json"

func SerializeToString(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func DeserializeFromString(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}
