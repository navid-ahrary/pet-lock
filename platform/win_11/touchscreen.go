package win_11

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func ToggleTouchScreen() {
	id, err := getTouchScreenID()
	if err != nil {
		log.Println("Error detecting TouchScreen:", err)
		return
	}
	if id == "" {
		log.Println("No TouchScreen found")
		return
	}

	enabled, err := isTouchScreenEnabled(id)
	if err != nil {
		log.Println("Error checking TouchScreen state:", err)
		return
	}

	if enabled {
		log.Println("TouchScreen is enabled — disabling it")
		disableTouchScreen()
	} else {
		log.Println("TouchScreen is disabled — enabling it")
		enableTouchScreen()
	}
}

func getTouchScreenID() (string, error) {
	cmd := exec.Command("xinput", "list")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(out), "\n")
	re := regexp.MustCompile(`id=([0-9]+)`)

	for _, line := range lines {
		lower := strings.ToLower(line)

		if strings.Contains(lower, "TouchScreen") || strings.Contains(lower, "synaptics") {
			match := re.FindStringSubmatch(line)
			if len(match) > 1 {
				return match[1], nil
			}
		}
	}

	return "", nil
}

func disableTouchScreen() {
	id, err := getTouchScreenID()
	if err != nil {
		log.Println("Error detecting TouchScreen:", err)
		return
	}
	if id == "" {
		log.Println("No TouchScreen found")
		return
	}

	log.Println("TouchScreen ID:", id)

	exec.Command("xinput", "disable", id).Run()
}

func enableTouchScreen() {
	id, err := getTouchScreenID()
	if err != nil {
		log.Println("Error detecting TouchScreen:", err)
		return
	}
	if id == "" {
		log.Println("No TouchScreen found")
		return
	}

	log.Println("TouchScreen ID:", id)

	exec.Command("xinput", "enable", id).Run()
}

func isTouchScreenEnabled(id string) (bool, error) {
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
