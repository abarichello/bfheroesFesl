package server

import (
	"net/http"

	"github.com/Synaxis/bfheroesFesl/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	logger = logrus.WithField("pkg", "server")
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

	r.Use(s.logRequestMiddleware)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	// r.Use(NewStructuredLogger(logrus.StandardLogger()))
	r.Use(middleware.Recoverer)
	// r.Use(middleware.Timeout(60 * time.Second))

	// TODO r.Post("/dc/submit", someHandler)
	// TODO: user/updateUserProfile/%d

	// Profiling
	// r.Route("/debug", func(r chi.Router) {
	// 	r.Get("/pprof/", pprof.Index)
	// 	r.Get("/pprof/cmdline", pprof.Cmdline)
	// 	r.Get("/pprof/profile", pprof.Profile)
	// 	r.Get("/pprof/symbol", pprof.Symbol)
	// 	r.Get("/pprof/trace", pprof.Trace)
	// })

	// Client-relationship (i.e. friends, server bookmarks)
	r.Route("/relationships", func(r chi.Router) {
		// /roster/nucleus:%d
		r.Get("/roster/nucleus:{userID}", s.relationshipRosterNucleus)
		// /roster/server:%s
		r.Get("/roster/server:{serverName}", s.relationshipRosterServer)
		// /roster/server:%s/bvip/1,3
		r.Get("/roster/server:{serverName}/bvip/1,3", s.relationshipRosterServer)

		// TODO: /acknowledge/nucleus:%d/%d
		// TODO: /acknowledge/server:%s/%d
		// TODO: /increase/nucleus:%d/nucleus:%d/%s
		// TODO: /increase/nucleus:%d/server:%s/%s
		// TODO: /increase/server:%s/nucleus:%d/%s
		// TODO: /decrease/nucleus:%d/nucleus:%d/%s
		// TODO: /decrease/nucleus:%d/server:%s/%s
		// TODO: /decrease/server:%s/nucleus:%d/%s
		// TODO: /status/nucleus:%d
		// TODO: /status/server:%s
	})

	// Nuclues (authentication and account data)
	r.Route("/nucleus", func(r chi.Router) {
		// /authToken
		r.Get("/authToken", s.nucleusAuthToken)

		// /entitlements/%d
		// /entitlements/%d?entitlementTag=%s
		r.Get("/entitlements/{heroID}", s.nucleusEntitlements)
		// TODO: /entitlement/%s/useCount/%d
		// TODO: /entitlement/%s/status/%s

		// TODO: /wallets/%d/%s/%d/%s
		// TODO: /wallets/%d
		r.Get("/wallets/{heroID}", s.walletsHandler)

		// TODO: /check/%s/%d
		r.Get("/check/user/{userID}", s.nucleusCheckUser)

		// TODO: /personas/%s
		// TODO: /refundAbilities/%d
		// TODO: /personas/%s
		// TODO: /name/%d
	})

	// Overlay
	r.Route("/ofb", func(r chi.Router) {
		r.Get("/products", s.ofbProducts)
		// TODO: purchase/%d/%s
	})

	// Game-client
	r.Route("/en/game", func(r chi.Router) {
		r.Get("/store", s.store)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		logrus.Warn("Not found URL: ", r.URL.String())
		logRequest(r)
		http.NotFound(w, r)
	})

	return r
}

func (s *Server) ListenAndServe(bind, bindSecure string) {
	r := s.registerRoutes()

	go func() { log.Println(http.ListenAndServe(bind, r)) }()
	go func() {
		srv := &http.Server{
			Addr:    bindSecure,
			Handler: r,
			//ErrorLog: stdlog.New(os.Stderr, "TLS: ", 0),
		}

		log.Println(
			srv.ListenAndServeTLS(s.certPath, s.privateKey),
		)
	}()
}
