package main

import (
	"flag"
	"log"
	"net/http"

	"scriptstudio/script-backend/internal/httpapi"
	"scriptstudio/script-backend/internal/query"
	"scriptstudio/script-backend/internal/service"
	"scriptstudio/script-backend/internal/store"
)

func main() {
	address := flag.String("address", "127.0.0.1:8080", "HTTP listen address")
	database := flag.String("database", "script-studio.db", "bbolt database path")
	flag.Parse()
	repo, err := store.Open(*database)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	commands := service.New(repo)
	api := httpapi.New(commands, query.New(repo))
	log.Printf("script studio listening on http://%s", *address)
	if err := http.ListenAndServe(*address, api.Handler()); err != nil {
		log.Fatal(err)
	}
}
