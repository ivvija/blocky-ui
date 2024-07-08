package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"blocky-ui/handlers"
	"blocky-ui/settings"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	fs := http.FileServer(http.Dir("assets"))
	router.Get("/assets/*", http.StripPrefix("/assets/", fs).ServeHTTP)

	// no-JS handlers
	router.Get("/", handlers.Get)
	router.Post("/", handlers.Post)

	// HTMX handlers
	router.Get("/status", handlers.Status)
	router.Post("/toggle", handlers.Toggle)
	router.Post("/togglePause", handlers.TogglePause)
	router.Post("/refresh", handlers.Refresh)
	router.Post("/flush", handlers.Flush)
	router.Post("/query", handlers.Query)

	addr := fmt.Sprintf("%s:%s", settings.Host, settings.Port)
	log.Printf("Listening on http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
