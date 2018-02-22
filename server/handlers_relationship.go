package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

type dtRelationship struct {
	ServerID string
}

// relationshipHandler is used by game-server
//
// FIXME: This request is retried 5 times in 10sec interval.
//
// See also TestServer_relationshipRosterServer as an example request.
func (s *Server) relationshipRosterServer(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplRelationship, &dtRelationship{
		ServerID: chi.URLParam(r, "serverName"),
	})
}

// relationshipRosterNucleus is used by game-client.
//
// FIXME: This request is retried 5 times in 10sec interval.
//
// See also TestServer_relationshipRosterNucleus as an example request.
func (s *Server) relationshipRosterNucleus(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplRelationship, &dtRelationship{
		ServerID: chi.URLParam(r, "userID"),
	})
}
