package server

import (
	"net/http"

	"github.com/go-chi/chi"
)

type dtRelationship struct {
	ServerID string
}

func (s *Server) relationshipHandler(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplRelationship, &dtRelationship{chi.URLParam(r, "id")})
}

type dtSession struct {
	Token string
}

func (s *Server) sessionHandler(w http.ResponseWriter, r *http.Request) {
	if serverKey := getServerHeader(r); serverKey != "" {
		s.rdr.renderXML(w, r, tplSession, nil)
		return
	}

	userKey, err := r.Cookie("magma")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	s.rdr.renderXML(w, r, tplSessionNew, &dtSession{userKey.Value})
}

type dtHero struct {
	HeroID string
}

func (s *Server) entitlementsHandler(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplSessionNew, &dtHero{chi.URLParam(r, "heroID")})
}

func (s *Server) offersHandler(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplProducts, nil)
}

func (s *Server) walletsHandler(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplWallets, nil)
}

func (s *Server) store(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplStore, nil)
}
