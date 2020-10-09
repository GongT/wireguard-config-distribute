package hostfile

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (w *Watcher) WriteBlock(vpnNetworkGroupName string, hosts map[string]string) {
	tools.Error("rewrite hosts file...")

	fullStart := COMMENT_START + " :: " + vpnNetworkGroupName
	generated := fullStart + "\n"
	for ip, hostline := range hosts {
		generated += hostLine(ip+" "+hostline) + "\n"
	}
	generated += COMMENT_END + " :: " + vpnNetworkGroupName + "\n"

	contents := ""
	skip := false
	found := false
	blockSet := false
	for _, oline := range strings.Split(w.current, "\n") {
		line := strings.TrimSpace(oline)
		if skip {
			if strings.HasPrefix(line, COMMENT_END) {
				tools.Debug("<< %s", line)
				skip = false
				if !blockSet {
					blockSet = true
					contents += generated
				}
			} else {
				// tools.Debug("!! %s", line)
			}
		} else if line == fullStart {
			tools.Debug(">>B %s", line)
			skip = true
			found = true
		} else if line == COMMENT_START {
			tools.Debug(">>A %s", line)
			skip = true
		} else {
			// tools.Debug("== %s", line)
			contents += oline + "\n"
		}
	}
	contents = strings.TrimSpace(contents) + "\n"

	if !found {
		contents += "\n" + generated + "\n"
	}

	if contents != w.current {
		tools.Debug("write hosts file content.")
		err := ioutil.WriteFile(w.file, []byte(contents), os.FileMode(0644))
		if err != nil {
			tools.Error("failed write hosts file: %v", err)
		}
	}
}
