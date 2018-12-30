package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/golang/protobuf/ptypes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "gopoll/proto"
)

var broadcastChan = make(chan *pb.ChatMessage)
var clients = make(map[pb.ChatService_ChatServer]bool)

type chatServer struct {
}

func broadcast() {
	for {
		m := <-broadcastChan
		fmt.Printf("broadcasting: %v", m)
		for client := range clients {
			if err := client.Send(&pb.BroadcastMessage{
				Timestamp: ptypes.TimestampNow(),
				Message:   m,
			}); err != nil {
				fmt.Println("failed to broadcast message")
			}
		}
	}
}

func (s *chatServer) Chat(stream pb.ChatService_ChatServer) error {
	for {
		clients[stream] = true

		m, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		broadcastChan <- m
	}
}

func run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 6453))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, &chatServer{})
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	fmt.Println("listening on", lis.Addr())
	grpcServer.Serve(lis)
}

func main() {

	go func() {
		PORT := 3000
		fs := http.FileServer(http.Dir("static"))
		http.Handle("/", fs)

		log.Println(fmt.Sprintf("web app listening on localhost:%d", PORT))
		http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	}()

	go run()       // start the gRPC server
	go broadcast() // print every message coming in on the broadcast channel then send out to clients

	// stop main from exiting while goroutine is executing
	select {}
}
