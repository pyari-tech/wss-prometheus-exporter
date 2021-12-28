package wsclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const pingPeriod = 4 * time.Second // send pings with this period

type WebSocketClient struct {
	configStr string
	sendBuf   chan []byte
	ctx       context.Context

	mu     sync.RWMutex
	wsconn *websocket.Conn
}

func NewWebSocketClient(wsUrl url.URL) (*WebSocketClient, error) {
	wsClient := WebSocketClient{
		sendBuf: make(chan []byte, 1),
	}
	wsClient.ctx = context.Background()
	wsClient.configStr = wsUrl.String()
	return &wsClient, nil
}

func (wsClient *WebSocketClient) String() string {
	return wsClient.configStr
}

func (wsClient *WebSocketClient) connect() *websocket.Conn {
	wsClient.mu.Lock()
	defer wsClient.mu.Unlock()
	if wsClient.wsconn != nil {
		return wsClient.wsconn
	}

	dialer := *websocket.DefaultDialer
	dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	ws, _, err := dialer.Dial(wsClient.configStr, nil)
	if err != nil {
		wsClient.log("connect", err, fmt.Sprintf("Cannot connect to websocket: %s", wsClient.configStr))
		return nil
	}
	wsClient.log("connect", nil, fmt.Sprintf("connected to websocket to %s", wsClient.configStr))
	wsClient.wsconn = ws
	return wsClient.wsconn
}

func (wsClient *WebSocketClient) Stop() {
	wsClient.closeWs()
}

func (wsClient *WebSocketClient) closeWs() {
	wsClient.mu.Lock()
	if wsClient.wsconn != nil {
		wsClient.wsconn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		wsClient.wsconn.Close()
		wsClient.wsconn = nil
	}
	wsClient.mu.Unlock()
}

func (wsClient *WebSocketClient) Ping() error {
	wsClient.log("ping", nil, "ping pong started")
	ws := wsClient.connect()
	if ws == nil {
		return fmt.Errorf("empty websocket connection")
	}
	err := wsClient.wsconn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(pingPeriod/2))
	if err != nil {
		wsClient.closeWs()
	}
	return err
}

func (wsClient *WebSocketClient) log(f string, err error, msg string) {
	if err != nil {
		fmt.Printf("Error in func: %s, err: %v, msg: %s\n", f, err, msg)
	} else {
		fmt.Printf("Log in func: %s, %s\n", f, msg)
	}
}
