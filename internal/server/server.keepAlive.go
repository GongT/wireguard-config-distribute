package server

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverImplement) KeepAlive(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	fmt.Println("Call to KeepAlive")
	return &emptypb.Empty{}, nil
}
