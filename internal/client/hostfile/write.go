package hostfile

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (w *Watcher) WriteBlock(hosts map[string]string) {
	tools.Error("rewrite hosts file...")

	generated := COMMENT_START + "\n"
	for ip, hostline := range hosts {
		generated += hostLine(ip+" "+hostline) + "\n"
	}
	generated += COMMENT_END + "\n"

	contents := ""
	skip := false
	found := false
	for _, oline := range strings.Split(w.current, "\n") {
		line := strings.TrimSpace(oline)
		if skip {
			if line == COMMENT_END {
				// tools.Debug("~! %s", line)
				skip = false
				contents += generated
			} else {
				// tools.Debug("~~ %s", line)
			}
		} else if line == COMMENT_START {
			// tools.Debug(">> %s", line)
			skip = true
			found = true
		} else {
			// tools.Debug("== %s", line)
			contents += oline + "\n"
		}
	}
	if !found {
		contents += generated
	}

	contents = strings.TrimSpace(contents) + "\n"
	if contents != w.current {
		tools.Debug("write hosts file content.")
		err := ioutil.WriteFile(w.file, []byte(contents), os.FileMode(0644))
		if err != nil {
			tools.Error("failed write hosts file: %v", err)
		}
	}
}
