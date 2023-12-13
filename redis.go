package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var pubsub *redis.PubSub
var rdb *redis.Client

func initRedis(ctx context.Context) {
	//redis://user:password@localhost:6789/3
	redisConnectStr := ""
	if tlsBool {
		redisConnectStr += "rediss://"
	} else {
		redisConnectStr += "redis://"
	}
	redisConnectStr += password + "@" + addr + "/" + db
	option, err := redis.ParseURL(redisConnectStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	rdb = redis.NewClient(option)
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
