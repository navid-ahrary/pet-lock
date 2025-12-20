package win_11

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func ToggleKeyboard() {
	id, err := getKeyboardID()
	if err != nil {
		log.Println("Error detecting Keyboard:", err)
		return
	}
	if id == "" {
		log.Println("No Keyboard found")
		return
	}

	enabled, err := isKeyboardEnabled(id)
	if err != nil {
		log.Println("Error checking Keyboard state:", err)
		return
	}

	if enabled {
		log.Println("Keyboard is enabled — disabling it")
		disableKeyboard()
	} else {
		log.Println("Keyboard is disabled — enabling it")
		enableKeyboard()
	}
}

func getKeyboardID() (string, error) {
	cmd := exec.Command("xinput", "list")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	re := regexp.MustCompile(`id=([0-9]+)`)

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.Contains(lower, "Keyboard") || strings.Contains(lower, "synaptics") {
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				return match[1], nil
			}
		}
	}

	return "", nil
}

func disableKeyboard() {
	id, err := getKeyboardID()
	if err != nil {
		log.Println("Error detecting Keyboard:", err)
		return
	}
	if id == "" {
		log.Println("No Keyboard found")
		return
	}

	log.Println("Keyboard ID:", id)

	exec.Command("xinput", "disable", id).Run()
}

func enableKeyboard() {
	id, err := getKeyboardID()
	if err != nil {
		log.Println("Error detecting Keyboard:", err)
		return
	}
	if id == "" {
		log.Println("No Keyboard found")
		return
	}

	log.Println("Keyboard ID:", id)

	exec.Command("xinput", "enable", id).Run()
}

func isKeyboardEnabled(id string) (bool, error) {
	cmd := exec.Command("xinput", "list-props", id)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(out), "\n")

	for _, line := range lines {
		if strings.Contains(line, "Device Enabled") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				return fields[len(fields)-1] == "1", nil
			}
		}
	}

	return false, fmt.Errorf("device Enabled property not found")
}
