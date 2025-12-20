package win_11

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func ToggleTouchPad() {
	id, err := getTouchpadID()
	if err != nil {
		log.Println("Error detecting touchpad:", err)
		return
	}
	if id == "" {
		log.Println("No touchpad found")
		return
	}

	enabled, err := isTouchpadEnabled(id)
	if err != nil {
		log.Println("Error checking touchpad state:", err)
		return
	}

	if enabled {
		log.Println("Touchpad is enabled — disabling it")
		disableTouchPad()
	} else {
		log.Println("Touchpad is disabled — enabling it")
		enableTouchPad()
	}
}

func getTouchpadID() (string, error) {
	cmd := exec.Command("xinput", "list")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	re := regexp.MustCompile(`id=([0-9]+)`)

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.Contains(lower, "touchpad") || strings.Contains(lower, "synaptics") {
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				return match[1], nil
			}
		}
	}

	return "", nil
}

func disableTouchPad() {
	id, err := getTouchpadID()
	if err != nil {
		log.Println("Error detecting touchpad:", err)
		return
	}
	if id == "" {
		log.Println("No touchpad found")
		return
	}

	log.Println("Touchpad ID:", id)

	exec.Command("xinput", "disable", id).Run()
}

func enableTouchPad() {
	id, err := getTouchpadID()
	if err != nil {
		log.Println("Error detecting touchpad:", err)
		return
	}
	if id == "" {
		log.Println("No touchpad found")
		return
	}

	log.Println("Touchpad ID:", id)

	exec.Command("xinput", "enable", id).Run()
}

func isTouchpadEnabled(id string) (bool, error) {
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
