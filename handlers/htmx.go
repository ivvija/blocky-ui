package handlers

import (
	"blocky-ui/components"
	"net/http"
)

func Toggle(w http.ResponseWriter, r *http.Request) {
	// TODO: call blocky API
	// TODO: real status value
	components.HeaderBar(1).Render(r.Context(), w)
}

func TogglePause(w http.ResponseWriter, r *http.Request) {
	// TODO: call blocky API
	// TODO: real status value
	components.HeaderBar(2).Render(r.Context(), w)
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// TODO: call blocky API
	w.WriteHeader(http.StatusNotImplemented)
}

func Flush(w http.ResponseWriter, r *http.Request) {
	// TODO: call blocky API
	w.WriteHeader(http.StatusNotImplemented)
}

func Query(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// TODO: call blocky API

	w.Write([]byte(r.Form.Encode()))
	w.WriteHeader(http.StatusOK)
}
