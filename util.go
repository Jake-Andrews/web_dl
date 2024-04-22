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

// reads txt file and returns urls
/* func readFile(filePath string) (URLslice []string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		URLslice = append(URLslice, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return URLslice
} */
