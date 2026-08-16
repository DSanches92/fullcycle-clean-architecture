package webserver

import (
	"database/sql"
	"net/http"

	"github.com/DSanches92/fullcycle-clean-architecture/internal/infra/web"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	OrderHandler  *web.WebOrderHandler
	DB            *sql.DB
	ServerAddress string
}

func NewWebServer(orderHandler *web.WebOrderHandler, db *sql.DB, serverAddress string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		OrderHandler:  orderHandler,
		DB:            db,
		ServerAddress: serverAddress,
	}
}

func (s *WebServer) Setup() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Get("/healthz", s.LivenessHandler)
	s.Router.Get("/readyz", s.ReadinessHandler)

	s.Router.Post("/orders", s.OrderHandler.Create)
	s.Router.Get("/orders", s.OrderHandler.List)
}

func (s *WebServer) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *WebServer) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.DB.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *WebServer) Start() error {
	return http.ListenAndServe(s.ServerAddress, s.Router)
}
