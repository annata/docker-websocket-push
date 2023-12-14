package main

import (
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"net"
)

var pubsub *redis.PubSub
var rdb *redis.Client

func initRedis(ctx context.Context) {
	//redis://user:password@localhost:6789/3
	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		h = addr
	}
	if h == "" {
		h = "localhost"
	}
	if p == "" {
		p = "6379"
	}
	var tlsConfig *tls.Config = nil
	if tlsBool {
		tlsConfig = &tls.Config{ServerName: h}
	}
	rdb = redis.NewClient(&redis.Options{
		Network:   "tcp",
		Addr:      net.JoinHostPort(h, p),
		Password:  password,
		DB:        db,
		TLSConfig: tlsConfig,
	})
	pubsub = rdb.Subscribe(ctx, prefix+"all")
}

func connectRedis(ctx context.Context) {
	defer pubsub.Close()
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			cancel()
			return
		}

		if msg.Channel == prefix+"all" {
			go sendAll(m, msg.Payload)
		} else {
			go send(msg.Channel, msg.Payload)
		}
	}
}

func add(user string) error {
	return pubsub.Subscribe(ctx, user)
}

func del(user string) error {
	return pubsub.Unsubscribe(ctx, user)
}
