package child_process

import (
	"os"
	"os/exec"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func MustSuccess(title, cmd string, args ...string) {
	err := runCmd(cmd, args...)
	if err != nil {
		tools.Die("failed %s: [%s %s]: %s", title, cmd, strings.Join(args, " "), err.Error())
	}
}

func ShouldSuccess(title, cmd string, args ...string) error {
	err := runCmd(cmd, args...)
	if err != nil {
		tools.Error("failed %s: [%s %s]: %s", title, cmd, strings.Join(args, " "), err.Error())
	}
	return err
}

func runCmd(cmd string, args ...string) error {
	tools.Debug("\x1B[2m%s %s\x1B[0m", cmd, strings.Join(args, " "))
	p := exec.Command(cmd, args...)
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	p.Env = append(os.Environ(), "LANG=C")
	return p.Run()
}

func RunGetOutput(title, cmd string, args ...string) string {
	p := exec.Command(cmd, args...)
	p.Env = append(os.Environ(), "LANG=C")
	ret, err := p.CombinedOutput()

	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			tools.Die("failed %s: [%s %s]: %s", title, cmd, strings.Join(args, " "), err.Error())
		}
	}

	return string(ret)
}

func RunGetStandardOutput(title, cmd string, args ...string) string {
	p := exec.Command(cmd, args...)
	p.Env = append(os.Environ(), "LANG=C")
	p.Stderr = os.Stderr
	ret, err := p.Output()

	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			tools.Die("failed %s: [%s %s]: %s", title, cmd, strings.Join(args, " "), err.Error())
		}
	}

	return string(ret)
}

func RunGetReturnCode(title, cmd string, args ...string) int {
	p := exec.Command(cmd, args...)
	p.Env = append(os.Environ(), "LANG=C")

	if tools.IsDevelopmennt() {
		tools.Error("%s %s", cmd, strings.Join(args, " "))
		p.Stderr = os.Stderr
		p.Stdout = os.Stdout
	}

	if err := p.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			tools.Error("failed %s: [%s %s]: %s", title, cmd, strings.Join(args, " "), err.Error())
			return 1
		}
	}

	tools.Debug("process exit with code %v", p.ProcessState.ExitCode())
	return p.ProcessState.ExitCode()
}
