package http

import "context"

type Http interface {
	Start(port string)
	Stop(ctx context.Context)
}
