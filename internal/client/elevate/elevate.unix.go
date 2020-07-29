// +build !windows

package elevate

import (
	"os"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func EnsureAdminPrivileges() {
	if os.Getuid() != 0 {
		tools.Die("root privilege is required.")
	}
}
