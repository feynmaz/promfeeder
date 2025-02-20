package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/feynmaz/pkg/http/middleware"
	"github.com/feynmaz/pkg/logger"
	"github.com/feynmaz/promfeeder/config"
	docs "github.com/feynmaz/promfeeder/openapi"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	cfg    *config.Config
	logger *logger.Logger
}

func New(cfg *config.Config, logger *logger.Logger) *Server {

	return &Server{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", s.cfg.Server.Port),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		Handler:      s.getRouter(),
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
	}

	s.logger.Info().Msgf("server started on port %d", s.cfg.Server.Port)
	return srv.ListenAndServe()
}

// @title Promfeeder API
// @version 1.0
// @description This is a service which feeds prometheus with metrics.
// @contact.name Nikolai Mazein
// @contact.email feynmaz@gmail.com
// @BasePath /
// @x-servers [{"url": "https://promfeeder.testshift.webtm.ru/", "description": "dev"}]
func (s *Server) getRouter() *chi.Mux {
	router := chi.NewMux()

	// Middleware
	router.Use(middleware.RequestIDMiddleware)
	router.Use(middleware.NewLoggingMiddleware(s.logger))

	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by status code",
		},
		[]string{"code"},
	)
	prometheus.MustRegister(requestCounter)
	router.Use(func(next http.Handler) http.Handler {
		return promhttp.InstrumentHandlerCounter(requestCounter, next)
	})

	// Profiler
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))

	// Metrics
	router.Handle("/metrics", promhttp.Handler())

	// Swagger
	docs.SwaggerInfo.BasePath = s.cfg.AppBasePath
	router.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		swaggerHandler := httpSwagger.Handler(
			httpSwagger.URL(
				fmt.Sprintf("%s/swagger/doc.json", s.cfg.AppBaseURL),
			),
		)
		swaggerHandler.ServeHTTP(w, r)
	})

	// Methods
	router.Get("/get/{code}", s.Get)

	return router
}

func (s *Server) Shutdown() {
	s.logger.Info().Msg("graceful server shutdown")
}
