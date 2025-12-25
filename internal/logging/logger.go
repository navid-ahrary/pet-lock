package logging

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

// Init initializes logging to a file (and stdout optionally)
func Init(appName string) (*os.File, error) {
	logDir := filepath.Join(os.TempDir(), appName)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logFile := filepath.Join(
		logDir,
		time.Now().Format("20060102_150405")+".log",
	)

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Logger initialized")
	log.Println("Log file:", logFile)

	return f, nil
}
