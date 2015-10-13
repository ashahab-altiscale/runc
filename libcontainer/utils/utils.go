package utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"syscall"
)

const (
	exitSignalOffset = 128
)

// GenerateRandomName returns a new name joined with a prefix.  This size
// specified is used to truncate the randomly generated value
func GenerateRandomName(prefix string, size int) (string, error) {
	id := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, id); err != nil {
		return "", err
	}
	if size > 64 {
		size = 64
	}
	return prefix + hex.EncodeToString(id)[:size], nil
}

// ResolveRootfs ensures that the current working directory is
// not a symlink and returns the absolute path to the rootfs
func ResolveRootfs(uncleanRootfs string) (string, error) {
	rootfs, err := filepath.Abs(uncleanRootfs)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(rootfs)
}

// ExitStatus returns the correct exit status for a process based on if it
// was signaled or existed cleanly.
func ExitStatus(status syscall.WaitStatus) int {
	if status.Signaled() {
		return exitSignalOffset + int(status.Signal())
	}
	return status.ExitStatus()
}

//Checks if host itself usernamespaced, to allow for
//containers in containers case
func IsHostUserns() (bool, error) {
	//scan uid map. should never be more than 5 lines long
	dat, err := ioutil.ReadFile("/proc/self/uid_map")
	reg, err := regexp.Compile("0\\s+0\\s+\\d+")
	if err != nil {
		return false, err
	}
	if reg.Match(dat) {
		return false, nil
	}

	return true, nil
}
