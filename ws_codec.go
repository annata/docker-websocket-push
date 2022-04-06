package main

import (
	"golang.org/x/net/websocket"
)

var WsPing = websocket.Codec{Marshal: marshal, Unmarshal: unmarshal}

var emptyMsg = make([]byte, 0)

func marshal(v interface{}) (msg []byte, payloadType byte, err error) {
	return emptyMsg, websocket.PingFrame, err
}

func unmarshal(msg []byte, payloadType byte, v interface{}) (err error) {
	return nil
}
