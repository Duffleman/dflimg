package rpc

import (
	"net/http"
)

// HealthCheck is an RPC handler for checking the app is responding to requests
func (r *RPC) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(204)

	return
}
