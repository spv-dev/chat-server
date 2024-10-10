package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
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
			log.Fatalf("error to close connection: %s", err)
		}
	}()

	c := desc.NewChatV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	respCreate, err := c.CreateChat(ctx, &desc.CreateChatRequest{
		Info: &desc.ChatInfo{
			Name:     "First Chat",
			UsersIds: []int64{5, 100, 1200},
		},
	})
	if err != nil {
		log.Fatalf("failed to create chat %s", err)
	}
	log.Printf("Created chat: \n%+v", respCreate.GetId())

	respDelete, err := c.DeleteChat(ctx, &desc.DeleteChatRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("failed to delete chat %s", err)
	}
	log.Printf("Deleted chat: \n%+v", respDelete)

	respSendMessage, err := c.SendMessage(ctx, &desc.SendMessageRequest{
		Info: &desc.MessageInfo{
			From:      100,
			Text:      "Hello, my friend!",
			Timestamp: timestamppb.Now(),
		},
	})
	if err != nil {
		log.Fatalf("failed to send message %s", err)
	}
	log.Printf("Send Message: \n%+v", respSendMessage)
}
