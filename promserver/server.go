package promserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/pyari-tech/wss-prometheus-exporter/wsclient"
)

const (
	success = 1
	failure = 0
)

var (
	SecondsInterval time.Duration
	PingForCheck    bool
	WebSocketClient *wsclient.WebSocketClient
)

func isUp() bool {
	if errPing := WebSocketClient.Ping(); errPing != nil {
		fmt.Printf("%s is down!\n%v\n", WebSocketClient.String(), errPing.Error())
		return false
	}
	WebSocketClient.Stop()
	fmt.Printf("%s is up!\n", WebSocketClient.String())
	return true
}

func recordMetrics() {
	/*
		uri, errParse := url.Parse(WebSocketClient.String())
		host, port, errSplit := net.SplitHostPort(uri.Host)
		gaugeName := "websocket_ping"
		if errParse == nil && errSplit == nil {
			gaugeName = fmt.Sprintf("%s_%s_%s", gaugeName, uri.Host, uri.Port)
		}
	*/
	WebSocketPing := promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_ping",
			Help: "Connection to WebSocket with Ping",
		},
	)

	for {
		if isUp() == true {
			WebSocketPing.Set(success)
		} else {
			WebSocketPing.Set(failure)
		}
		time.Sleep(SecondsInterval)
	}
}

func ping(w http.ResponseWriter, req *http.Request) {
	if PingForCheck {
		if isUp() == false {
			fmt.Fprintf(w, "crash")
		}
	}
	fmt.Fprintf(w, "pong")
}

func Server(listenAt string) {
	go recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/ping", ping)

	http.ListenAndServe(listenAt, nil)
}
