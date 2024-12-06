package main

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/net/websocket"
	"strconv"
)

func parseToken(ws *websocket.Conn) map[string]int {
	m := make(map[string]int)
	for i := 0; i < 64; i++ {
		token := ws.Request().URL.Query().Get("token" + strconv.Itoa(i))
		if token == "" {
			token = ws.Request().Header.Get("token" + strconv.Itoa(i))
			if token == "" {
				if i == 0 {
					token = ws.Request().URL.Query().Get("token")
					if token == "" {
						token = ws.Request().Header.Get("token")
						if token == "" {
							break
						}
					}
				} else {
					break
				}
			}
		}
		sum := md5.Sum([]byte(token))
		topic := prefix + hex.EncodeToString(sum[:])
		m[topic] = 0
	}
	return m
}

func wsConnect(sn string, closeFlag <-chan any, ws *websocket.Conn) {
	//sn := strconv.FormatUint(atomic.AddUint64(&snn, 1), 10)
	for i := 0; i < 64; i++ {
		token := ws.Request().URL.Query().Get("token" + strconv.Itoa(i))
		if token == "" {
			token = ws.Request().Header.Get("token" + strconv.Itoa(i))
			if token == "" {
				if i == 0 {
					token = ws.Request().URL.Query().Get("token")
					if token == "" {
						token = ws.Request().Header.Get("token")
						if token == "" {
							break
						}
					}
				} else {
					break
				}
			}
		}
		sum := md5.Sum([]byte(token))
		topic := prefix + hex.EncodeToString(sum[:])
		go topicMap(closeFlag, sn, topic, ws)
	}
}

func globalMap(closeFlag <-chan any, sn string, ws *websocket.Conn) {
	m.Set(sn, ws)
	defer m.Remove(sn)
	<-closeFlag
}

func topicMap(closeFlag <-chan any, sn, topic string, ws *websocket.Conn) {
	addTopic(topic, sn, ws)
	defer removeTopic(topic, sn)
	<-closeFlag
}
