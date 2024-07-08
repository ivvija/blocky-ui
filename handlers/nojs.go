package handlers

import (
	"blocky-ui/api"
	"blocky-ui/components"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	status, err := api.GetStatus(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	components.Dash(status.Status).Render(r.Context(), w)
}

func Post(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// TODO: maybe handle HTMX by checking the headers:
	//   r.Header.Get("HX-Request") == "true"

	// TODO: real status value
	status := 0

	if r.Form.Has("toggle") {
		// TODO: call blocky API
		status = 1
	}

	if r.Form.Has("togglePause") {
		// TODO: call blocky API
		status = 2
	}

	if r.Form.Has("flush") {
		// TODO: call blocky API
	}

	if r.Form.Has("refresh") {
		// TODO: call blocky API
	}

	if r.Form.Has("query") && r.Form.Has("type") {
		// TODO: call blocky API
		// TODO: insert result
	}

	components.Dash(status).Render(r.Context(), w)
}
