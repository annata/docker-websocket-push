package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"os"
)

var pubsub *redis.PubSub
var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	pubsub = rdb.Subscribe(context.TODO(), prefix+"all")
}

func connectRedis() {
	for {
		msg, err := pubsub.ReceiveMessage(context.TODO())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
			return
		}

		if msg.Channel == prefix+"all" {
			go sendAll(m, msg.Payload)
		} else {
			go send(msg.Channel, msg.Payload)
		}
	}
}

func add(user string) {
	pubsub.Subscribe(context.TODO(), user)
}

func del(user string) {
	pubsub.Unsubscribe(context.TODO(), user)
}
