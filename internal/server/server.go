package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	shutdown   chan os.Signal
}

func New(port string, issueHandler http.Handler) *Server {
	router := setupRouter(issueHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
		shutdown: make(chan os.Signal, 1),
	}
}

func (s *Server) Start() error {

	signal.Notify(s.shutdown, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("Server starting on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case <-s.shutdown:
		return s.Stop()
	}
}

func (s *Server) Stop() error {

	log.Println("Initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
