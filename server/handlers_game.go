package server

import "net/http"

func (s *Server) store(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplStore, nil)
}
