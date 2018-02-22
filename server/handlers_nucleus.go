package server

import (
	"net/http"

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
	logger.Println(r.Cookies())
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
