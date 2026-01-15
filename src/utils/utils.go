package utils

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
)

// gets active window title
func GetActiveWindow() string {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		return ""
	}

	buf := make([]uint16, 256)
	procGetWindowTextW.Call(
		hwnd,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)

	return syscall.UTF16ToString(buf)
}

// retrieves the absolute path to the current executable
func GetSelfDir() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	return dir, nil
}

// retrieves the name to the current executable
func GetSelfName() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	fileName := filepath.Base(execPath)

	name := strings.TrimSuffix(fileName, ".exe")

	return name, nil
}
