package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var (
	sessionsCounter = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "janus_sessions",
		Help: "Monitoring janus sessions",
	}, []string{"node", "namespace"})
)

func recordMetrics() {
	go func() {
		for {
			sessionsCounter.WithLabelValues("node-1", "namespace-b").Set(getJanusSessionsCount(janusHost, janusAdminToken))
			time.Sleep(time.Second * 5)
		}
	}()
}

func init() {
	prometheus.MustRegister(sessionsCounter)
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
		w.Write([]byte(`<html>
			<head><title>Janus Exporter</title></head>
			<body>
			<h1>Janus Exporter</h1>
			<p><a href='` + janusExporterPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})

	if err := http.ListenAndServe(janusExporterHost, srv); err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}
