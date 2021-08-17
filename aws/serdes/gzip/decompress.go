package gzip

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
)

func Decompress(block []byte) []byte {
	buf := bytes.NewReader(block)

	r, err := gzip.NewReader(buf)
	if err != nil {
		log.Fatal(err)
	}

	output, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	return output
}
