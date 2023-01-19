package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"gitlab.com/napspan/SampleCompany/api/handlers"
	"gitlab.com/napspan/SampleCompany/api/middlewares"
)

var (
	addr    = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	version string
)

func init() {
	version = os.Getenv("IMAGE_TAG")
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// logs
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("version", version).
		Logger()

	// graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	r := mux.NewRouter()
	c := alice.New(hlog.NewHandler(log), hlog.AccessHandler(accessLogger))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	r.HandleFunc("/", handlers.Index).Methods("GET")

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/", handlers.Index).Methods("GET")
	s.HandleFunc("/GetAllComputers", middlewares.Chain(handlers.GetAllSavedComputers, middlewares.ValidateAuthorization())).Methods("GET")
	s.HandleFunc("/computer", middlewares.Chain(handlers.CreateNewComputer, middlewares.ValidateContentType(), middlewares.ValidateAuthorization())).Methods("POST")
	s.HandleFunc("/computer/assign", middlewares.Chain(handlers.AssignComputer, middlewares.ValidateContentType(), middlewares.ValidateAuthorization())).Methods("POST")
	s.HandleFunc("/computer/unsassign", middlewares.Chain(handlers.DisassignComputer, middlewares.ValidateContentType(), middlewares.ValidateAuthorization())).Methods("POST")
	s.HandleFunc("/computer/by-employee", middlewares.Chain(handlers.GetAllSavedComputersFromEmployee, middlewares.ValidateContentType(), middlewares.ValidateAuthorization())).Methods("POST")
	s.HandleFunc("/computer", middlewares.Chain(handlers.UpdateComputer, middlewares.ValidateContentType(), middlewares.ValidateAuthorization())).Methods("PUT")
	s.HandleFunc("/computer", middlewares.Chain(handlers.DeleteComputer, middlewares.ValidateContentType(), middlewares.ValidateAuthorization())).Methods("DELETE")

	http.Handle("/favicon.ico", http.NotFoundHandler())

	srv := &http.Server{
		Addr:    *addr,
		Handler: c.Then(r),
	}
	go serveHTTP(srv)

	<-quit

	log.Info().Msg("Shutting down API...")

	// Gracefully shutdown connections
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}

// accessLogger register every requests except invocations to resource "api/v1/customlog"
func accessLogger(r *http.Request, status, size int, dur time.Duration) {
	hlog.FromRequest(r).Info().
		Str("host", r.Host).
		Int("status", status).
		Str("url", r.RequestURI).
		Str("method", r.Method).
		Int("size", size).
		Dur("duration_ms", dur).
		Msg("request")
}

func serveHTTP(srv *http.Server) {
	log.Info().Msgf("API started at %s", srv.Addr)
	err := srv.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Error().Err(err).Msg("Starting Server listener failed")
	}
}
