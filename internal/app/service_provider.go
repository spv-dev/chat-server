package app

import (
	"context"
	"log"

	"github.com/spv-dev/chat-server/internal/api/chat"
	"github.com/spv-dev/chat-server/internal/client/db"
	"github.com/spv-dev/chat-server/internal/client/db/pg"
	"github.com/spv-dev/chat-server/internal/client/db/transaction"
	"github.com/spv-dev/chat-server/internal/closer"
	"github.com/spv-dev/chat-server/internal/config"
	"github.com/spv-dev/chat-server/internal/repository"
	chatRepository "github.com/spv-dev/chat-server/internal/repository/chat"
	"github.com/spv-dev/chat-server/internal/service"
	chatService "github.com/spv-dev/chat-server/internal/service/chat"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repository.ChatRepository

	chatService service.ChatService

	chatServer *chat.Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}

		s.pgConfig = cfg
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}

		s.grpcConfig = cfg
	}
	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database : %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err)
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) ChatServer(ctx context.Context) *chat.Server {
	if s.chatServer == nil {
		s.chatServer = chat.NewServer(s.ChatService(ctx))
	}

	return s.chatServer
}
