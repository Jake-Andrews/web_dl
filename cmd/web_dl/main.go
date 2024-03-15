package main

import (
	"github.com/Jake-Andrews/web_dl"
)

/*
Future options:

	History
	Validate URL given by user (github.com/asaskevich/govalidator govalidator.IsURL(url))
	Duplicate check
	Auth
	Options (proxy, timeouts, etc...)
	Dockerize
	Github workflow
	Display progress

Current Goals:

	CLI
	Requests
	Save data
	Logging
	Concurency

	Proper Tests
	HEAD request and deal with range response
	More specificity when dealing with servers responses

Current To Do Order:

	cmd/web_dl/main.go
	client.go
	requests.go
	download.go
	util.go

	1. Pass URL as arg
	2. Perform GET request on arg and save the media
	4. Logging
	5. Concurency
*/

func main() {
	//First arg must be a URL
	//if urlArg := os.Args[1]; govalidator.IsURL(urlArg) {
	//	log.Fatalf("Faulty url given: %s", urlArg)
	//}
	web_dl.Start(make([]string, 0))
}
