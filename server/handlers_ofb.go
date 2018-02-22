package server

import "net/http"

func (s *Server) ofbProducts(w http.ResponseWriter, r *http.Request) {
	s.rdr.renderXML(w, r, tplProducts, nil)
}
