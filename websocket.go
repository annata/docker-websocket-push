package main

import (
	"github.com/orcaman/concurrent-map/v2"
	"golang.org/x/net/websocket"
	"sync/atomic"
	"time"
)

var m = cmap.New[*websocket.Conn]()

var n = cmap.New[*WsMap]()

type WsMap struct {
	mm    cmap.ConcurrentMap[string, *websocket.Conn]
	count int32
}

var snn uint64 = 0

func send(key string, value string) {
	mm, ok := n.Get(key)
	if ok {
		sendAll(mm.mm, value)
	}
}

func sendAll(mm cmap.ConcurrentMap[string, *websocket.Conn], value string) {
	tuple := mm.IterBuffered()
	number := (cap(tuple) / 256) + 1
	for i := 0; i < number; i++ {
		go func() {
			for t := range tuple {
				ws := t.Val
				go func() {
					err := websocket.Message.Send(ws, value)
					if err != nil {
						ws.Close()
					}
				}()
			}
		}()
	}
}

func websocketHandle(ws *websocket.Conn) {
	defer ws.Close()
	closeFlag := make(chan any)
	defer close(closeFlag)
	go wsConnect(closeFlag, ws)
	go ping(closeFlag, ws)
	for {
		e := WsPing.Receive(ws, nil)
		if e != nil {
			return
		}
	}
}

func addTopic(topic, sn string, ws *websocket.Conn) {
	var mm *WsMap
	var ok bool
	shard := n.GetShard(topic)
	shard.RLock()
	for mm, ok = n.Get(topic); !ok; mm, ok = n.Get(topic) {
		shard.RUnlock()
		mm = &WsMap{
			mm:    cmap.New[*websocket.Conn](),
			count: 0,
		}
		n.Upsert(topic, mm, func(exist bool, valueInMap, newValue *WsMap) *WsMap {
			if exist {
				return valueInMap
			} else {
				add(topic)
				return newValue
			}
		})
		shard.RLock()
	}
	setRes := mm.mm.SetIfAbsent(sn, ws)
	if setRes {
		atomic.AddInt32(&mm.count, 1)
	}
	shard.RUnlock()
}

func removeTopic(topic, sn string) {
	mm, ok := n.Get(topic)
	if ok {
		exist := mm.mm.RemoveCb(sn, func(key string, ws *websocket.Conn, exists bool) bool {
			return exists
		})
		if exist {
			number := atomic.AddInt32(&mm.count, -1)
			if number == 0 {
				go delTopic(topic)
			}
		}
	}
}

func delTopic(topic string) {
	n.RemoveCb(topic, func(key string, v *WsMap, exists bool) bool {
		if exists && atomic.LoadInt32(&v.count) == 0 {
			del(topic)
			return true
		} else {
			return false
		}
	})
}

func ping(closeFlag <-chan any, ws *websocket.Conn) {
	ticker := time.NewTicker(40 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-closeFlag:
			return
		case <-ticker.C:
			err := WsPing.Send(ws, nil)
			if err != nil {
				ws.Close()
				return
			}
		}
	}
}
