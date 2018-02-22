package server

import (
	"net/http"
	"testing"
)

func TestServer_relationshipRosterServer(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+`/relationships/roster/server:`+serverName+`/bvip/1,3`,
		nil,
	)
	setCommonTestHeaders(req)
	addTestGameServerHeaders(req)

	res, _ := ts.Client().Do(req)

	t.Log(res)
}

func TestServer_relationshipRosterNucleus(t *testing.T) {
	ts := startTestServer()
	defer ts.Close()

	req, _ := http.NewRequest(
		http.MethodGet,
		ts.URL+`/relationships/roster/nucleus:`+nucleusID,
		nil,
	)
	setCommonTestHeaders(req)
	addTestGameClientHeaders(req)

	res, _ := ts.Client().Do(req)

	t.Log(res)
}
