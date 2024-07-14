package handlers

import (
	"blocky-ui/api"
	"blocky-ui/components"
	"blocky-ui/settings"
	"fmt"
	"net/http"
	"time"
)

func Get(w http.ResponseWriter, r *http.Request) {
	isHX := "true" == r.Header.Get("HX-Request")

	handleStatus(w, r, isHX)
}

func Post(w http.ResponseWriter, r *http.Request) {
	isHX := "true" == r.Header.Get("HX-Request")

	r.ParseForm()

	switch {
	case r.Form.Has("toggle"):
		handleToggle(w, r, isHX)

	case r.Form.Has("togglePause"):
		handleTogglePause(w, r, isHX, settings.PauseDuration)

	case r.Form.Has("flush"):
		handleFlush(w, r, isHX)

	case r.Form.Has("refresh"):
		handleRefresh(w, r, isHX)

	case r.Form.Has("query") && r.Form.Has("type"):
		handleQuery(w, r, isHX, r.Form.Get("query"), r.Form.Get("type"))

	default:
		if isHX {
			http.Error(w, "HX-POST: unexpected data", http.StatusUnprocessableEntity)
			return
		}

		handleStatus(w, r, false)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request, isHX bool) {
	status, err := api.Status(r.Context())
	if err != nil {
		err = fmt.Errorf("handleStatus() %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if isHX {
		components.HeaderBar(*status).Render(r.Context(), w)
	} else {
		components.Dash(*status).Render(r.Context(), w)
	}
}

func handleToggle(w http.ResponseWriter, r *http.Request, isHX bool) {
	newStatus, err := api.Toggle(r.Context())
	if err != nil {
		err = fmt.Errorf("handleToggle() %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isHX {
		components.HeaderBar(*newStatus).Render(r.Context(), w)
	} else {
		components.Dash(*newStatus).Render(r.Context(), w)
	}
}

func handleTogglePause(w http.ResponseWriter, r *http.Request, isHX bool, duration time.Duration) {
	newStatus, err := api.TogglePause(r.Context(), duration)
	if err != nil {
		err = fmt.Errorf("handleTogglePause() %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isHX {
		components.HeaderBar(*newStatus).Render(r.Context(), w)
	} else {
		components.Dash(*newStatus).Render(r.Context(), w)
	}
}

func handleFlush(w http.ResponseWriter, r *http.Request, isHX bool) {
	err := api.Flush(r.Context())
	if err != nil {
		err = fmt.Errorf("handleFlush() %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if isHX {
		w.WriteHeader(http.StatusNoContent)
	} else {
		handleStatus(w, r, false)
	}
}
func handleRefresh(w http.ResponseWriter, r *http.Request, isHX bool) {
	err := api.Refresh(r.Context())
	if err != nil {
		err = fmt.Errorf("handleRefresh() %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if isHX {
		w.WriteHeader(http.StatusNoContent)
	} else {
		handleStatus(w, r, false)
	}
}

func handleQuery(w http.ResponseWriter, r *http.Request, isHX bool, query string, recordType string) {
	status, err := api.Status(r.Context())
	if err != nil {
		err = fmt.Errorf("handleQuery() status: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := api.Query(r.Context(), query, recordType)
	if err != nil {
		err = fmt.Errorf("handleQuery() query %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if isHX {
		components.QueryResult(*result).Render(r.Context(), w)
	} else {
		components.DashQueryResult(*status, *result).Render(r.Context(), w)
	}
}
