package rpc

import (
	"net/http"
)

// HeadResource gets a resource and handles the response for it
func (r *RPC) HeadResource(w http.ResponseWriter, req *http.Request) {
	fake := NewFake()

	r.HandleResource(fake, req)

	for key, value := range fake.headers {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	if fake.status >= 100 && fake.status <= 999 {
		w.WriteHeader(fake.status)
	}

	return
}

type FakeResponse struct {
	headers http.Header
	body    []byte
	status  int
}

func NewFake() *FakeResponse {
	return &FakeResponse{
		headers: make(http.Header),
	}
}

func (r *FakeResponse) Header() http.Header {
	return r.headers
}

func (r *FakeResponse) Write(body []byte) (int, error) {
	r.body = body
	return len(body), nil
}

func (r *FakeResponse) WriteHeader(status int) {
	r.status = status
}
