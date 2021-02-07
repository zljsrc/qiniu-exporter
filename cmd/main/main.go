package main

import (
	"flag"
	"log"
	"net/http"

	"qiniu-exporter/cmd/qiniu"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress  = flag.String("web.listen-address", ":9901", "Address to listen on for telemetry")
	metricsPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	qiniuAccessKey = flag.String("qiniu.access-key", "", "Path under which to expose metrics")
	qiniuSecretKey = flag.String("qiniu.secret-key", "", "Path under which to expose metrics")
)

func main() {
	if *qiniuAccessKey == "" {
		log.Fatal("qiniuAccessKey can't be empty")
		return
	}
	if *qiniuSecretKey == "" {
		log.Fatal("qiniuSecretKey can't be empty")
		return
	}
	qiniuMetricsController := qiniu.QiniuMetricsController(*qiniuAccessKey, *qiniuSecretKey)
	prometheus.Register(qiniuMetricsController)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Mirth Channel Exporter</title></head>
             <body>
             <h1>Mirth Channel Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))

}
