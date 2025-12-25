package win_11

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type pnpDevice struct {
	FriendlyName string `json:"FriendlyName"`
	InstanceId   string `json:"InstanceId"`
}

func ToggleTouchpad() {
	touchpadInstanceID, err := getTouchpadDeviceID()
	if err != nil {
		log.Println("Error in finding touchpad id", err)
	}

	enabled, err := isTouchpadEnabled(touchpadInstanceID)
	if err != nil {
		log.Println("Error checking touchpad state:", err)
		return
	}

	var ps string

	if enabled {
		log.Println("Touchpad is enabled — disabling it")
		ps = fmt.Sprintf(`Disable-PnpDevice -InstanceId "%s" -Confirm:$false`, touchpadInstanceID)
	} else {
		log.Println("Touchpad is disabled — enabling it")
		ps = fmt.Sprintf(`Enable-PnpDevice -InstanceId "%s" -Confirm:$false`, touchpadInstanceID)
	}

	executePowerShell(ps)
}

func getTouchpadDeviceID() (string, error) {
	ps := `
Get-PnpDevice -PresentOnly |
Where-Object { $_.Class -eq 'HIDClass' } |
Select FriendlyName, InstanceId |
ConvertTo-Json
`

	stdout, stderr, err := executePowerShell(ps)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(stderr) != "" {
		return "", errors.New("powershell error: " + stderr)
	}

	raw := strings.TrimSpace(stdout)
	if raw == "" || raw == "null" {
		return "", errors.New("no devices returned from PowerShell")
	}

	var devices []pnpDevice

	// Handle object vs array output
	if strings.HasPrefix(raw, "{") {
		var single pnpDevice
		if err := json.Unmarshal([]byte(raw), &single); err != nil {
			return "", err
		}
		devices = append(devices, single)
	} else {
		if err := json.Unmarshal([]byte(raw), &devices); err != nil {
			return "", err
		}
	}

	type scoredDevice struct {
		id    string
		score int
	}

	var candidates []scoredDevice

	for _, d := range devices {
		id := strings.ToLower(d.InstanceId)
		name := strings.ToLower(d.FriendlyName)

		score := 0

		// Exclusions
		if strings.HasPrefix(id, "usb\\") {
			continue
		}
		if strings.HasPrefix(id, "acpi\\") {
			continue
		}

		// Strong vendor matches
		if strings.Contains(id, "elan") {
			score += 50
		}
		if strings.Contains(id, "syn") {
			score += 50
		}
		if strings.Contains(id, "alps") {
			score += 50
		}
		if strings.Contains(id, "msft") {
			score += 40
		}

		// Name-based hints
		if strings.Contains(name, "touchpad") {
			score += 40
		} else if strings.Contains(name, "touch") {
			score += 25
		}

		// Pen / digitizer penalty
		if strings.Contains(name, "pen") || strings.Contains(name, "digitizer") {
			score -= 40
		}

		if score > 0 {
			candidates = append(candidates, scoredDevice{
				id:    d.InstanceId,
				score: score,
			})
		}
	}

	if len(candidates) == 0 {
		return "", errors.New("touchpad not found")
	}

	best := candidates[0]
	for _, c := range candidates {
		if c.score > best.score {
			best = c
		}
	}

	return best.id, nil
}

func isTouchpadEnabled(touchpadInstanceID string) (bool, error) {
	if touchpadInstanceID == "" {
		return false, errors.New("empty touchpad instance id")
	}

	ps := `
$dev = Get-PnpDevice -InstanceId "` + touchpadInstanceID + `"
if ($null -eq $dev) {
    Write-Error "Device not found"
    exit 1
}
$dev.Status
`
	stdout, stderr, err := executePowerShell(ps)
	if err != nil {
		return false, err
	}

	if stderr != "" {
		return false, errors.New(stderr)
	}

	status := strings.TrimSpace(strings.ToLower(stdout))

	// Windows reports enabled devices as "ok"
	return status == "ok", nil
}

func executePowerShell(script string) (string, string, error) {
	log.Println("Executing PowerShell:")
	log.Println(script)

	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy", "Bypass",
		"-Command", script,
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	log.Println("PowerShell STDOUT:", stdout.String())
	log.Println("PowerShell STDERR:", stderr.String())

	if err != nil {
		log.Println("PowerShell ERROR:", err)
	}

	return stdout.String(), stderr.String(), err
}
