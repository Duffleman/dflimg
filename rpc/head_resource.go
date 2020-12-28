package rpc

import (
	"net/http"

	"dflimg/lib/fakehttp"
)

// HeadResource gets a resource and handles the response for it
func (r *RPC) HeadResource(w http.ResponseWriter, req *http.Request) {
	fake := fakehttp.NewResponse()

	r.HandleResource(fake, req)

	for key, value := range fake.Headers {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}

	if fake.Status >= 100 && fake.Status <= 999 {
		w.WriteHeader(fake.Status)
	}

	return
}
