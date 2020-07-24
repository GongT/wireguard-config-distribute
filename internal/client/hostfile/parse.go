package hostfile

import (
	"net"
	"strings"
)

func ParseServices(hosts string) map[string]string {
	var set = make(map[string]string)

	for _, line := range strings.Split(hosts, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		s := strings.Fields(line)
		if len(s) < 2 {
			continue
		}

		ip := net.ParseIP(s[0])
		if ip == nil {
			continue
		}

		for _, k := range s[1:] {
			if strings.HasPrefix(k, "localhost") {
				continue
			}
			if strings.HasPrefix(k, "#") {
				break
			}
			set[k] = ip.String()
		}
	}

	return set
}

func ToArray(set map[string]string) []string {
	var keys = make([]string, len(set))
	i := 0
	for k := range set {
		keys[i] = k
		i = i + 1
	}
	return keys
}
