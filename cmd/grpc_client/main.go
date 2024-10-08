package main

import (
	"context"
	"log"
	"time"

	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	host = "localhost:50051"
	id   = 13
)

func main() {
	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connection to host %s . error: %s", host, err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("error when close connection: %s", err)
		}
	}()

	c := desc.NewChatV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	respCreate, err := c.Create(ctx, &desc.CreateRequest{
		Info: &desc.ChatInfo{
			Usernames: []string{"a", "b", "c"},
		},
	})

	if err != nil {
		log.Fatalf("failed when create chat %s", err)
	}
	log.Printf("Created chat: \n%+v", respCreate.GetId())

	respDelete, err := c.Delete(ctx, &desc.DeleteRequest{
		Id: id,
	})

	if err != nil {
		log.Fatalf("failed when delete chat %s", err)
	}
	log.Printf("Deleted chat: \n%+v", respDelete)

	respSendMessage, err := c.SendMessage(ctx, &desc.SendMessageRequest{
		Message: &desc.Message{
			Info: &desc.MessageInfo{
				From: "Max",
				Text: "Hello, my friend!",
			},
			Timestamp: timestamppb.Now(),
		},
	})

	if err != nil {
		log.Fatalf("failed when send message %s", err)
	}
	log.Printf("Send Message: \n%+v", respSendMessage)
}
