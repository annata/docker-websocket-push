package main

import (
	"flag"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"strconv"
)

var addr = ""
var password = ""
var db = 0
var prefix = "ws_push."

func main() {
	parse()
	initRedis()
	go connectRedis()
	http.Handle("/ws", websocket.Handler(websocketHandle))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		os.Exit(0)
	}
}

func parse() {
	flag.StringVar(&addr, "addr", "localhost:6379", "redis连接地址")
	flag.StringVar(&password, "password", "", "redis密码")
	flag.IntVar(&db, "db", 0, "redis数据库")
	flag.Parse()
	addrStr := os.Getenv("addr")
	if addrStr != "" {
		addr = addrStr
	}
	passwordStr := os.Getenv("password")
	if passwordStr != "" {
		password = passwordStr
	}
	dbStr := os.Getenv("db")
	if dbStr != "" {
		dbInt, e := strconv.Atoi(dbStr)
		if e == nil {
			db = dbInt
		}
	}
}
