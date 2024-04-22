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

type DownloadInfo struct {
	Dirname string
	dfiles  []DownloadFile
}

type DownloadFile struct {
	Filename    string
	ContentSize uint64
	URI         string
}

func SetDownloaderArgs(args map[string][]string) *DownloadInfo {
	URISlice := args["tail"]
	info := &DownloadInfo{
		Dirname: args["dirname"][0],
	}

	for i, URI := range URISlice {
		filename := "generic_fname" + strconv.Itoa(i)
		info.dfiles = append(info.dfiles, DownloadFile{Filename: filename, ContentSize: 0, URI: URI})
	}

	// Display the configured DownloadInfo instance
	fmt.Printf("Dirname: %s\n", info.Dirname)
	for i, file := range info.dfiles {
		fmt.Printf("Download File %d, Filename: %s\n", i, file.Filename)
		fmt.Printf("Download File %d, ContentSize: %d\n", i, file.ContentSize)
		fmt.Printf("Download File %d, URI: %s\n", i, file.URI)
	}

	return info
}

func downloadFiles(c *http.Client, d *DownloadInfo) {
	createDirectory(d.Dirname)
	for i, dfile := range d.dfiles {
		fmt.Printf("File# %d\n", i)
		getFile(c, &dfile, d)

	}
}

func getFile(c *http.Client, d *DownloadFile, dinfo *DownloadInfo) {
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
	d.Filename = dinfo.Dirname + d.Filename
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
