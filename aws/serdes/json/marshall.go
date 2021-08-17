package json

import (
	"encoding/json"
	"log"
)

func Marshall(o interface{}) []byte {
	buf, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}

	return buf
}
