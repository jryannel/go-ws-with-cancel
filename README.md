# Go WebSocket with Cancel

A websocket client/server model with cancel support.

This is my learning attempt to create a websocket client/server model with cancel support. I have used the following resources to learn and create this project:

* %[https://github.com/gorilla/websocket]
* %[https://pkg.go.dev/context]
* %[https://www.youtube.com/watch?v=LSzR0VEraWw]
* %[https://www.amazon.com/Concurrency-Go-Tools-Techniques-Developers/dp/1491941197]

## Concept

* `Connection` is the websocket connection used in client and server hub
* `Hub` is the websocket server hub, manages all connections and creates a request stream
* `Request` is the websocket request, contains message and possible error. Allows to reply to sender or to broadcast to all connections

## Usage

The demo will send a message with a count property and increment the value by 1 on the server and then on the client side. The client will stop when the count reaches 100000.

You can play with the max count and timeout settings to see which is triggered first.

### Server

```go run cmd/server/main.go```

### Client

```go run cmd/client/main.go```