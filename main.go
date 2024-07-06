package main

import (
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"

	"blocky-ui/components"
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

	router.Get("/", templ.Handler(components.Dash()).ServeHTTP)

	log.Fatal(http.ListenAndServe(":3000", router))
}
