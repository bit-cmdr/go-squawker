package broker

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   int
	conn *websocket.Conn
}

type MessageData struct {
	MessageType int
	Payload     []byte
}

type Bus struct {
	clients  []Client
	messages []MessageData
}

var counter int
var counterMu sync.Mutex

func GenID() int {
	counterMu.Lock()
	defer counterMu.Unlock()
	counter++
	return counter
}

func NewClient(conn *websocket.Conn) Client {
	return Client{ID: GenID(), conn: conn}
}

func (c *Client) Send(md MessageData) error {
	if err := c.conn.WriteMessage(md.MessageType, md.Payload); err != nil {
		log.Println("Unable to write message", err)
		return err
	}
	return nil
}

func (b *Bus) AddClient(c Client) {
	log.Println("sending messages to new client", len(b.messages))
	for _, m := range b.messages {
		if err := c.conn.WriteMessage(m.MessageType, m.Payload); err != nil {
			log.Println("Unable to rebroadcast to client", c.ID)
			return
		}
	}
	b.clients = append(b.clients, c)
}

func (b *Bus) RemoveClient(client Client) {
	for i, c := range b.clients {
		if c.ID == client.ID {
			b.clients = append(b.clients[:i], b.clients[i+1:]...)
		}
	}
}

func (b *Bus) Publish(md MessageData) {
	b.messages = append(b.messages, md)
	for _, c := range b.clients {
		if err := c.Send(md); err != nil {
			b.RemoveClient(c)
		}
	}
}
