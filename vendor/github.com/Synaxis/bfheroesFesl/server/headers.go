package server

import (
	"net/http"
)

const (
	xServerKey = "X-SERVER-KEY"
)

func getServerHeader(r *http.Request) string {
	return r.Header.Get(xServerKey)
}
