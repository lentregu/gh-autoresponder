package main

import (
	"log"

	"github.com/lentregu/gh-autoresponder/internal/server"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting application...")

	srv := server.New("8080")
	srv.Start()
}
