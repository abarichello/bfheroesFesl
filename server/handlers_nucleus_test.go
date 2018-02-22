package server

import (
	"net/http"
	"testing"
)

func TestServer_nucleusCheckUser(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+`/nucleus/check/user/`+nucleusID,
		nil,
	)
	setCommonTestHeaders(req)
	addTestGameClientHeaders(req)

	res, _ := ts.Client().Do(req)

	t.Log(res)
}
