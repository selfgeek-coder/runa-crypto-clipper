package autorun

import (
	"fmt"
	"os/exec"
)

// adds the given path to the Task Scheduler to run at logon
func AddToSchelduler(path, name string) error {
	cmd := exec.Command("schtasks", 
		"/Create", 
		"/F", 
		"/SC", "ONLOGON", 
		"/TN", name, 
		"/TR", fmt.Sprintf(`"%s"`, path), 
		"/RL", "HIGHEST",
	)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}