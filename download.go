package web_dl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

//0 to 18446744073709551616 = 18446744073.709552765 GB
//uint32 0 to 4294967295 roughly = 4.29 GB

type DownloadFile struct {
	Filename    string
	ContentSize uint64
	URI         string
}

func SetDownloaderArgs(config *Config) []DownloadFile {
	filesToDownload := []DownloadFile{}
	for i, URI := range config.URIToFiles {
		filename := "generic_fname" + strconv.Itoa(i) //set a fname incase a fname could not be generated later
		filesToDownload = append(filesToDownload, DownloadFile{Filename: filename, ContentSize: 0, URI: URI})
	}

	// Display the configured DownloadInfo instance
	fmt.Printf("Dirname: %s\n", config.Dirname)
	for i, file := range filesToDownload {
		fmt.Printf("Download File %d, Filename: %s\n", i, file.Filename)
		fmt.Printf("Download File %d, ContentSize: %d\n", i, file.ContentSize)
		fmt.Printf("Download File %d, URI: %s\n", i, file.URI)
	}

	return filesToDownload
}

func downloadFiles(c *http.Client, config *Config, filesToDownload []DownloadFile) {
	createDirectory(config.Dirname)
	for i, dfile := range filesToDownload {
		fmt.Printf("File# %d\n", i)
		// build filename
		dfile.Filename = extractFilename(dfile.URI, dfile.Filename)
		dfile.Filename = config.Dirname + dfile.Filename

		// if the file exists and flags set to false, don't download file
		if pathExists(dfile.Filename) && !config.DownloadExistingFilenames {
			fmt.Printf("File already exists, not downloading: %s", dfile.Filename)
			continue
		}
		getFile(c, config, &dfile)
	}
}

func getFile(c *http.Client, config *Config, d *DownloadFile) {
	//*https://pkg.go.dev/net/http#Get GET url
	resp, err := c.Get(d.URI)
	if err != nil {
		log.Fatalf("%q\nDownload url: %q\n", err, d.URI)
	} else {
		fmt.Printf("Success downloading url: %q\n", d.URI)
	}
	//resp.Body ReadCloser interface, which contains Reader and Closer interfaces
	defer resp.Body.Close()

	//try to parse the filename & ext from URI, if this fails, use a generic filename
	d.Filename = extractFilename(d.URI, d.Filename)
	d.Filename = config.Dirname + d.Filename
	d.Filename = getUniqueFilename(d.Filename)

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
