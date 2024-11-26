package app

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spv-dev/platform_common/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/spv-dev/chat-server/internal/config"
	"github.com/spv-dev/chat-server/internal/interceptor"
	"github.com/spv-dev/chat-server/internal/logger"
	"github.com/spv-dev/chat-server/internal/metric"
	desc "github.com/spv-dev/chat-server/pkg/chat_v1"
)

// App структура приложения
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp инициализизует зависимости и создаёт экземпляр структуры приложения
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает gRPC сервер
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run gRPC server")
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runPrometheus()
		if err != nil {
			log.Fatalf("failed to run Prometheus server")
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.MetricsInterceptor,
				//a.serviceProvider.AccessInterceptor(ctx).AccessInterceptor,
			),
		),
	)

	reflection.Register(a.grpcServer)

	desc.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatServer(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is runnign on %v", a.serviceProvider.GRPCConfig().Address())

	flag.Parse()

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = metric.DefaultInit()
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}
	logger.DefaultInit()

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:              "localhost:2112",
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}

	log.Printf("Prometheus server is running on %v", "localhost:2112")

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
