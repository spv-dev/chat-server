package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/spv-dev/chat-server/database"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedChatV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})
	log.Printf("server listening at %s", lis.Addr())

	//инициализация базы данных
	database.InitDB()
	defer database.CloseDB()

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}

// CreateChat creates a new chat
//
// Return id of created chat
func (s *server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	return database.CreateChatDB(ctx, req)
}

// DeleteChat deletes the chat by id
func (s *server) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	return database.DeleteChatDB(ctx, req)
}

// SendMessage sends the user's message to chat
func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	return database.SendMessageDB(ctx, req)
}

// GetChatMessages получает все сообщения чата
func (s *server) GetChatMessages(ctx context.Context, req *desc.GetChatMessagesRequest) (*desc.GetChatMessagesResponse, error) {
	return database.GetChatMessagesDB(ctx, req)
}
