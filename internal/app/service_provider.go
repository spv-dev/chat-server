package app

import (
	"context"
	"log"

	"github.com/spv-dev/platform_common/pkg/closer"
	"github.com/spv-dev/platform_common/pkg/db"
	"github.com/spv-dev/platform_common/pkg/db/pg"
	"github.com/spv-dev/platform_common/pkg/db/transaction"

	"github.com/spv-dev/chat-server/internal/api/chat"
	"github.com/spv-dev/chat-server/internal/config"
	"github.com/spv-dev/chat-server/internal/interceptor"
	"github.com/spv-dev/chat-server/internal/repository"
	chatRepository "github.com/spv-dev/chat-server/internal/repository/chat"
	"github.com/spv-dev/chat-server/internal/service"
	chatService "github.com/spv-dev/chat-server/internal/service/chat"
)

type serviceProvider struct {
	pgConfig         config.PGConfig
	grpcConfig       config.GRPCConfig
	prometheusConfig config.PrometheusConfig

	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repository.ChatRepository

	chatService service.ChatService

	chatServer *chat.Server

	accessInterceptor *interceptor.AccessInterceptor
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// PGConfig получение конфигурации подключения к Постгрессу
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

// GRPCConfig получение конфигурации для gRPC
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

// PrometheusConfig получение конфигурации для Prometheus
func (s *serviceProvider) PrometheusConfig() config.PrometheusConfig {
	if s.prometheusConfig == nil {
		cfg, err := config.NewPrometheusConfig()
		if err != nil {
			log.Fatalf("failed to get prometheus config: %v", err)
		}

		s.prometheusConfig = cfg
	}
	return s.prometheusConfig
}

// DBClient получение клиента БД
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

// TxManager получение менеджера транзакций
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

// ChatRepository получение слоя repo
func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

// ChatService получение сервисного слоя
func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

// ChatServer получение указателя на сам сервис
func (s *serviceProvider) ChatServer(ctx context.Context) *chat.Server {
	if s.chatServer == nil {
		s.chatServer = chat.NewServer(s.ChatService(ctx))
	}

	return s.chatServer
}

func (s *serviceProvider) AccessInterceptor(_ context.Context) *interceptor.AccessInterceptor {
	if s.accessInterceptor == nil {
		s.accessInterceptor = interceptor.NewAccessInterceptor()
	}

	return s.accessInterceptor
}
