package main

import (
	"log"

	"github.com/lentregu/gh-autoresponder/internal/server"
	"github.com/lentregu/gh-autoresponder/internal/webhook"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	issueHandler := webhook.NewIssueHandler()

	srv := server.New("8080", issueHandler)
	srv.Start()
}
