package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

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

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

// CreateChat creates a new chat
//
// Return id of created chat
func (s *server) CreateChat(ctx context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	log.Printf("Method: %s\nRequest: %v\nContext: %v\n", "CreateChat", req, ctx)
	return &desc.CreateChatResponse{}, nil
}

// DeleteChat deletes the chat by id
func (s *server) DeleteChat(ctx context.Context, req *desc.DeleteChatRequest) (*emptypb.Empty, error) {
	log.Printf("Method: %s\nRequest: %v\nContext: %v\n", "DeleteChat", req, ctx)
	return nil, nil
}

// SendMessage sends the user's message to chat
func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Method: %s\nRequest: %v\nContext: %v\n", "SendMessage", req, ctx)
	return nil, nil
}
