package antivirus

import (
	"os/exec"
	"strings"
	"syscall"
)

// gets a list of installed antivirus products on the system
// e.g. "Windows Defender" ...
func GetInstalledAntiviruses() []string {
	cmd := exec.Command(
		"powershell", "-Command", 
		`Get-CimInstance -Namespace root/SecurityCenter2 -ClassName AntiVirusProduct | Select-Object -ExpandProperty displayName`,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}
	
	result := []string{}
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}