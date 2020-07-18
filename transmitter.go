package main

import (
	"log"
	"net/http"

	"github.com/bit-cmdr/go-squawker/pkg/broker"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var messageBus = &broker.Bus{}
var broadcast = make(chan broker.MessageData)

func transmitter(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := broker.NewClient(conn)
	log.Println("Client connected", client.ID)
	messageBus.AddClient(client)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Unable to read message", err)
			messageBus.RemoveClient(client)
			return
		}
		broadcast <- broker.MessageData{MessageType: messageType, Payload: p}
	}
}

func broadcastMessages() {
	for {
		md := <-broadcast
		messageBus.Publish(md)
	}
}
