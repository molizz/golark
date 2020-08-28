package api

import (
	"net/http"
)

type BearerToken struct {
	Base http.RoundTripper

	Token string // Bearer token
}

func (t *BearerToken) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Header.Get("Authorization") != "" {
		return t.base().RoundTrip(r)
	}

	r2 := cloneRequest(r)
	r2.Header.Set("Authorization", "Bearer "+t.Token)
	return t.base().RoundTrip(r2)
}

func (t *BearerToken) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func cloneRequest(r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
