
## wss-prometheus-exporter

> a Prometheus' Exporter for WebSocket service

* it connects to the configured WS://URI or WSS://URI, checks for `ping`

> * if everything goes smooth publishes `1` for gauge `websocket_ping`
>
> * otherwise sets it to `0`

* alerts can be raised for `expr: websocket_ping == 0`


### How To

* build using `./build.sh`

* all available configurations are provided via config switch, could be checked as `./out/wss-exporter -help`

* could run as `./out/wss-exporter`; place it in `/usr/local/bin/wss-exporter` for infra instances to run easily

* example usage:

```
wss-exporter \
        -bind 0.0.0.0:10101 \
        -channel myweb \
        -scheme wss \
        -sec 5 \
        -uri 127.0.0.1:8080
```

> here `bind` is for where exporter will serve at; `uri` is where websocket service is available

* sample [systemd service](extras/wss-exporter.systemd.service) definition could be edited for correct config and placed at `/etc/systemd/system/wss-exporter.service`

> then a simple `systemctl daemon-reload && systemctl start wss-exporter && systemctl status wss-exporter`; shall start the service and show status given the built binary has been placed at an available `$PATH`

---
