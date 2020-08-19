package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

const distPath = `C:\Program Files\WireGuard\config-client.exe`

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

func stopService(s *mgr.Service) error {
	for {
		if status, err := s.Query(); err != nil {
			return fmt.Errorf("failed query service: %v", err)
		} else if status.State == svc.Stopped {
			break
		} else if status.State == svc.StopPending {
			time.Sleep(1)
		} else {
			if _, err := s.Control(svc.Stop); err != nil {
				return fmt.Errorf("failed stop service: %v", err)
			}
		}
	}

	return nil
}

func startService(s *mgr.Service) error {
	for {
		if status, err := s.Query(); err != nil {
			return fmt.Errorf("failed query service: %v", err)
		} else if status.State == svc.Running {
			break
		} else if status.State == svc.StartPending {
			time.Sleep(1)
		} else {
			if _, err := s.Control(svc.Continue); err != nil {
				return fmt.Errorf("failed start service: %v", err)
			}
		}
	}

	return nil
}

func installService(opts elevateOptions, install bool) error {
	serviceName := opts.GetInterfaceName()

	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed connect windows service manager: %s", err.Error())
	}
	defer m.Disconnect()

	service, serr := m.OpenService(serviceName)
	exists := serr == nil
	if exists {
		defer service.Close()
	}

	if install {
		if exists {
			return fmt.Errorf("service already exists")
		}

		exec, err := exePath()
		if err != nil {
			return fmt.Errorf("invalid executable file: %s", err.Error())
		}

		if err := tools.CopyFile(exec, distPath); err != nil {
			return fmt.Errorf("failed copy binary: %v", err)
		}

		if err := _install(m, distPath, serviceName); err != nil {
			return fmt.Errorf("failed install windows service: %v", err)
		}

		if err := startService(service); err != nil {
			tools.Error("failed start windows service: %v", err)
		}
	} else {
		if !exists {
			tools.Error("service did not exists")
			return nil
		}

		if err := stopService(service); err != nil {
			return fmt.Errorf("failed stop windows service: %v", err)
		}
		if err := _uninstall(m, serviceName); err != nil {
			return fmt.Errorf("failed uninstall windows service: %v", err)
		}
	}

	return nil
}

func _install(m *mgr.Mgr, execPath string, serviceName string) error {
	config.ResetInternalOption()
	s, err := m.CreateService(serviceName, execPath, mgr.Config{
		DisplayName: "Wireguard Config (" + serviceName + ")",
		Description: "Wireguard auto config service\r\nWireGuard VPN自动配置服务",
		StartType:   mgr.StartAutomatic,
		// ServiceStartName: "NetworkService",
		DelayedAutoStart: true,
		Dependencies:     []string{"WireGuardManager"},
	}, config.StringifyOptions(false))
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

	if _, err := s.Control(svc.Stop); err != nil {
		tools.Error("Stop service: %v", err)
	}

	if err := s.Delete(); err != nil {
		return err
	}

	eventlog.Remove(serviceName)
	return nil
}
