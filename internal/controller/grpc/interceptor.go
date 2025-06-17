package grpc_handler

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()

		resp, err = handler(ctx, req)

		duration := time.Since(start)
		logger.Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.String("status", status.Code(err).String()),
			zap.Duration("duration", duration),
			zap.Error(err),
		)

		return resp, err
	}
}
