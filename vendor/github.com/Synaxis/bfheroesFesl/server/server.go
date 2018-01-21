package server

import (
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	stdlog "log"

	"github.com/Synaxis/bfheroesFesl/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// Server (also called as magma)
type Server struct {
	certPath   string
	privateKey string
	rdr        renderer
}

func New(cert config.Fixtures) *Server {
	rdr := newRenderer()
	return &Server{cert.Path, cert.PrivateKey, rdr}
}

func (s *Server) registerRoutes() http.Handler {
	r := chi.NewRouter()
	// r := mux.NewRouter()

	// r.Use(s.logRequestMiddleware)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	// r.Use(NewStructuredLogger(logrus.StandardLogger()))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Register pprof handlers
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	r.Get("/nucleus/authToken", s.sessionHandler)
	// r.Post("/dc/submit", s.someHandler)
	// r.HandleFunc("/nucleus/check/user/{userID}", s.someHandler)
	// 00A96454: 'nucleus/personas/%s',0
	// 00A96164: 'nucleus/entitlements/%I64d',0
	// 00A96180: 'ofb/products',0
	// 00A96190: 'nucleus/wallets/%I64d',0
	// 00A961A8: 'ofb/purchase/%I64d/%s',0
	// 00A961C0: 'nucleus/wallets/%I64d/%s/%d/%s',0
	// 00A961E0: 'nucleus/refundAbilities/%I64d',0
	// 00A96200: 'relationships/acknowledge/nucleus:%I64d/%I64d',0
	// 00A96230: 'relationships/acknowledge/server:%s/%I64d',0
	// 00A9625C: 'relationships/increase/nucleus:%I64d/nucleus:%I64d/%s',0
	// 00A96294: 'relationships/increase/nucleus:%I64d/server:%s/%s',0
	// 00A962C8: 'relationships/increase/server:%s/nucleus:%I64d/%s',0
	// 00A962FC: 'relationships/decrease/nucleus:%I64d/nucleus:%I64d/%s',0
	// 00A96334: 'relationships/decrease/nucleus:%I64d/server:%s/%s',0
	// 00A96368: 'relationships/decrease/server:%s/nucleus:%I64d/%s',0
	// 00A9639C: 'relationships/status/nucleus:%I64d',0
	// 00A963F0: 'relationships/status/server:%s',0
	// 00A96454: 'nucleus/personas/%s',0
	// 00A964F8: 'user/updateUserProfile/%I64d',0
	// 00A96518: 'nucleus/entitlement/%s/useCount/%d',0
	// 00A9653C: 'nucleus/entitlement/%s/status/%s',0
	// 00A96570: 'nucleus/check/%s/%I64d',0
	// 00A96588: 'nucleus/authToken',0
	// 00A9659C: 'nucleus/name/%I64d',0
	// 00A965BC: 'nucleus/entitlements/%I64d',0
	// 00A965D8: 'nucleus/entitlements/%I64d?entitlementTag=%s',0
	// 00A96468: 'relationships/roster/nucleus:%I64d',0
	// 00A9648C: 'relationships/roster/server:%s/bvip/1,3',0
	// 00A964B4: 'relationships/roster/nucleus:%I64d',0
	// 00A964D8: 'relationships/roster/server:%s',0
	r.Get("/relationships/roster/nucleus:{id}", s.relationshipHandler)
	r.Get("/relationships/roster/server:{id}", s.relationshipHandler)
	r.Get("/relationships/roster/server:{id}/bvip/1,3", s.relationshipHandler)

	r.HandleFunc("/nucleus/entitlements/{heroID}", s.entitlementsHandler)
	r.HandleFunc("/nucleus/wallets/{heroID}", s.walletsHandler)
	r.HandleFunc("/ofb/products", s.offersHandler)
	r.HandleFunc("/en/game/store", s.store)

	return r
}

func (s *Server) ListenAndServe(bind, bindSecure string) {
	r := s.registerRoutes()

	go func() { log.Println(http.ListenAndServe(bind, r)) }()
	go func() {
		srv := &http.Server{
			Addr:     bindSecure,
			Handler:  r,
			ErrorLog: stdlog.New(os.Stderr, "TLS: ", 0),
		}

		log.Println(
			srv.ListenAndServeTLS(s.certPath, s.privateKey),
		)
	}()
}
