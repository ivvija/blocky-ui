package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"blocky-ui/settings"
)

const (
	Enabled = iota
	Disabled
	Paused
)

type StatusResponse struct {
	Status        int
	PausedSeconds int
}
type QueryResponse struct {
	RecordType    string
	Reason        string
	ResponseTable [][]string
	ResponseType  string
	ReturnCode    string
}

var client = &http.Client{Timeout: 30 * time.Second}

func sendRequest(req *http.Request, v any) error {
	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
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

type statusApiResponse struct {
	AutoEnableInSec int      `json:"autoEnableInSec"`
	DisabledGroups  []string `json:"disabledGroups"`
	Enabled         bool     `json:"enabled"`
}

func Status(ctx context.Context) (*StatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settings.ApiBaseUrl+"/blocking/status", nil)
	if err != nil {
		return nil, err
	}

	apiRes := statusApiResponse{}
	if err := sendRequest(req, &apiRes); err != nil {
		return nil, err
	}

	if apiRes.AutoEnableInSec > 0 {
		return &StatusResponse{Paused, apiRes.AutoEnableInSec}, nil
	} else if !apiRes.Enabled {
		return &StatusResponse{Disabled, -1}, nil
	}

	return &StatusResponse{Enabled, -1}, nil
}

func enable(ctx context.Context) (*StatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settings.ApiBaseUrl+"/blocking/enable", nil)
	if err != nil {
		return nil, err
	}

	if err := sendRequest(req, nil); err != nil {
		return nil, err
	}
	return &StatusResponse{Enabled, -1}, nil
}

func disable(ctx context.Context) (*StatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settings.ApiBaseUrl+"/blocking/disable", nil)
	if err != nil {
		return nil, err
	}

	if err := sendRequest(req, nil); err != nil {
		return nil, err
	}
	return &StatusResponse{Disabled, -1}, nil
}

func pause(ctx context.Context, duration time.Duration) (*StatusResponse, error) {
	params := fmt.Sprintf("?duration=%s", duration)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, settings.ApiBaseUrl+"/blocking/disable"+params, nil)
	if err != nil {
		return nil, err
	}

	if err := sendRequest(req, nil); err != nil {
		return nil, err
	}

	return &StatusResponse{Paused, int(duration.Seconds())}, nil
}

func Toggle(ctx context.Context) (*StatusResponse, error) {
	status, err := Status(ctx)
	if err != nil {
		return nil, err
	}

	if status.Status == Enabled {
		status, err = disable(ctx)
	} else {
		status, err = enable(ctx)
	}
	if err != nil {
		return nil, err
	}

	return status, nil
}

func TogglePause(ctx context.Context, duration time.Duration) (*StatusResponse, error) {
	status, err := Status(ctx)
	if err != nil {
		return nil, err
	}

	if status.Status == Paused {
		status, err = enable(ctx)
	} else {
		status, err = pause(ctx, duration)
	}
	if err != nil {
		return nil, err
	}

	return status, nil
}

func Refresh(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, settings.ApiBaseUrl+"/lists/refresh", nil)
	if err != nil {
		return err
	}

	if err := sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}

func Flush(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, settings.ApiBaseUrl+"/cache/flush", nil)
	if err != nil {
		return err
	}

	if err := sendRequest(req, nil); err != nil {
		return err
	}
	return nil
}

type queryApiResponse struct {
	Reason       string `json:"reason"`
	Response     string `json:"response"`
	ResponseType string `json:"responseType"`
	ReturnCode   string `json:"returnCode"`
}

func Query(ctx context.Context, query string, recordType string) (*QueryResponse, error) {
	queryJson, err := json.Marshal(map[string]string{
		"query": query,
		"type":  recordType,
	})
	if err != nil {
		return nil, err
	}

	queryBody := bytes.NewBuffer(queryJson)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, settings.ApiBaseUrl+"/query", queryBody)
	if err != nil {
		return nil, err
	}

	apiRes := queryApiResponse{}
	if err := sendRequest(req, &apiRes); err != nil {
		return nil, err
	}

	log.Println(apiRes)
	log.Println(apiRes.Response)

	return &QueryResponse{
		recordType,
		apiRes.Reason,
		parseQueryResponse(apiRes.Response),
		apiRes.ResponseType,
		apiRes.ReturnCode,
	}, nil
}

func parseQueryResponse(response string) [][]string {
	re := regexp.MustCompile("(A|AAAA|CNAME|PTR) \\((.*)\\)")

	recordStrings := strings.Split(response, ", ")
	table := make([][]string, 0)

	for _, recordString := range recordStrings {
		if recordString == "" {
			continue
		}

		match := re.FindStringSubmatch(recordString)
		if match == nil {
			table = append(table, strings.Split(recordString, "\t"))
		} else {
			table = append(table, []string{match[1], match[2]})
		}
	}
	return table
}
