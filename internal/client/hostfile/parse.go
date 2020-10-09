package hostfile

import (
	"net"
	"strings"
)

func ParseServices(hosts string) map[string]string {
	var set = make(map[string]string)
	skip := true

	for _, line := range strings.Split(hosts, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			if strings.HasPrefix(line, COMMENT_START) {
				skip = true
			} else if strings.HasPrefix(line, COMMENT_END) {
				skip = false
			}
			continue
		}
		if skip {
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
