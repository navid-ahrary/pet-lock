package main

import (
	"image/color"
	"log"

	"pet-lock/internal/logging"
	"pet-lock/platform"
	linux_x11 "pet-lock/platform/linux-x11"
	"pet-lock/platform/win_11"
	"pet-lock/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/hotkey"
)

const (
	windowWidth  = 250
	windowHeight = 200

	DeviceKeyboard    = "⌨️ Keyboard"
	DeviceTouchPad    = "🖱️ Touch Pad"
	DeviceTouchScreen = "👆 Touch Screen"
)

var (
	locked bool
	os     = platform.DetectPlatform()
)

func main() {
	logFile, err := logging.Init("PetLock")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.Println("Application starting")

	if os == "win" {
		if !win_11.IsRunningAsAdmin() {
			log.Println("Requesting Administrator privileges...")
			_ = win_11.RelaunchAsAdmin()
			return
		}
	}

	myApp := app.New()
	setIcon(myApp)

	mainWindow := myApp.NewWindow("Pet Lock! 🐾")
	mainWindow.Resize(resizeWindow())
	mainWindow.SetFixedSize(true)

	labelWidget := widget.NewLabel("What would you like to lock?")
	deviceWidgets := &widget.CheckGroup{
		Required:   true,
		Horizontal: true,
		Options:    []string{DeviceKeyboard, DeviceTouchPad, DeviceTouchScreen},
		Selected:   []string{DeviceKeyboard, DeviceTouchPad, DeviceTouchScreen},
	}

	hintLabel := widget.NewLabel("Shift+Ctrl+L to lock/unlock!")
	status := ui.NewStatus("Unlocked")
	status.Text.TextStyle.Bold = true
	bottomWidget := container.NewBorder(nil, nil, hintLabel, status.Text)

	mainWindow.SetContent(container.NewBorder(labelWidget, bottomWidget, nil, deviceWidgets))

	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyL)
		if err := hk.Register(); err != nil {
			log.Fatalf("The hotkey not registered: %v", err)
		}

		for range hk.Keydown() {
			fyne.DoAndWait(func() {
				locked = !locked

				if locked {

					deviceWidgets.Disable()
				} else {
					deviceWidgets.Enable()
				}

				updateStatus(locked, status)
			})

			toggleSelectedDevices(deviceWidgets)

		}
	}()

	mainWindow.ShowAndRun()

	tidyUp()
}

func tidyUp() {
	log.Println("Exited")
}

func setIcon(myApp fyne.App) {
	iconResource, err := fyne.LoadResourceFromPath("./assets/app_icon.png")
	if err != nil {
		log.Println("Error loading icon:", err)
	} else {
		myApp.SetIcon(iconResource)
	}
}

func resizeWindow() fyne.Size {
	return fyne.Size{Width: windowWidth, Height: windowHeight}
}

func toggleSelectedDevices(checkGroup *widget.CheckGroup) {
	selected := checkGroup.Selected

	if contains(selected, DeviceTouchPad) {

		switch os {
		case "linux":
			linux_x11.ToggleTouchpad()
		case "win":
			win_11.ToggleTouchpad()
		}
	}
	if contains(selected, DeviceKeyboard) {
		switch os {
		case "linux":
			linux_x11.ToggleKeyboard()
		case "win":
			win_11.ToggleKeyboard()
		}
	}
	if contains(selected, DeviceTouchScreen) {
		switch os {
		case "linux":
			linux_x11.ToggleTouchScreen()
		case "win":
			win_11.ToggleTouchScreen()
		}
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func updateStatus(locked bool, status *ui.Status) {
	if locked {
		status.Set("Locked", color.RGBA{200, 0, 0, 255}) // red
	} else {
		status.Set("Unlocked", color.RGBA{0, 200, 0, 255}) // green
	}
}
