package main

import (
	"golang.org/x/net/websocket"
)

var WS_PING = websocket.Codec{Marshal: marshal, Unmarshal: unmarshal}

func marshal(v interface{}) (msg []byte, payloadType byte, err error) {
	return msg, websocket.PingFrame, err
}

func unmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return nil
}
