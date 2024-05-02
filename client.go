package web_dl

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

// enables injecting later if needed
type loggingTransport struct {
	Transport http.RoundTripper
}

func (t *loggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	resp, roundTripErr := t.Transport.RoundTrip(r)
	requestBytes, errReq := httputil.DumpRequestOut(r, false)
	respBytes, errResp := httputil.DumpResponse(resp, false)

	if errReq == nil && errResp == nil {
		requestBytes = append(requestBytes, respBytes...)
		fmt.Printf("\n\n%s\n\n", requestBytes)
	}

	return resp, roundTripErr
}

func newClient(timeout int) *http.Client {
	return &http.Client{
		Transport: &loggingTransport{Transport: http.DefaultTransport},
		Timeout:   time.Duration(timeout) * time.Second,
	}
}
