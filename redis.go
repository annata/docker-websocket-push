package main

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var pubsub *redis.PubSub
var rdb *redis.Client

func initRedis(ctx context.Context) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
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
