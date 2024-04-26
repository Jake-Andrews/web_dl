package web_dl

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// enables injecting later if needed
type loggingTransport struct {
	Transport http.RoundTripper
}

func (t *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	resp, roundTripErr := http.DefaultTransport.RoundTrip(r)

	requestBytes, errReq := httputil.DumpRequestOut(r, false)

	respBytes, errResp := httputil.DumpResponse(resp, false)

	if errReq == nil && errResp == nil {
		requestBytes = append(requestBytes, respBytes...)
		fmt.Printf("\n\n%s\n\n", requestBytes)
	}

	return resp, roundTripErr
}

func newClient() *http.Client {
	//DefaultClient, DefaultTransport, etc...
	//return &http.Client{}
	return &http.Client{
		Transport: &loggingTransport{},
	}

}
