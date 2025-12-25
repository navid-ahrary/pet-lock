package win_11

import (
	"fmt"
	"os"
	"os/exec"
)

// IsRunningAsAdmin checks if the current process has admin privileges
func IsRunningAsAdmin() bool {
	cmd := exec.Command("cmd", "/C", "net session")
	return cmd.Run() == nil
}

// RelaunchAsAdmin relaunches the current executable with UAC
func RelaunchAsAdmin() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	ps := fmt.Sprintf(`Start-Process "%s" -Verb RunAs`, exe)

	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-NonInteractive",
		"-Command", ps,
	)

	return cmd.Run()
}
