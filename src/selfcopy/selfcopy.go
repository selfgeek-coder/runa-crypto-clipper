package selfcopy

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunFromTemp(selfPath string) bool {
	tempDir := os.TempDir()

	if strings.HasPrefix(strings.ToLower(selfPath), strings.ToLower(tempDir)) {
		return false
	}

	randomBytes := make([]byte, 6)
	rand.Read(randomBytes)
	randomName := hex.EncodeToString(randomBytes) + ".exe"
	tempExe := filepath.Join(tempDir, randomName)

	src, err := os.Open(selfPath)
	if err != nil {
		return false
	}
	defer src.Close()

	dst, err := os.Create(tempExe)
	if err != nil {
		return false
	}
	defer dst.Close()

	io.Copy(dst, src)
	dst.Close()

	cmd := exec.Command(tempExe)
	cmd.Start()

	parentDir := filepath.Dir(selfPath)
	if files, _ := os.ReadDir(parentDir); len(files) == 0 {
		os.Remove(parentDir)
	}

	return true
}