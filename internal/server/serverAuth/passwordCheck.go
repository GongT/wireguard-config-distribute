package serverAuth

import (
	"context"
	"errors"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type passwordCheck struct {
	password string
}

// Return value is mapped to request headers.
func (t *passwordCheck) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authorize(t.password, stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}
func (t *passwordCheck) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authorize(t.password, ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func CreatePasswordCheck(password string) *passwordCheck {
	return &passwordCheck{password}
}

func check(input string, password string, salt string) bool {
	input = strings.TrimPrefix(input, "Bearer ")
	err := bcrypt.CompareHashAndPassword([]byte(input), []byte(password+salt))
	// tools.Debug("check: %v", err)
	return err == nil
}

func authorize(password string, ctx context.Context) error {
	if md, ok := metadata.FromIncomingContext(ctx); !ok {
		return errors.New("Failed get metadata from incoming context")
	} else if ip := tools.GetRemoteFromContext(ctx); len(ip) == 0 {
		return errors.New("Failed get peer info from incoming context")
	} else if len(md["authorization"]) > 0 && len(md["authorization-salt"]) > 0 {
		if check(md["authorization"][0], password, md["authorization-salt"][0]) {
			return nil
		} else {
			tools.Error("auth fail: from %s <%v:%v>", ip, md["authorization-salt"][0], md["authorization"][0])
			return errors.New("Invalid authorization data")
		}
	} else {
		tools.Error("auth fail: from %s <empty>", ip)

		return errors.New("Empty authorization")
	}
}
