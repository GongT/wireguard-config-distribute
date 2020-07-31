package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/config"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func isServiceExists(m *mgr.Mgr, serviceName string) (exists bool, err error) {
	s, err := m.OpenService(serviceName)
	if err == nil {
		s.Close()
		err = fmt.Errorf("service %s already exists", serviceName)
		exists = true
	} else {
		err = fmt.Errorf("service %s did not exists", serviceName)
		exists = false
	}

	return
}

func installService(opts elevateOptions, install bool) error {
	serviceName := opts.GetInterfaceName()

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed connect windows service manager: %s", err.Error())
	}
	defer m.Disconnect()

	if ex, err := isServiceExists(m, serviceName); ex == install {
		return fmt.Errorf("failed install/uninstall service: %s", err.Error())
	}

	if install {
		exec, err := exePath()
		if err != nil {
			return fmt.Errorf("invalid executable file: %s", err.Error())
		}

		err = _install(m, exec, serviceName)
	} else {
		err = _uninstall(m, serviceName)
	}
	if err != nil {
		return fmt.Errorf("failed install/uninstall windows service: %s", err.Error())
	}

	return nil
}

func _install(m *mgr.Mgr, execPath string, serviceName string) error {
	s, err := m.CreateService(serviceName, execPath, mgr.Config{
		DisplayName: "Wireguard Config (" + serviceName + ")",
		Description: "Wireguard auto config service\r\nWireGuard VPN自动配置服务",
		StartType:   mgr.StartAutomatic,
		// ServiceStartName: "NetworkService",
		DelayedAutoStart: true,
		Dependencies:     []string{"WireGuardManager"},
	}, config.StringifyOptions())
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		if !strings.Contains(err.Error(), "exists") {
			s.Delete()
			return fmt.Errorf("InstallAsEventCreate() failed: %s", err)
		}
	}
	return nil
}

func _uninstall(m *mgr.Mgr, serviceName string) error {
	s, err := m.OpenService(serviceName)
	if err != nil {
		return err
	}
	defer s.Close()

	if err := s.Delete(); err != nil {
		return err
	}

	eventlog.Remove(serviceName)
	return nil
}
