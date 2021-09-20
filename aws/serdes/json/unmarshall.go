package json

import (
	"encoding/json"
	"log"
)

func Unmarshall(buf []byte, o interface{}) {
	if err := json.Unmarshal(buf, o); err != nil {
		log.Fatal(err)
	}
}
