package web_dl

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DownloadExistingFilenames bool
	Dirname                   string
	URIToFiles                []string
}

// parse args
// create http client
// create downloader
// create get request using http client + downloader

func Start() {
	config := parseArgs()
	client := newClient()
	filesToDownload := SetDownloaderArgs(config)
	downloadFiles(client, config, filesToDownload)
}

func parseArgs() *Config {
	config := &Config{}
	fmt.Printf("Args: %s\n", os.Args)
	// Create flags for settings for the downloader, http_client, etc...
	// Retry count, timeout, etc...

	// dirname flag
	const (
		defaultDirname = "../../media/"
		usageDirname   = "Path to the folder to download files to"
	)
	flag.StringVar(&config.Dirname, "dirname", defaultDirname, usageDirname)
	flag.StringVar(&config.Dirname, "d", defaultDirname, usageDirname+" (shorthand)")

	flag.BoolVar(&config.DownloadExistingFilenames, "E", false, `Set this flag to true, so files 
	with filenames already existing in the download directory are downloaded by appending a number to the filename`)

	flag.Parse()

	// process non-flag args
	if flag.NArg() == 0 {
		log.Fatalln("No URI's provided as non-flag arguments, please provide URI's.")
	}
	fmt.Printf("Number of URI's provided: %d\n", flag.NArg())

	config.URIToFiles = flag.Args()
	fmt.Printf("Tail: %s\n", config.URIToFiles)

	return config
}
