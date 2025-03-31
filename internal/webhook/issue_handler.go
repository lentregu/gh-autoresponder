package webhook

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/go-github/v57/github"
)

type IssueHandler struct {
}

func NewIssueHandler() http.Handler {
	return &IssueHandler{}
}

func (h *IssueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var event github.IssuesEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("Error decoding webhook: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if event.GetAction() == "opened" {
		log.Printf("The issue #%d has been opened", event.GetIssue().GetNumber())
	}

	w.WriteHeader(http.StatusOK)
}
