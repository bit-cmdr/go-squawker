# Squawker

A quick in-memory message sharing/broadcasting platform in Go using web sockets and static html from the server

## Usage

Get the modules

```sh
go mod download
```

Run the server

```sh
go run .
```

Open a few browser tabs

uri | purpose | abilities
---|---|---
http://localhost:3000/squawker | Server WebSocket Port | Brokers Websocket messages
http://localhost:3000/listen | Lists messages as they come in and previous ones that have been sent | Listen Only
http://localhost:3000/spam | Sends a single message or spams a message n times | Sends Only
http://localhost:3000/squawk | Lists messages as they come in and previous ones that have been sent, also can send a single message to all clients | Send and Listen

## Testing

```sh
go test ./...
```

### Test Coverage

```sh
go test -coverprofile=cover.out ./...
go tool cover -func=cover.out
go tool cover -html=cover.out
```
