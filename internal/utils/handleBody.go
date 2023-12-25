package utils

import (
	"io"
	"log"
)

func HandleBody(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Fatal(err)
	}
}
