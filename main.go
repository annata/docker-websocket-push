package main

import (
	"context"
	"flag"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var addr = ""
var password = ""
var db = 0
var port = ""
var prefix = "ws_push."
var ctx context.Context
var cancel context.CancelFunc

func main() {
	ctx, cancel = context.WithCancel(context.Background())
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	parse()
	initRedis(ctx)
	go connectRedis(ctx)
	http.Handle("/ws", websocket.Handler(websocketHandle))
	http.HandleFunc("/", defaultRoute)
	server := &http.Server{Addr: ":" + port, Handler: nil}
	go stopHttp(server)
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

func stopHttp(server *http.Server) {
	<-ctx.Done()
	server.Shutdown(context.TODO())
}

func parse() {
	flag.StringVar(&addr, "addr", "localhost:6379", "redis连接地址")
	flag.StringVar(&password, "password", "", "redis密码")
	flag.StringVar(&port, "port", "8080", "端口")
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
	portStr := os.Getenv("port")
	if portStr != "" {
		port = portStr
	}
}

func defaultRoute(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json;charset=UTF-8")
	response.Write([]byte("{\"code\":\"0\"}"))
}
