package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/orcaman/concurrent-map"
	"golang.org/x/net/websocket"
)

var m = cmap.New()

func send(key string, value string) {
	wsObj, ok := m.Get(key)
	if ok {
		ws := wsObj.(*websocket.Conn)
		err := websocket.Message.Send(ws, value)
		if err != nil {
			ws.Close()
		}
	}
}

func sendAll(value string) {
	for t := range m.IterBuffered() {
		wsObj := t.Val
		ws := wsObj.(*websocket.Conn)
		err := websocket.Message.Send(ws, value)
		if err != nil {
			ws.Close()
		}
	}
}

func websocketHandle(ws *websocket.Conn) {
	defer ws.Close()
	token := ws.Request().URL.Query().Get("token")
	tokenTmp := ws.Request().Header.Get("token")
	if tokenTmp != "" {
		token = tokenTmp
	}
	if token == "" {
		return
	}
	sum256 := sha256.Sum256([]byte(token))
	user := prefix + hex.EncodeToString(sum256[:])

	if ok := m.SetIfAbsent(user, ws); !ok {
		return
	}
	defer m.Remove(user)

	add(user)
	defer del(user)

	var data string
	for {
		e := websocket.Message.Receive(ws, &data)
		if e != nil {
			return
		}
	}
}
