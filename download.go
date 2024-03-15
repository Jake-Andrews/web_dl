package web_dl

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type DownloadInfo struct {
	Filename    string
	Dirname     string
	ContentSize uint64
	URL         string
}

//0 to 18446744073709551616 = 18446744073.709552765 GB
//uint32 0 to 4294967295 roughly = 4.29 GB

// Option function type to customize DownloadInfo fields
type DownloadInfoOption func(*DownloadInfo)

// Option function to set the filename
func WithFilename(filename string) DownloadInfoOption {
	return func(di *DownloadInfo) {
		di.Filename = filename
	}
}

// Option function to set the dirname
func WithDirname(dirname string) DownloadInfoOption {
	return func(di *DownloadInfo) {
		di.Dirname = dirname
	}
}

// Option function to set the content size
func WithContentSize(contentSize uint64) DownloadInfoOption {
	return func(di *DownloadInfo) {
		di.ContentSize = contentSize
	}
}

// Option function to set the URL
func WithURL(url string) DownloadInfoOption {
	return func(di *DownloadInfo) {
		di.URL = url
	}
}

// Function to create a new DownloadInfo instance with provided options
func NewDownloadInfo(options ...DownloadInfoOption) *DownloadInfo {
	di := &DownloadInfo{}

	// Apply each option to the DownloadInfo instance
	for _, option := range options {
		option(di)
	}

	return di
}

func SetDownloaderArgs(args []string) *DownloadInfo {
	info := NewDownloadInfo(
		WithFilename("../../test/test_image.jpg"),
		WithDirname("../../test"),
		WithContentSize(10000000),
		WithURL("https://preview.redd.it/4dqyhtrsjrmc1.jpeg?auto=webp&s=093d3be09624e47cb9b90d011c50a20fede99e52"),
	)

	// Display the configured DownloadInfo instance
	fmt.Printf("Filename: %s\n", info.Filename)
	fmt.Printf("Dirname: %s\n", info.Dirname)
	fmt.Printf("Content Size: %d\n", info.ContentSize)
	fmt.Printf("URL: %s\n", info.URL)

	return info
}

func getFile(d *DownloadInfo) {
	resp, err := http.Get(d.URL)

	if err != nil {
		fmt.Printf("%q\nDownload url: %q\n", d.URL, err)
	}
	defer resp.Body.Close()

	createDirectory(d.Dirname)

	file, err := os.Create(d.Filename)
	if err != nil {
		log.Fatalf("%q\n", err)
	}
	defer file.Close()

	//byte_slice, err := io.ReadAll(resp.Body)
	//fmt.Printf("Body info-> len:%d, cap:%d", len(byte_slice), cap(byte_slice))
	// Use io.Copy to just dump the response body to the file. This supports huge files

	written, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("%q\nBytes written:%d", err, written)
	} else {
		fmt.Printf("Success writing to file, bytes written: %d\n", written)
	}
}
