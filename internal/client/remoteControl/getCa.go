package remoteControl

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (tool *ToolObject) GetCA(target string) {
	ret, err := tool.server.GetSelfSignedCertFile(&protocol.GetCertFileRequest{})
	if err != nil {
		tools.Die("Failed request server: %s", err.Error())
	}

	if !filepath.IsAbs(target) {
		pwd, err := os.Getwd()
		if err != nil {
			tools.Die("Failed get working directory: %s", err.Error())
		}

		target = filepath.Join(pwd, target)
	}

	tools.Error("Write file: %s", target)

	if err := os.MkdirAll(filepath.Dir(target), os.FileMode(0755)); err != nil {
		tools.Die("Failed make parent directory of file: %s", err.Error())
	}
	err = ioutil.WriteFile(target, ret.CertFileText, os.FileMode(0644))
	if err != nil {
		tools.Die("Failed write file: %s", err.Error())
	}

	tools.Error("Complete!")
}
