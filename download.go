package web_dl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type DownloadFile struct {
	Filename    string
	ContentSize int
	URI         string
}

// prepares a slice of DownloadFile based on the provided Config and creates directory.
func SetDownloaderArgs(config *Config) []DownloadFile {
	var filesToDownload []DownloadFile
	for _, URI := range config.URIToFiles {
		filename, skip := processURI(URI, config)
		if skip {
			continue
		}
		filesToDownload = append(filesToDownload, DownloadFile{Filename: filename, ContentSize: 0, URI: URI})
	}
	createDirectory(config.Dirname)
	logDownloadConfig(config, filesToDownload)
	return filesToDownload
}

func logDownloadConfig(config *Config, files []DownloadFile) {
	fmt.Printf("Dirname: %s\n", config.Dirname)
	for i, file := range files {
		fmt.Printf("Download File %d, Filename: %s\n", i, file.Filename)
		fmt.Printf("Download File %d, ContentSize: %d\n", i, file.ContentSize)
		fmt.Printf("Download File %d, URI: %s\n", i, file.URI)
	}
	fmt.Println()
}

func downloadFiles(c *http.Client, config *Config, filesToDownload []DownloadFile) {
	createDirectory(config.Dirname)

	jobs := make(chan DownloadFile, len(filesToDownload))
	results := make(chan error, len(filesToDownload))
	var wg sync.WaitGroup

	for i := 0; i < config.MaxConcurrentDownloads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for dfile := range jobs {
				log.Printf("Inside of goroutine: %s\n", dfile.URI)

				//HEAD, if Accept-Ranges do range requests, otherwise getFile
				resp, err := getHead(c, dfile.URI)
				if err != nil || resp == nil {
					results <- fmt.Errorf("failed to get head request for %s: %v", dfile.URI, err)
					continue
				}
				//range requests
				if strings.EqualFold(resp.Header.Get("Accept-Ranges"), "bytes") && resp.Header.Get("Content-Length") != "" {
					contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
					if err != nil {
						results <- fmt.Errorf("error parsing content length for %s: %v", dfile.URI, err)
						continue
					}
					err = buildRangeRequests(c, contentLength, &dfile, config.NumConnections)
					if err != nil {
						results <- err
					} else {
						err = joinFiles(&dfile, config.NumConnections)
						results <- err
					}
				} else { //otherwise one get request
					results <- getFile(c, &dfile)
				}
			}
		}()
	}

	//distribute jobs to workers
	for _, dfile := range filesToDownload {
		jobs <- dfile
	}
	close(jobs)
	wg.Wait()

	//handle results
	close(results)
	for err := range results {
		if err != nil {
			log.Println(err)
		}
	}
}

func buildRangeRequests(c *http.Client, contentLength int, d *DownloadFile, parts int) error {
	partSize := contentLength / parts
	var wg sync.WaitGroup
	errors := make(chan error, parts)

	for i := 0; i < parts; i++ {
		startSize := i * partSize
		endSize := (i+1)*partSize - 1
		if i == parts-1 {
			//ensure the last part goes up to the end of the content length
			//it should...still be valid even if this wasn't done, as long as start is in bounds
			//server should assume endSize > contentlength implies startSize to endSize
			endSize = contentLength - 1
		}

		req, err := http.NewRequest("GET", d.URI, nil)
		if err != nil {
			return err
		}
		rangeHeader := fmt.Sprintf("bytes=%d-%d", startSize, endSize)
		req.Header.Add("Range", rangeHeader)

		wg.Add(1)
		go func(req *http.Request, partNum int) {
			defer wg.Done()
			err := getRangeRequest(c, req, d, partNum)
			errors <- err
		}(req, i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

func getRangeRequest(c *http.Client, req *http.Request, d *DownloadFile, partNum int) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	ext := filepath.Ext(d.Filename)
	name := strings.TrimSuffix(d.Filename, ext)
	newFileName := fmt.Sprintf("%s.part%d%s", name, partNum, ext)

	//create and write to the file
	return withCreate(newFileName, resp.Body)
}

func joinFiles(d *DownloadFile, parts int) error {
	ext := filepath.Ext(d.Filename)
	name := strings.TrimSuffix(d.Filename, ext)
	finalFileName := fmt.Sprintf("%s%s", name, ext)

	finalFile, err := os.Create(finalFileName)
	if err != nil {
		return err
	}
	defer finalFile.Close()

	for i := 0; i < parts; i++ {
		partFileName := fmt.Sprintf("%s.part%d%s", name, i, ext)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(finalFile, partFile)
		partFile.Close()
		os.Remove(partFileName) //cleanup part file regardless of copy success
		if err != nil {
			return err
		}
	}
	return nil
}

func getFile(c *http.Client, d *DownloadFile) error {
	//*https://pkg.go.dev/net/http#Get GET url
	fmt.Printf("Attempting to download: %s\n", d.URI)
	resp, err := c.Get(d.URI)
	if err != nil {
		log.Printf("Failed to download URL %q: %v\n", d.URI, err)
		return err
	} else {
		fmt.Printf("Success downloading url: %q\n", d.URI)
	}
	//resp.Body ReadCloser interface, which contains Reader and Closer interfaces
	defer resp.Body.Close()

	//create file w/ mode (0666) default, and transfer body to it
	if fileErr := withCreate(d.Filename, resp.Body); fileErr != nil {
		return err
	}
	return nil
}

func getHead(c *http.Client, url string) (*http.Response, error) {
	resp, err := c.Head(url)
	if err != nil {
		log.Printf("%v\ngetHead for: %s\n", err, url)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("%d\ngetHead for: %s\n", resp.StatusCode, url)
		return nil, err
	}

	return resp, nil
}
