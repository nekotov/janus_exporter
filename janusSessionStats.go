package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type SessionStats struct {
	In  Stats `json:"in"`
	Out Stats `json:"out"`
}

type Stats struct {
	Packets int `json:"packets"`
	Bytes   int `json:"bytes"`
}

type SelectedPair struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

type Stream struct {
	Type        string `json:"type"`
	MIndex      int    `json:"mindex"`
	Mid         string `json:"mid"`
	Codec       string `json:"codec"`
	Subscribers int    `json:"subscribers"`
}

type PluginSpecific struct {
	Bitrate int      `json:"bitrate"`
	Streams []Stream `json:"streams"`
}

type WebRTC struct {
	ICE struct {
		SelectedPair string `json:"selected-pair"`
	} `json:"ice"`
	DTLS struct {
		Stat Stats `json:"stats"`
	} `json:"dtls"`
}

type JanusSession struct {
	Plugin         string         `json:"plugin"`
	PluginSpecific PluginSpecific `json:"plugin_specific"`
	WebRTC         WebRTC         `json:"webrtc"`
	Stats          SessionStats   `json:"stats"`
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

//sessions := getJanusSessionsList(janusHost, janusAdminToken)
////	//var wg sync.WaitGroup
////	//var mu sync.Mutex
////
////	for _, session := range sessions {
////		handlers := getJanusHandlersList(janusHost, janusAdminToken, session)
////		for _, handler := range handlers {
////			//wg.Add(100)
////			go func(session int64, handler int64) {
////				//defer wg.Done()
////				s := getJanusHandlerInfo(janusHost, janusAdminToken, session, handler)
////				//if s.PluginSpecific.Bitrate != 0 {
////				//	for _, stream := range s.PluginSpecific.Streams {
////				//		if stream.Subscribers > 0 {
////				//			//mu.Lock()
////				//			s.WebRTC.ICE.SelectedPair = strings.Split(strings.Split(s.WebRTC.ICE.SelectedPair, "<->")[1], ":")[0]
////				//			fmt.Println(s)
////				//			//mu.Unlock()
////				//			break
////				//		}
////				//	}
////				//}
////				fmt.Println(s)
////			}(session, handler)
////		}
////	}
////
////	//wg.Wait()
