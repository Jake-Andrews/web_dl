package web_dl

//send args off
//create downloader
//run downloader on url

func Start(args []string) {
	DownloadInfo := SetDownloaderArgs(args)
	getFile(DownloadInfo)
}
