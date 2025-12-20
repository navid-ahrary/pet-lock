package main

import (
	"log"

	linux_x11 "pet-lock/platform/linux-x11"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	hook "github.com/robotn/gohook"
)

const (
	windowWidth  = 250
	windowHeight = 200

	KeyL        = 108
	KeyLeftAlt  = 65513
	KeyRightAlt = 65514

	DeviceKeyboard    = "⌨️ Keyboard"
	DeviceTouchPad    = "🖱️ Touch Pad"
	DeviceTouchScreen = "🖥️ Touch Screen"
)

var (
	deviceWidgets *widget.CheckGroup
	locked        bool
)

func main() {
	myApp := app.New()
	setIcon(myApp)

	mainWindow := myApp.NewWindow("Pet Lock! 🐾")
	mainWindow.Resize(resizeWindow())
	mainWindow.SetFixedSize(true)

	evChan := hook.Start()
	defer hook.End()

	var altDown bool
	go func() {
		for ev := range evChan {
			switch ev.Kind {
			case hook.KeyDown:
				if isAlt(ev.Rawcode) {
					altDown = true
				}

				if isL(ev.Rawcode) && altDown {
					log.Println("Alt+L pressed")

					fyne.DoAndWait(func() {
						locked = !locked
						if locked {
							deviceWidgets.Disable()
						} else {
							deviceWidgets.Enable()
						}
					})

					toggleSelectedDevices()
				}

			case hook.KeyUp:
				if isAlt(ev.Rawcode) {
					altDown = false
				}
			}
		}
	}()

	labelWidget := widget.NewLabel("What would you like to lock?")
	deviceWidgets = generateDeviceWidgets()
	bottomContainer := widget.NewLabel("Active/Inactive")

	mainWindow.SetContent(container.NewBorder(labelWidget, bottomContainer, nil, nil, deviceWidgets))
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

func generateDeviceWidgets() *widget.CheckGroup {
	return &widget.CheckGroup{
		Required:   true,
		Horizontal: true,
		Options:    []string{DeviceKeyboard, DeviceTouchPad, DeviceTouchScreen},
		Selected:   []string{DeviceKeyboard, DeviceTouchPad, DeviceTouchScreen},
	}
}

func isAlt(raw uint16) bool {
	return raw == KeyLeftAlt || raw == KeyRightAlt
}

func isL(raw uint16) bool {
	return raw == KeyL
}

func toggleSelectedDevices() {
	selected := deviceWidgets.Selected
	if contains(selected, DeviceTouchPad) {
		log.Println("ToggleTouchPad")
		linux_x11.ToggleTouchPad()
	}
	if contains(selected, DeviceKeyboard) {
		linux_x11.ToggleKeyboard()
	}
	if contains(selected, DeviceTouchScreen) {
		linux_x11.ToggleTouchScreen()
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
