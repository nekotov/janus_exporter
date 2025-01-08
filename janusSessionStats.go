package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Stream struct {
	Type        string `json:"type"`
	Subscribers int    `json:"subscribers"`
}

type PluginSpecific struct {
	Bitrate int      `json:"bitrate"`
	Streams []Stream `json:"streams"`
}

type IOStats struct {
	Packets int `json:"packets"`
	Bytes   int `json:"bytes"`
}

type WebRTC struct {
	ICE struct {
		SelectedPair string `json:"selected-pair"`
	} `json:"ice"`
	DTLS struct {
		STATS struct {
			IN  IOStats `json:"in"`
			OUT IOStats `json:"out"`
		}
	} `json:"dtls"`
}

type JanusSession struct {
	Plugin         string         `json:"plugin"`
	PluginSpecific PluginSpecific `json:"plugin_specific"`
	WebRTC         WebRTC         `json:"webrtc,omitempty"`
}

type ResponseHandlers struct {
	Janus       string  `json:"janus"`
	SessionID   int64   `json:"session_id"`
	Transaction string  `json:"transaction"`
	Handles     []int64 `json:"handles"`
}

type ResponseHandlerInfo struct {
	Janus       string       `json:"janus"`
	SessionID   int64        `json:"session_id"`
	Transaction string       `json:"transaction"`
	HandleID    int64        `json:"handle_id"`
	JanusSes    JanusSession `json:"info"`
}

func getJanusHandlersList(host string, token string, session int64) []int64 {

	data := PayloadSessions{
		Janus:       "list_handles",
		Transaction: "prometheus",
		AdminSecret: token,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, host+"/"+strconv.Itoa(int(session)), body)
	if err != nil {
		return []int64{}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []int64{}
	}

	var response ResponseHandlers
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return []int64{}
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// ignoring
		}
	}(resp.Body)

	return response.Handles

}

func getJanusHandlerInfo(host string, token string, session int64, handler int64) JanusSession {

	data := PayloadSessions{
		Janus:       "handle_info",
		Transaction: "prometheus",
		AdminSecret: token,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return JanusSession{}
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, host+"/"+strconv.Itoa(int(session))+"/"+strconv.Itoa(int(handler)), body)
	if err != nil {
		return JanusSession{}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return JanusSession{}
	}

	var response ResponseHandlerInfo

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {

	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// ignoring
		}
	}(resp.Body)

	return response.JanusSes
}
