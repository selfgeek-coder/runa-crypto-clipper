package hide

import (
    "os"
    "syscall"
)

func HideFile(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return err
    }

    pathPtr, err := syscall.UTF16PtrFromString(path)
    if err != nil {
        return err
    }

    attrs, err := syscall.GetFileAttributes(pathPtr)
    if err != nil {
        return err
    }

    hiddenAttrs := attrs | syscall.FILE_ATTRIBUTE_HIDDEN

    err = syscall.SetFileAttributes(pathPtr, hiddenAttrs)
    if err != nil {
        return err
    }

    return nil
}