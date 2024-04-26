package web_dl

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func createDirectory(dirPath string) {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

// takes a URL string and returns the base filename component if it can be found,
// otherwise returns defaultURLStr (2nd argument)
func extractFilename(urlStr string, dirname string) string {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Error parsing filename from: %s, not downloading\nError: %v\n", urlStr, err)
		return ""
	}
	// ensure the base path doesn't return useless garbage "." or "/"
	basePath := path.Base(parsedUrl.Path)
	if basePath == "/" || basePath == "" {
		log.Printf("Error parsing filename from: %s, got: %s, not downloading\n", urlStr, basePath)
		return ""
	}
	return dirname + basePath
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func getUniqueFilename(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path // filepath doesn't exist already, return it
	}

	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	base := filepath.Base(path)
	name := base[0 : len(base)-len(ext)]

	// append number until the filepath is unique
	for i := 1; ; i++ {
		newName := fmt.Sprintf("%s(%d)%s", name, i, ext)
		newPath := filepath.Join(dir, newName)

		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath // path is unique
		}
	}
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
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
