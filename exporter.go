package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	sessionsCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "janus_sessions",
		Help: "Monitoring janus sessions",
	}, []string{})

	handlersCounterInt int
	handlersCounter    = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "janus_handlers",
		Help: "Monitoring janus handlers",
	}, []string{})

	dynamicIPListMu sync.Mutex
	dynamicIPList   = make(map[string]int)
	ipList          = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "janus_clients_ips",
		Help: "Monitoring janus clients ips",
	}, []string{"ip"})
)

func AddIP(ip string) {
	dynamicIPListMu.Lock()
	defer dynamicIPListMu.Unlock()
	dynamicIPList[ip]++
}

func recordMetrics() {
	go func() {
		for {
			sessionsCounter.WithLabelValues().Set(getJanusSessionsCount(janusHost, janusAdminToken))
			time.Sleep(time.Second * 5)
		}
	}()

	go func() {
		for {

			handlersCounterInt = 0
			dynamicIPList = make(map[string]int)
			ipList.Reset()

			sessions := getJanusSessionsList(janusHost, janusAdminToken)
			var wg sync.WaitGroup

			for _, session := range sessions {
				handlers := getJanusHandlersList(janusHost, janusAdminToken, session)
				handlersCounterInt += len(handlers) // Counting handlers
				for _, handler := range handlers {
					wg.Add(1)
					go func(session int64, handler int64) {
						defer wg.Done()
						s := getJanusHandlerInfo(janusHost, janusAdminToken, session, handler)
						if s.PluginSpecific.Bitrate != 0 {
							for _, stream := range s.PluginSpecific.Streams {
								if stream.Subscribers > 0 {
									s.WebRTC.ICE.SelectedPair = strings.Replace(strings.Split(strings.Split(s.WebRTC.ICE.SelectedPair, "<->")[1], ":")[0], " ", "", -1)
									AddIP(s.WebRTC.ICE.SelectedPair)
									fmt.Println(s)
									break
								}
							}
						}
					}(session, handler)
				}
			}

			wg.Wait()

			handlersCounter.WithLabelValues().Set(float64(handlersCounterInt))

			dynamicIPListMu.Lock()
			for ip, count := range dynamicIPList {
				ipList.With(prometheus.Labels{"ip": ip}).Set(float64(count))
			}
			dynamicIPListMu.Unlock()

			time.Sleep(time.Second * 30)
		}
	}()
}

func init() {
	prometheus.MustRegister(sessionsCounter)
	prometheus.MustRegister(handlersCounter)
	prometheus.MustRegister(ipList)
}

var (
	janusHost         = "http://localhost:7088/admin"
	janusAdminToken   = "janusoverlord"
	janusExporterHost = ":8090"
	janusExporterPath = "/metrics"
)

func main() {

	flag.String("janus-host", janusHost, "Janus host")
	flag.String("janus-admin-token", janusAdminToken, "Janus admin token")
	flag.String("janus-exporter-host", janusExporterHost, "Janus exporter host")
	flag.String("janus-exporter-path", janusExporterPath, "Janus exporter path")
	flag.Parse()

	recordMetrics()

	srv := http.NewServeMux()
	srv.Handle(janusExporterPath, promhttp.Handler())
	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
					<head><title>Janus Exporter</title></head>
					<body>
					<h1>Janus Exporter</h1>
					<p><a href='` + janusExporterPath + `'>Metrics</a></p>
					</body>
					</html>`))
		if err != nil {
			return
		}
	})

	if err := http.ListenAndServe(janusExporterHost, srv); err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}
