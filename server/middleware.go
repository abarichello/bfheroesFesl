package server

import (
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

func (s *Server) logRequestMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func logRequest(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		logrus.Debugf("Cannot call DumpRequest: %v", err)
		return
	}

	logrus.
		WithFields(logrus.Fields{
			"req": string(dump),
		}).
		Debugf("Dump %s", r.URL.Path)
}
