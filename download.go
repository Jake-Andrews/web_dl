package web_dl

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

//0 to 18446744073709551616 = 18446744073.709552765 GB
//uint32 0 to 4294967295 roughly = 4.29 GB

type DownloadFile struct {
	Filename    string
	ContentSize uint64
	URI         string
}

// prepares a slice of DownloadFile based on the provided Config.
func SetDownloaderArgs(config *Config) []DownloadFile {
	var filesToDownload []DownloadFile
	for _, URI := range config.URIToFiles {
		filename, skip := processURI(URI, config)
		if skip {
			continue
		}
		filesToDownload = append(filesToDownload, DownloadFile{Filename: filename, ContentSize: 0, URI: URI})
	}

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
	var wg sync.WaitGroup

	// worker goroutines
	for i := 0; i < config.MaxConcurrentDownloads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for dfile := range jobs {
				fmt.Printf("Inside of goroutine: %s\n", dfile.URI)
				getFile(c, &dfile)
			}
		}()
	}

	// distribute jobs to workers
	for _, dfile := range filesToDownload {
		jobs <- dfile
	}
	close(jobs) // close jobs channel to signal no more jobs are coming

	wg.Wait()
}

func getFile(c *http.Client, d *DownloadFile) {
	//*https://pkg.go.dev/net/http#Get GET url
	fmt.Printf("Attempting to download: %s\n", d.URI)
	resp, err := c.Get(d.URI)
	if err != nil {
		log.Printf("Failed to download URL %q: %v\n", d.URI, err)
		return
	} else {
		fmt.Printf("Success downloading url: %q\n", d.URI)
	}
	//resp.Body ReadCloser interface, which contains Reader and Closer interfaces
	defer resp.Body.Close()

	//Create file w/ mode (0666) default, and transfer body to it
	if fileErr := withCreate(d.Filename, resp.Body); fileErr != nil {
		return
	}
}
