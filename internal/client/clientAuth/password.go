package clientAuth

import (
	"context"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"golang.org/x/crypto/bcrypt"
)

type passwordAuth struct {
	password string
}

func hashPassword(password string) string {
	rawPasswd, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		tools.Die("Failed encrypt password: %s", err.Error())
	}
	return string(rawPasswd)
}

// Return value is mapped to request headers.
func (t *passwordAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	salt := tools.RandString(5)
	// tools.Debug("  * request salt: %s", salt)
	return map[string]string{
		"authorization":      "Bearer " + hashPassword(t.password+salt),
		"authorization-salt": salt,
	}, nil
}

func (*passwordAuth) RequireTransportSecurity() bool {
	return true
}

func CreatePasswordAuth(password string) *passwordAuth {
	// tools.Debug("using password %s.", password)
	return &passwordAuth{password}
}
