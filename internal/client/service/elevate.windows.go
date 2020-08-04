// +build windows

package service

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/config"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/natefinch/npipe"
	"golang.org/x/sys/windows"
)

type elevateOptions interface {
	GetInterfaceName() string
	GetInstallService() bool
	GetUnInstallService() bool
}

// https://github.com/golang/go/issues/28804
func EnsureAdminPrivileges(opts elevateOptions) {
	var sid *windows.SID

	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		log.Fatalf("SID Error: %s", err)
		return
	}
	defer windows.FreeSid(sid)

	// This appears to cast a null pointer so I'm not sure why this
	// works, but this guy says it does and it Works for Me™:
	// https://github.com/golang/go/issues/28804#issuecomment-438838144
	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		tools.Die("Token Membership Error: %s", err.Error())
	}

	if member /*&& token.IsElevated() */ {
		if install, uninstall := opts.GetInstallService(), opts.GetUnInstallService(); install || uninstall {
			if install == uninstall {
				tools.Die("Can not use /install and /uninstall at same time")
			}
			var err error
			if install {
				tools.Error("Install Windows Service...")
				err = installService(opts, true)
			} else if uninstall {
				tools.Error("Uninstall Windows Service...")
				err = installService(opts, false)
			}
			if err == nil {
				tools.Error("Install success!")
			} else {
				tools.Error("Failed install service! %s", err.Error())
			}
			os.Exit(0)
		}
		return
	}

	tools.Error("member=%v ; IsElevated=%v", member, token.IsElevated())
	if config.IsSuRunning() {
		tools.Error("Failed start with admin permission")
		os.Exit(1)
	} else {
		tools.Error("Restart self with admin permission...")
		runMeElevated()
	}
}

// https://stackoverflow.com/questions/31558066/how-to-ask-for-administer-privileges-on-windows-with-go
func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()

	config.InternalOption.IsElevated = true
	config.InternalOption.StandardOutputPath = `\\.\pipe\wireguard-config-client-elevate`
	config.CommonOption.LogFilePath = ""

	wg := sync.WaitGroup{}
	wg.Add(3)
	go handleStream(createPipe("stdout"), os.Stdout, &wg)
	go handleStream(createPipe("stderr"), os.Stderr, &wg)
	go handleControl(createPipe("control"), &wg)

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(config.StringifyOptions())

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		tools.Die("windows self execute failed: %s", err.Error())
	}

	wg.Wait()
	os.Exit(0)
}

func createPipe(chType string) *npipe.PipeListener {
	ln, err := npipe.Listen(config.InternalOption.StandardOutputPath + "." + chType)
	if err != nil {
		waitExit("failed create pipes", err)
	}

	return ln
}

func waitExit(reason string, err error) {
	tools.Error(reason+": %v", err)
	time.Sleep(10 * time.Second)
	tools.Die(reason)
}

func handleControl(lis *npipe.PipeListener, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := lis.Accept()
	lis.Close()
	if err != nil {
		waitExit("failed accpet pipes", err)
	}
	defer conn.Close()

	tools.Debug("Control socket connected.")

	go func() {
		<-tools.WaitForCtrlC()
		conn.Write([]byte("sigint\n"))
	}()

	fscanner := bufio.NewScanner(conn)
	for fscanner.Scan() {
		text := fscanner.Text()
		if strings.HasPrefix(text, "exit:") {
			code, _ := strconv.ParseInt(text[5:], 10, 32)
			tools.HasError(int(code))
			tools.Exit()
		}
	}

	tools.Debug("Control socket finished.")
}

func handleStream(lis *npipe.PipeListener, out *os.File, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := lis.Accept()
	lis.Close()
	if err != nil {
		waitExit("failed accpet pipes", err)
	}
	defer conn.Close()

	tools.Debug("Output socket (%s) connected.", out.Name())

	_, err = io.Copy(out, conn)
	tools.Debug("Output socket (%s) finished.", out.Name())

	if err != nil {
		waitExit("failed read pipes", err)
	}
}
