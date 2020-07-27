package vpnManager

import (
	"strconv"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type NumberBasedIp uint32

func FromNumber(ip string) NumberBasedIp {
	ips := strings.Split(ip, ".")
	l := 0
	var v uint32

	for i := len(ips) - 1; i >= 0; i -= 1 {
		vv, err := strconv.ParseUint(ips[i], 10, 8)

		if err != nil {
			tools.Error("Failed parse IP [%s]: %s", ip, err.Error())
			return 0
		}

		v += uint32(vv << l)

		l += 8
	}

	return NumberBasedIp(v)
}

func (v NumberBasedIp) String(part uint) (ret string) {
	for i := uint(0); i < part; i += 1 {
		ret = "." + strconv.FormatUint(uint64(uint8(v>>(8*i))), 10) + ret
	}
	return ret[1:]
}
