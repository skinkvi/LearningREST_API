package grpc

import (
	"context"
	"log/slog"
	"rest_api_app/internal/lib/logger/sl"
	"time"

	ssov1 "github.com/skinkvi/protosSTT/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	api ssov1.AuthClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "client.sso.grpc.New"

	log.Info("creating grpc client", slog.String("addr", addr), slog.Duration("timeout", timeout))

	creds, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		log.Error("failed to create TLS credentials", sl.Err(err))
		return nil, err
	}

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(creds), grpc.WithBlock(), grpc.WithTimeout(timeout))
	if err != nil {
		log.Error("Failed to dial grpc server", slog.String("addr", addr), op, sl.Err(err))
		return nil, err
	}

	client := &Client{
		api: ssov1.NewAuthClient(conn),
		log: log,
	}

	return client, nil
}
