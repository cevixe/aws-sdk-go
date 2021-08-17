package json

import (
	"encoding/json"
	"log"
)

func Unmarshall(buf []byte, o interface{}) {
	err := json.Unmarshal(buf, o)
	if err != nil {
		log.Fatal(err)
	}
}
