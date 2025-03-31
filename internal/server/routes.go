package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func setupRouter(issueHandler http.Handler) *mux.Router {
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.Handle("/webhook", issueHandler).Methods("POST")
	return r
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
