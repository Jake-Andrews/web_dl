package web_dl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

// builds filename from URI. If skip=True, skip file (can't parse URI for filename).
func processURI(URI string, config *Config) (filename string, skip bool) {
	filename = extractFilename(URI, config.Dirname)
	if filename == "" {
		return "", true // skip if no filename could be extracted
	}
	if pathExists(filename) {
		if !config.DownloadExistingFilenames {
			fmt.Printf("File already exists and DownloadExistingFilenames flag set to false, not downloading: %s\n", filename)
			return "", true // skip download if file exists and DownloadExistingFilenames flag is set to false
		}
		filename = getUniqueFilename(filename) // file exists, flag=True, create unique filename and don't skip file
	}
	return filename, false
}

func logDownloadConfig(config *Config, files []DownloadFile) {
	fmt.Printf("Dirname: %s\n", config.Dirname)
	for i, file := range files {
		fmt.Printf("Download File %d, Filename: %s\n", i, file.Filename)
		fmt.Printf("Download File %d, ContentSize: %d\n", i, file.ContentSize)
		fmt.Printf("Download File %d, URI: %s\n", i, file.URI)
	}
}

func downloadFiles(c *http.Client, config *Config, filesToDownload []DownloadFile) {
	fmt.Println("downloadFiles")
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
		// check if a file is a duplicate
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

	//Create file w/ mode (0666)
	file, err := os.Create(d.Filename)
	if err != nil {
		log.Fatalf("%q\n", err)
	} else {
		fmt.Printf("Success creating file: %q\n", d.Filename)
	}
	defer file.Close()

	fmt.Printf("Response Body Len: %d\n", resp.ContentLength)
	d.ContentSize = uint64(resp.ContentLength)
	// ContentLength, -1 if length is unknown, unless Request.Method = HEAD, >= 0 means said # of bytes may be read from the body
	resp_len := resp.ContentLength
	written, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("%q\nBytes written:%d\n", err, written)
	} else if written != int64(resp_len) {
		log.Fatalf(`Error writing to file, bytes written: 
		%d Bytes, Number of bytes expected: %d Bytes\n`, written, resp_len)
	} else {
		fmt.Printf("Success writing to file, bytes written: %d Bytes\n", written)
		fmt.Printf("Success writing to file, MB written: %.2f MB\n", float64(written)/1000000)
	}
}
