// +build !windows

package service

import (
	"os"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func EnsureAdminPrivileges(_ interface{}) {
	if os.Getuid() != 0 {
		tools.Die("root privilege is required.")
	}
}
