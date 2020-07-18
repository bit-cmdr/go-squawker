package broker

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGenID(t *testing.T) {
	counter = 0
	tests := []struct {
		name string
		want int
	}{
		{
			name: "first",
			want: 1,
		},
		{
			name: "second",
			want: 2,
		},
		{
			name: "third",
			want: 3,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := GenID(); got != tt.want {
				t.Errorf("GenID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	counter = 0
	type args struct {
		conn *websocket.Conn
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{
			name: "first",
			args: args{conn: nil},
			want: Client{
				ID:   1,
				conn: nil,
			},
		},
		{
			name: "second",
			args: args{conn: nil},
			want: Client{
				ID:   2,
				conn: nil,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.conn); !reflect.DeepEqual(got.ID, tt.want.ID) {
				t.Errorf("NewClient() = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

func TestBus_AddClient(t *testing.T) {
	var clients []Client
	type args struct {
		c Client
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "add works as expected for adding one",
			args: args{c: Client{ID: 1, conn: nil}},
		},
		{
			name: "add works as expected for adding a second one",
			args: args{c: Client{ID: 2, conn: nil}},
		},
	}
	b := &Bus{
		clients:  clients,
		messages: []MessageData{},
	}
	for i, tt := range tests {
		tt := tt
		i := i
		b := b
		t.Run(tt.name, func(t *testing.T) {
			b.AddClient(tt.args.c)
			if len(b.clients) != i+1 {
				t.Errorf("AddClient() = %v, want %v", len(b.clients), i+1)
			}
			if b.clients[i].ID != tt.args.c.ID {
				t.Errorf("AddClient() = %v, want %v", b.clients[i].ID, tt.args.c.ID)
			}
		})
	}
}

func TestBus_Publish(t *testing.T) {
	var upgrader = websocket.Upgrader{}

	echo := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	}
	server := httptest.NewServer(http.HandlerFunc(echo))
	defer server.Close()
	uri := "ws" + strings.TrimPrefix(server.URL, "http")
	getSocket := func() *websocket.Conn {
		ws, _, err := websocket.DefaultDialer.Dial(uri, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		return ws
	}
	var messages []MessageData
	var clients []Client
	b := &Bus{
		messages: messages,
		clients:  clients,
	}
	type args struct {
		md MessageData
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "publishes messages and saves them in memory",
			args: args{
				md: MessageData{websocket.TextMessage, []byte("test")},
			},
		},
	}
	for i, tt := range tests {
		i := i
		b := b
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ws := getSocket()
			defer ws.Close()
			b.AddClient(NewClient(ws))
			b.Publish(tt.args.md)
			if err := ws.WriteMessage(tt.args.md.MessageType, tt.args.md.Payload); err != nil {
				t.Errorf("%v", err)
			}
			_, p, err := ws.ReadMessage()
			if err != nil {
				t.Errorf("%v", err)
			}
			if string(p) != "test" {
				t.Errorf("Publish() sent %v, want \"test\"", string(tt.args.md.Payload))
			}
			if len(b.messages) != i+1 {
				t.Errorf("Publish() len %v, want %v", len(b.messages), i+1)
			}
		})
	}
}
