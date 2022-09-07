# Go WebSocket with Cancel

A websocket client/server model with cancel support.

This is my learning attempt to create a websocket client/server model with cancel support. I have used the following resources to learn and create this project:

* https://github.com/gorilla/websocket
* https://pkg.go.dev/context
* https://www.youtube.com/watch?v=LSzR0VEraWw

## Usage

The demo will send a message with a count property and increment the value by 1 on the server and then on the client side. The client will stop when the count reaches 10000.

### Server

```go run cmd/server/main.go```

### Client

```go run cmd/client/main.go```