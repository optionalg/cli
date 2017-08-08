package uaa

import "net/http"

// Response represents an UAA response object.
type Response struct {
	// Result represents the resource entity type that is expected in the
	// response JSON.
	Result interface{}

	// RawResponse represents the response body.
	RawResponse []byte

	// HTTPResponse represents the HTTP response object.
	HTTPResponse *http.Response

	// ResourceLocationURL represents the Location header value
	ResourceLocationURL string
}

func (r *Response) reset() {
	r.RawResponse = []byte{}
	r.HTTPResponse = nil
}
