package main

import (
	"log"
	"net/http"
)

const port = ":3000"

func main() {
	go broadcastMessages()
	http.HandleFunc("/spam", spam)
	http.HandleFunc("/listen", listen)
	http.HandleFunc("/squawk", squawk)
	http.HandleFunc("/squawker", transmitter)

	log.Println("Server is running on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Println("Unable start server", err)
	}
}
