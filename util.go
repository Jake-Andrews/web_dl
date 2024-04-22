package web_dl

import (
	"log"
	"net/url"
	"os"
	"path"
)

func createDirectory(dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// extractFilename takes a URL string and returns the base filename component
func extractFilename(urlStr string, defaultUrlStr string) string {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		return defaultUrlStr
	}
	return path.Base(parsedUrl.Path)
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
