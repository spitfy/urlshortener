package handler

import (
	"context"
	"github.com/spitfy/urlshortener/internal/auth"
	"net"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/pkg/shortener"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer — структура для gRPC-сервера
type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
}

// NewGRPCServer создает и настраивает gRPC-сервер
func NewGRPCServer(cfg config.Config, service ServiceShortener, auth *auth.Manager) (*GRPCServer, error) {
	grpcServer := grpc.NewServer()

	shortener.RegisterShortenerServiceServer(grpcServer, newGRPC(service, auth))

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.Handlers.GRPCAddr)
	if err != nil {
		return nil, err
	}

	return &GRPCServer{
		server:   grpcServer,
		listener: listener,
	}, nil
}

// Serve запускает gRPC-сервер
func (g *GRPCServer) Serve() error {
	return g.server.Serve(g.listener)
}

// GracefulStop останавливает gRPC-сервер
func (g *GRPCServer) GracefulStop() {
	g.server.GracefulStop()
}

// Shutdown останавливает сервер с таймаутом
func (g *GRPCServer) Shutdown(ctx context.Context) error {
	done := make(chan bool, 1)
	go func() {
		g.server.GracefulStop()
		done <- true
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
