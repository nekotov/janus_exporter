package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type PayloadSessions struct {
	Janus       string `json:"janus"`
	Transaction string `json:"transaction"`
	AdminSecret string `json:"admin_secret"`
}

type ResponseSessions struct {
	Janus       string  `json:"janus"`
	Transaction string  `json:"transaction"`
	Sessions    []int64 `json:"sessions"`
}

func getJanusSessionsList(host string, token string) []int64 {
	data := PayloadSessions{
		Janus:       "list_sessions",
		Transaction: "prometheus",
		AdminSecret: token,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, host, body)
	if err != nil {
		return []int64{}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []int64{}
	}

	var response ResponseSessions
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return []int64{}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// ignoring
		}
	}(resp.Body)

	return response.Sessions

}

func getJanusSessionsCount(host string, token string) float64 {

	sessionCount := float64(len(getJanusSessionsList(host, token)))

	return sessionCount

}
