package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/orcaman/concurrent-map"
	"golang.org/x/net/websocket"
	"sync/atomic"
)

var m = cmap.New()
var c int32 = 0

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
	tuple := m.IterBuffered()
	number := (cap(tuple) / 256) + 1
	for i := 0; i < number; i++ {
		go func() {
			for t := range tuple {
				wsObj := t.Val
				ws := wsObj.(*websocket.Conn)
				err := websocket.Message.Send(ws, value)
				if err != nil {
					ws.Close()
				}
			}
		}()
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
	sum := md5.Sum([]byte(token))
	user := prefix + hex.EncodeToString(sum[:])

	if ok := m.SetIfAbsent(user, ws); !ok {
		return
	}
	defer m.Remove(user)

	add(user)
	defer del(user)

	atomic.AddInt32(&c, 1)
	defer atomic.AddInt32(&c, -1)

	for {
		e := WS_PING.Receive(ws, nil)
		if e != nil {
			return
		}
	}
}
