package web_dl

//send args off
//create downloader
//run downloader on url

func Start(args []string) {
	client := newClient()
	DownloadInfo := SetDownloaderArgs(args)
	getFile(client, DownloadInfo)
}
