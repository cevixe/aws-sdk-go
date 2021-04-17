package util

import "encoding/json"

func UnmarshalJsonString(data string, v interface{}) {
	buffer := []byte(data)
	err := json.Unmarshal(buffer, v)
	if err != nil {
		panic(err)
	}
}

func MarshalJsonString(v interface{}) string {
	buffer, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(buffer)
}

func UnmarshalJson(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		panic(err)
	}
}

func MarshalJson(v interface{}) []byte {
	buffer, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return buffer
}
