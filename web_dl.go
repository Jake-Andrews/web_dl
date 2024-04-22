package web_dl

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// parse args
// create http client
// create downloader
// create get request using http client + downloader

func Start() {
	args := parseArgs()
	fmt.Printf("Args: %q\n", args)
	client := newClient()
	DownloadInfo := SetDownloaderArgs(args)
	downloadFiles(client, DownloadInfo)
	//getFile(client, DownloadInfo)

}

func parseArgs() map[string][]string {
	fmt.Printf("Args: %s\n", os.Args)
	// Create flags for settings for the downloader, http_client, etc...
	// Retry count, timeout, etc...

	// dirname flag
	var dirname string
	const (
		defaultDirname = "../../media/"
		usageDirname   = "Path to the folder to download files to"
	)
	flag.StringVar(&dirname, "dirname", defaultDirname, usageDirname)
	flag.StringVar(&dirname, "d", defaultDirname, usageDirname+" (shorthand)")
	var dirnameSlice = []string{dirname}
	flag.Parse()

	args := make(map[string][]string)
	args["dirname"] = dirnameSlice
	fmt.Printf("Dirname: %s\n", args["dirname"])

	// process non-flag args
	if flag.NArg() == 0 {
		log.Fatalln("No URI's provided as non-flag arguments")
	}
	fmt.Printf("Number of URI's provided: %d\n", flag.NArg())

	args["tail"] = flag.Args()
	fmt.Printf("Tail: %s\n", flag.Args())

	return args
}
