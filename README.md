
# gopoll

current issues:

```none
* dockerize
* add etcd as a message store
```

attempting to build a real time polling service using gRPC and bi-directional streaming. Purposely not using Websockets.

Once we get basic messaging working (simple live chat), we can generalize for other payload types.

Use case: Think HQ Trivia, but built in Go.

## Run locally

### start the server

```bash
go run server/main.go
```

the server prints all messages to stdout and broadcasts messages to all connected clients

### start the client

```bash
go run client/main.go
```

the client dials the gRPC server and connects to the stream, printing messages coming in and sending new messages sent from a `sender` message channel.

the client also starts an input loop reading text from stdin and sending in to the `sender` channel to be broadcast.