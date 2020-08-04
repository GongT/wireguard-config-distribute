package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (opts *serverProgramOptions) Sanitize() error {
	if !tools.GetSystemHostName(&opts.ServerName) {
		return fmt.Errorf("HOSTNAME and COMPUTERNAME is empty, please set --server-name")
	}

	if len(opts.GetStorageLocation()) == 0 {
		if path, exists := os.LookupEnv("STATE_DIRECTORY"); exists {
			fmt.Println("use storage path from STATE_DIRECTORY")
			opts.StorageLocation = path
		} else {
			fmt.Println("use storage path from user home dir")
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("Failed get user HOME: %s", err.Error())
			}
			opts.StorageLocation = filepath.Join(home, ".wireguard-config-server")
		}
	}
	if len(opts.GetStorageLocation()) == 0 {
		return fmt.Errorf("Need a storage path, please set --storage")
	}
	return nil
}
