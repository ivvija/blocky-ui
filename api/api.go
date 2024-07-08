package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"blocky-ui/settings"
)

const (
	Enabled = iota
	Disabled
	Paused
)

type Status struct {
	Status        int
	PausedSeconds int
}

type Query struct {
	Error string
}

var client = &http.Client{Timeout: 10 * time.Second}

func sendRequest(req *http.Request, v any) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}

	if v == nil {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

type statusResponse struct {
	AutoEnableInSec int      `json:"autoEnableInSec"`
	DisabledGroups  []string `json:"disabledGroups"`
	Enabled         bool     `json:"enabled"`
}

func GetStatus(ctx context.Context) (*Status, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settings.ApiBaseUrl+"/blocking/status", nil)
	if err != nil {
		return nil, err
	}

	res := statusResponse{}
	if err := sendRequest(req, &res); err != nil {
		return nil, err
	}

	if res.AutoEnableInSec > 0 {
		return &Status{Paused, res.AutoEnableInSec}, nil
	} else if !res.Enabled {
		return &Status{Disabled, -1}, nil
	}

	return &Status{Enabled, -1}, nil
}

func SetEnabled(ctx context.Context) (*Status, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settings.ApiBaseUrl+"/blocking/enable", nil)
	if err != nil {
		return nil, err
	}

	if err := sendRequest(req, nil); err != nil {
		return nil, err
	}
	return &Status{Enabled, -1}, nil
}

func SetDisabled(ctx context.Context) (*Status, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settings.ApiBaseUrl+"/blocking/disable", nil)
	if err != nil {
		return nil, err
	}

	if err := sendRequest(req, nil); err != nil {
		return nil, err
	}
	return &Status{Disabled, -1}, nil
}

func SetPaused(ctx context.Context, seconds int) (*Status, error) {
	params := "?duration=5m"
	if seconds > 0 {
		params = fmt.Sprintf("?duration=%ds", seconds)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settings.ApiBaseUrl+"/blocking/disable"+params, nil)
	if err != nil {
		return nil, err
	}

	if err := sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &Status{Paused, seconds}, nil
}

type queryRequest struct {
	Query string `json:"query"`
	Type  string `json:"type"`
}

type QueryResponse struct {
	Reason       string `json:"reason"`
	Response     string `json:"response"`
	ResponseType string `json:"responseType"`
	ReturnCode   string `json:"returnCode"`
}
