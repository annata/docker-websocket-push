package main

import (
	"encoding/json"
	"golang.org/x/net/websocket"
)

var WsSubscribe = websocket.Codec{Marshal: marshalSubscribe, Unmarshal: unmarshalSubscribe}

func marshalSubscribe(v interface{}) (msg []byte, payloadType byte, err error) {
	bytes, err := json.Marshal(v)
	return bytes, websocket.TextFrame, err
}

func unmarshalSubscribe(msg []byte, payloadType byte, v interface{}) (err error) {
	if payloadType == websocket.TextFrame {
		_ = json.Unmarshal(msg, v)
	}
	return nil
}

type Subscribe struct {
	Op    string `json:"op"`
	Token string `json:"token"`
}
