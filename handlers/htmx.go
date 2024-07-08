package handlers

import (
	"blocky-ui/api"
	"blocky-ui/components"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
	status, err := api.GetStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	components.HeaderBar(status.Status).Render(r.Context(), w)
}

func Toggle(w http.ResponseWriter, r *http.Request) {
	status, err := api.GetStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status.Status == api.Enabled {
		status, err = api.SetDisabled(r.Context())
	} else {
		status, err = api.SetEnabled(r.Context())
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	components.HeaderBar(status.Status).Render(r.Context(), w)
}

func TogglePause(w http.ResponseWriter, r *http.Request) {
	status, err := api.GetStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if status.Status == api.Paused {
		status, err = api.SetEnabled(r.Context())
	} else {
		status, err = api.SetPaused(r.Context(), 10)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	components.HeaderBar(status.Status).Render(r.Context(), w)
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
