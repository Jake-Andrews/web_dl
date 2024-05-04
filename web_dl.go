package web_dl

import (
	"flag"
	"fmt"
	"log"
)

type Config struct {
	DownloadExistingFilenames bool
	Dirname                   string
	URIToFiles                []string
	MaxConcurrentDownloads    int
	Timeout                   int
	NumConnections            int
}

func Start() {
	config := parseArgs()
	client := newClient(config.Timeout)
	filesToDownload := SetDownloaderArgs(config)
	downloadFiles(client, config, filesToDownload)
}

func parseArgs() *Config {
	config := &Config{}

	// dirname flag
	const (
		defaultDirname = "../../media/"
		usageDirname   = "Path to the folder to download files to"
	)
	flag.StringVar(&config.Dirname, "dirname", defaultDirname, usageDirname)
	flag.StringVar(&config.Dirname, "d", defaultDirname, usageDirname+" (shorthand)")

	flag.BoolVar(&config.DownloadExistingFilenames, "E", false, `Set this flag to true, so files 
	with filenames already existing in the download directory are downloaded by appending a number to the filename`)

	// maximum concurrent downloads flag.
	const (
		defaultMaxConcurrentDownloads = 5
		usageMaxConcurrentDownloads   = "Maximum number of concurrent downloads."
	)
	flag.IntVar(&config.MaxConcurrentDownloads, "maxconcurrentdownloads", defaultMaxConcurrentDownloads, usageMaxConcurrentDownloads)
	flag.IntVar(&config.MaxConcurrentDownloads, "M", defaultMaxConcurrentDownloads, usageMaxConcurrentDownloads+" (shorthand)")

	// timeout flag
	const (
		defaultTimeout = 0
		usageTimeout   = "Timeout for the HTTP client in seconds (0 = no timeout)"
	)
	flag.IntVar(&config.Timeout, "timeout", defaultTimeout, usageTimeout)
	flag.IntVar(&config.Timeout, "t", defaultTimeout, usageTimeout+" (shorthand)")

	// num-connections flag
	const (
		defaultNumConnections = 4
		usageNumConnections   = "Number of connections per download."
	)
	flag.IntVar(&config.NumConnections, "num-connections", defaultNumConnections, usageNumConnections)
	flag.IntVar(&config.NumConnections, "C", defaultNumConnections, usageNumConnections+" (shorthand)")

	flag.Parse()

	// process non-flag args
	if flag.NArg() == 0 {
		log.Fatalln("No URI's provided as non-flag arguments, please provide URI's.")
	}

	// remove non valid URI's from tail args
	count := 0
	for _, URI := range flag.Args() {
		if IsUrl(URI) {
			config.URIToFiles = append(config.URIToFiles, URI)
		} else {
			count += 1
		}
	}
	fmt.Printf("Number of URI's provided: %d\n", flag.NArg()-count)
	fmt.Printf("Tail URI's: %s\n\n", config.URIToFiles)

	return config
}
