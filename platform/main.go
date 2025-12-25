package platform

import "runtime"

func DetectPlatform() string {
	switch runtime.GOOS {
	case "windows":
		return "win"
	case "linux":
		return "linux"
	default:
		return "unknown"
	}
}
