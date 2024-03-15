package web_dl

import (
	"log"
	"os"
)

func createDirectory(dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
