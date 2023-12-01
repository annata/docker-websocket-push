package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
)

func pushRoute(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	body, err := io.ReadAll(request.Body)
	for i := 0; i < 64; i++ {
		token := request.URL.Query().Get("token" + strconv.Itoa(i))
		if token == "" {
			if i == 0 {
				token = request.URL.Query().Get("token")
				if token == "" {
					break
				}
			} else {
				break
			}
		}
		sum := md5.Sum([]byte(token))
		topic := prefix + hex.EncodeToString(sum[:])
		if err == nil {
			go rdb.Publish(ctx, topic, body)
		}
	}
	header := response.Header()
	header.Set("Content-Type", "application/json;charset=UTF-8")
	header.Set("Access-Control-Allow-Origin", "*")

	response.Write([]byte("{\"code\":\"0\"}"))
}
