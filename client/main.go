package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "gopoll/proto"
)

var client struct {
}

func chat(client pb.ChatServiceClient, from string, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, _ := client.Chat(ctx)

	stream.Send(&pb.ChatMessage{From: from, Message: msg})
	stream.CloseSend()
}

// wait for input, then send as gRPC message
func main() {
	var sender = make(chan string)

	// connect to gRPC server
	conn, _ := grpc.Dial(":6453", grpc.WithInsecure())
	client := pb.NewChatServiceClient(conn)

	// get the stream
	stream, _ := client.Chat(context.Background())

	// print messages coming in on the stream
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a message : %v", err)
				continue
			}
			fmt.Printf("BROADCAST: \t\t\t %v\n", in.GetMessage())
			fmt.Print("-> ")
		}
	}()

	// take input from the user and send to sender channel
	go func() {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Chat client listening for input...")

		for {
			fmt.Print("-> ")
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("* failed to parse input, try again *")
				continue
			}

			input = strings.Replace(input, "\n", "", -1)
			if len(input) < 1 {
				fmt.Println("\t\t\t" + "(no content)" + " <-")
				continue
			}
			sender <- input
		}
	}()

	// get messages from sender channel and send to server
	for {
		rawmessage := <-sender
		stream.Send(&pb.ChatMessage{
			From:    "Christopher Wong",
			Message: rawmessage,
		})
	}

	stream.CloseSend()
	<-waitc
}
