package grpcImplements

import (
	"context"
	"errors"
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/peer"
)

func (s *serverImplement) GetSelfSignedCertFile(ctx context.Context, req *protocol.GetCertFileRequest) (*protocol.GetCertFileResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("Failed get peer info")
	}

	if bcrypt.CompareHashAndPassword(req.Password, []byte(s.password)) != nil {
		if !p.Addr.(*net.TCPAddr).IP.IsLoopback() {
			return nil, errors.New("Password wrong!")
		}
	}

	if s.insecure {
		return nil, errors.New("Server is insecure mode")
	}

	return &protocol.GetCertFileResponse{
		CertFileText: s.storage.GetCaCertFileContent(),
	}, nil
}
