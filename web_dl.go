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
	MaxConcurrentDownloads    int
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
	// create flags for settings for the downloader, http_client, etc...
	// retry count, timeout, etc...

	// dirname flag
	const (
		defaultDirname = "../../media/"
		usageDirname   = "Path to the folder to download files to"
	)
	flag.StringVar(&config.Dirname, "dirname", defaultDirname, usageDirname)
	flag.StringVar(&config.Dirname, "d", defaultDirname, usageDirname+" (shorthand)")

	flag.BoolVar(&config.DownloadExistingFilenames, "E", false, `Set this flag to true, so files 
	with filenames already existing in the download directory are downloaded by appending a number to the filename`)

	// setup the maximum concurrent downloads flag.
	const (
		defaultMaxConcurrentDownloads = 5
		usageMaxConcurrentDownloads   = "Maximum number of concurrent downloads."
	)
	flag.IntVar(&config.MaxConcurrentDownloads, "maxconcurrentdownloads", defaultMaxConcurrentDownloads, usageMaxConcurrentDownloads)
	flag.IntVar(&config.MaxConcurrentDownloads, "M", defaultMaxConcurrentDownloads, usageMaxConcurrentDownloads+" (shorthand)")

	flag.Parse()

	// process non-flag args
	if flag.NArg() == 0 {
		log.Fatalln("No URI's provided as non-flag arguments, please provide URI's.")
	}

	// remove non valid URI's
	count := 0
	for _, URI := range flag.Args() {
		if IsUrl(URI) {
			config.URIToFiles = append(config.URIToFiles, URI)
		} else {
			count += 1
		}
	}
	fmt.Printf("Number of URI's provided: %d\n", flag.NArg()-count)
	fmt.Printf("Tail URI's: %s\n", config.URIToFiles)

	return config
}
