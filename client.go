package web_dl

import "net/http"

func newClient() *http.Client {
	//DefaultClient, DefaultTransport, etc...
	//return &http.Client{}
	return http.DefaultClient
}
