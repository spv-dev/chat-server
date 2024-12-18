package interceptor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	accessDesc "github.com/spv-dev/auth/pkg/access_v1"
	"github.com/spv-dev/chat-server/internal/serviceerror"
)

const (
	accessPrefix = "Bearer "
	accessPort   = "50053"
)

// AccessInterceptor интерцептор доступа
type AccessInterceptor struct {
	AccessV1Client accessDesc.AccessV1Client
}

// AccessInterceptor функция для проверки доступа
func (accessInterceptor AccessInterceptor) AccessInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New(serviceerror.ErrorInMetadata)
	}

	accessHeader, ok := md["authorization"]
	if !ok || len(accessHeader) == 0 {
		return nil, errors.New(serviceerror.ErrorInHeader)
	}

	if !strings.HasPrefix(accessHeader[0], accessPrefix) {
		return nil, errors.New(serviceerror.FailedInBearer)
	}

	clientCtx := context.Background()
	clientCtx = metadata.NewOutgoingContext(clientCtx, md)

	_, err := accessInterceptor.AccessV1Client.Check(clientCtx, &accessDesc.CheckRequest{
		Endpoint: info.FullMethod,
	})
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

// NewAccessInterceptor создание интерцептора доступа
func NewAccessInterceptor() *AccessInterceptor {
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%s", accessPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("error create connect to access service")
	}

	client := accessDesc.NewAccessV1Client(conn)
	return &AccessInterceptor{
		AccessV1Client: client,
	}
}
