package grpcImplements

import (
	"context"
	"errors"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"golang.org/x/crypto/bcrypt"
)

func (s *serverImplement) GetSelfSignedCertFile(_ context.Context, req *protocol.GetCertFileRequest) (*protocol.GetCertFileResponse, error) {
	if bcrypt.CompareHashAndPassword(req.Password, []byte(s.password)) != nil {
		return nil, errors.New("Password wrong!")
	}

	if s.insecure {
		return nil, errors.New("Server is insecure mode")
	}

	return &protocol.GetCertFileResponse{
		CertFileText: s.storage.GetCaCertFileContent(),
	}, nil
}
