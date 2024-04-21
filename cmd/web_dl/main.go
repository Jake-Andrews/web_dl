package main

import (
	"github.com/Jake-Andrews/web_dl"
)

func main() {
	//First arg must be a URL
	//if urlArg := os.Args[1]; govalidator.IsURL(urlArg) {
	//	log.Fatalf("Faulty url given: %s", urlArg)
	//}
	web_dl.Start(make([]string, 0))
}
