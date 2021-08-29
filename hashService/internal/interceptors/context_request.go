package interceptors

import (
	"context"
	guid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

func ContextRequestInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestId, ok := ctx.Value("requestID").(string)
		if !ok || requestId == "" {
			uuid := guid.Must(guid.NewV4(), nil)
			requestId = uuid.String()
		}

		return handler(context.WithValue(ctx, "requestID", requestId), req)
	}
}
