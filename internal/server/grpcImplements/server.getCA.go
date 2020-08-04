package grpcImplements

import (
	"context"
	"errors"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
)

func (s *Implements) GetSelfSignedCertFile(ctx context.Context, req *protocol.GetCertFileRequest) (*protocol.GetCertFileResponse, error) {
	if s.insecure {
		return nil, errors.New("Server is insecure mode")
	}

	return &protocol.GetCertFileResponse{
		CertFileText: s.storage.GetCaCertFileContent(),
	}, nil
}
