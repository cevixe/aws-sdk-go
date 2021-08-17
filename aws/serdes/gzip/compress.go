package gzip

import (
	"bytes"
	"compress/gzip"
	"log"
)

func Compress(block []byte) []byte {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	if _, err := zw.Write(block); err != nil {
		log.Fatal(err)
	}

	if err := zw.Close(); err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
