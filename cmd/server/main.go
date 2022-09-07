package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"wswc/log"
	"wswc/model"
	"wswc/ws"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h := ws.NewHub(ctx)
	go func() {
		for r := range h.Requests() {
			var msg model.Message
			err := r.AsJSON(&msg)
			if err != nil {
				log.Warnf("unmarshal message failed: %v", err)
				continue
			}
			log.Debugf("server received: %d", msg.Count)
			msg.Count++
			log.Debugf("server send: %d", msg.Count)
			r.ReplyJSON(&msg)
		}
	}()
	go h.Run()
	http.HandleFunc("/", h.ServeHTTP)
	fmt.Printf("listen at %s\n", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		panic(err)
	}
}
