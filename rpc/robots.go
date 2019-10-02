package rpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

func (r *RPC) Robots(w http.ResponseWriter, req *http.Request) {
	fileContent, err := ioutil.ReadFile("resources/robots.txt")
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	modTime, err := time.Parse(time.RFC3339, "2019-10-02T12:00:00Z")
	if err != nil {
		r.handleError(w, req, err)
		return
	}

	reader := bytes.NewReader(fileContent)

	http.ServeContent(w, req, "robots.txt", modTime, reader)
	return
}
