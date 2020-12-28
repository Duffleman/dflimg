package fakehttp

import (
	"net/http"
)

type FakeResponse struct {
	Headers http.Header
	Body    []byte
	Status  int
}

func NewResponse() *FakeResponse {
	return &FakeResponse{
		Headers: make(http.Header),
	}
}

func (r *FakeResponse) Header() http.Header {
	return r.Headers
}

func (r *FakeResponse) Write(body []byte) (int, error) {
	r.Body = body
	return len(body), nil
}

func (r *FakeResponse) WriteHeader(status int) {
	r.Status = status
}
