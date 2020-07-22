package server

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serverImplement) Close(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
