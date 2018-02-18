package server

import (
	"net/http"
	//"testing"

	"github.com/go-chi/chi"
)

type dtSession struct {
	Token string
}

// nucleusAuthToken authorizes the client by assigning a `magma` cookie.
func (s *Server) nucleusAuthToken(w http.ResponseWriter, r *http.Request) {
	if serverKey := getServerHeader(r); serverKey != "" {
		s.rdr.renderXML(w, r, tplSession, nil)
		return
	}

	userKey, err := r.Cookie("magma")
	//logger.Println(r.Cookies())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	s.rdr.renderXML(w, r, tplSessionNew, dtSession{userKey.Value})
}

// nucleusCheckUser is requested by both: game-server and game-client.
//
// See also TestServer_nucleusCheckUser as an example request
// made by game-client.
func (s *Server) nucleusCheckUser(w http.ResponseWriter, r *http.Request) {
	// userID := chi.URLParam("userID")
}

type dtHero struct {
	HeroID string
}

func (s *Server) nucleusEntitlements(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplEntitlements, &dtHero{chi.URLParam(r, "heroID")})
}

func (s *Server) walletsHandler(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplWallets, nil)
}

func (s *Server) store(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplStore, nil)
}

func (s *Server) ofbProducts(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplProducts, nil)
}

/*func TestServer_relationshipRosterServer(t *testing.T) {
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
*/

/*func TestServer_relationshipRosterNucleus(t *testing.T) {
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
*/

type dtRelationship struct {
	ServerID string
}

// relationshipHandler is used by game-server

// FIXME: This request is retried 5 times in 10sec interval.

// See also TestServer_relationshipRosterServer as an example request.
func (s *Server) relationshipRosterServer(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplRelationship, &dtRelationship{
		ServerID: chi.URLParam(r, "serverName"),
	})
}

// relationshipRosterNucleus is used by game-client.

// FIXME: This request is retried 5 times in 10sec interval.

// See also TestServer_relationshipRosterNucleus as an example request.
func (s *Server) relationshipRosterNucleus(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplRelationship, &dtRelationship{
		ServerID: chi.URLParam(r, "userID"),
	})
}
