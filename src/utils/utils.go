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

	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGetUserGeoID = kernel32.NewProc("GetUserGeoID")
	procGetGeoInfoW  = kernel32.NewProc("GetGeoInfoW")
)

// retrieves the active window title (e.g. Google Chrome)
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

	for i, v := range buf {
		if v == 0 {
			return syscall.UTF16ToString(buf[:i])
		}
	}
	return syscall.UTF16ToString(buf)
}

// retrieves full path of the current executable (e.g. D:/test/geek/clipper.exe)
func GetSelfPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(execPath)
}

// retrieves the name of the current executable (without .exe extension)
func GetSelfName() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	fileName := filepath.Base(execPath)
	return strings.TrimSuffix(fileName, ".exe"), nil
}

// retrieves users geographic location code (e.g., RU, UA)
func GetGeo() string {
	geoID, _, _ := procGetUserGeoID.Call(16) // GEOCLASS_NATION = 16

	buffer := make([]uint16, 256)
	ret, _, _ := procGetGeoInfoW.Call(
		uintptr(geoID),
		4, // GEO_ISO2 = 4
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		0,
	)

	if ret == 0 {
		return "unknown"
	}

	for i, v := range buffer {
		if v == 0 {
			return syscall.UTF16ToString(buffer[:i])
		}
	}
	return syscall.UTF16ToString(buffer)
}

func IsElevated() bool {
	ret, _, _ := syscall.NewLazyDLL("shell32.dll").NewProc("IsUserAnAdmin").Call()
	return ret != 0
}