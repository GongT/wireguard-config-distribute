package client

import (
	"google.golang.org/grpc"
)

type Client struct {
	grpcConn *grpc.ClientConn
}
