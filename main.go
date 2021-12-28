package main

import (
	"flag"
	"net/url"
	"time"

	"github.com/pyari-tech/wss-prometheus-exporter/promserver"
	"github.com/pyari-tech/wss-prometheus-exporter/wsclient"
)

var bind = flag.String("bind", "0.0.0.0:10101", "prometheus exporter server:port")
var wsScheme = flag.String("scheme", "ws", "websocket scheme ws/wss")
var uri = flag.String("uri", "127.0.0.1:8080", "websocket server:port")
var wsChannel = flag.String("channel", "frontend", "websocket channel")
var secondsInterval = flag.Int("sec", 60, "interval between checks, minimum allowed is 5sec")
var pingForCheck = flag.Bool("pingcheck", false, "use /ping for websocket check as well")

func main() {
	var err error
	flag.Parse()
	if *secondsInterval < 5 {
		*secondsInterval = 5
	}
	promserver.SecondsInterval = time.Duration(*secondsInterval) * time.Second
	promserver.PingForCheck = *pingForCheck
	wsUrl := url.URL{Host: *uri, Path: *wsChannel, Scheme: *wsScheme}
	promserver.WebSocketClient, err = wsclient.NewWebSocketClient(wsUrl)
	if err != nil {
		panic(err)
	}
	promserver.Server(*bind)
}
