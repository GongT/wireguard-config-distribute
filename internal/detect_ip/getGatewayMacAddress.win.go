// +build windows

package detect_ip

import (
	"errors"
	"net"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func findGatewayAddr(gwIp net.IP) (string, error) {
	localIp, err := findRouteFromIp(gwIp)
	if err != nil {
		return "", err
	}
	localIpStr := localIp.String()
	gatewayIpStr := gwIp.String()

	tools.Debug("arp -a -N %s", localIpStr)
	routeCmd := exec.Command("arp", "-a", "-N", localIpStr)
	routeCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := routeCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 0 && fields[0] == gatewayIpStr {
			return fields[1], nil
		}
	}

	return "", errors.New("Default gateway did not appear on ARP table.")
}
