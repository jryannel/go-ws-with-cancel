package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wswc/log"
	"wswc/model"
	"wswc/ws"
)

const (
	maxCount = 100000
	timeout  = 10 * time.Second
)

var addr = flag.String("addr", "ws://127.0.0.1:8080", "ws service address")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := ws.Dial(ctx, *addr)
	if err != nil {
		cancel()
		panic(err)
	}

	// handle incoming messages
	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				log.Debugln("waiting for message")
				var msg model.Message
				err := client.ReadJSON(&msg)
				if err != nil {
					log.Warnf("read message failed: %v", err)
					return
				}
				log.Debugf("count: %d\n", msg.Count)
				if msg.Count >= maxCount {
					log.Infof("client received %d messages, exit\n", maxCount)
					return
				}
				msg.Count++
				client.WriteJSON(&msg)
				if err != nil {
					log.Errorf("write message failed: %v", err)
					return
				}
			}
		}
	}()

	// send initial message
	msg := model.Message{Count: 0}
	client.WriteJSON(msg)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		<-sigs
		cancel()
	}()
	time.AfterFunc(timeout, func() {
		log.Infof("timeout, ...")
		cancel()
	})
	<-ctx.Done()
	log.Infoln("client exit")

}
