// +build windows

package hostfile

import "strings"

func hostLine(hostline string) string {
	p := strings.SplitN(hostline, "#", 2)
	if len(p) == 1 {
		return hostline
	} else {
		return "#" + p[1] + "\n" + strings.TrimSpace(p[0])
	}
}
