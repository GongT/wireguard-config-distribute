package service

import (
	"bufio"
	"os"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func pause() {
	tools.Error("\nPress enter to quit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
