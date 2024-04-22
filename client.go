package web_dl

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

type loggingTransport struct{}

func newClient() *http.Client {
	//DefaultClient, DefaultTransport, etc...
	//return &http.Client{}
	return &http.Client{
		Transport: &loggingTransport{},
	}

}

// https://www.jvt.me/posts/2023/03/11/go-debug-http/
// prints to the http.Client on every request
func (s *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(r, false)

	resp, err := http.DefaultTransport.RoundTrip(r)
	// err is returned after dumping the response

	respBytes, _ := httputil.DumpResponse(resp, false)
	bytes = append(bytes, respBytes...)

	fmt.Printf("\n\n%s\n\n", bytes)

	return resp, err
}
