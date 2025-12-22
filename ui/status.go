package ui

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
)

type Status struct {
	Text *canvas.Text
}

func NewStatus(initial string) *Status {
	t := canvas.NewText(initial, color.RGBA{0, 200, 0, 255}) // default green
	t.TextSize = 14

	return &Status{Text: t}
}

func (s *Status) Set(text string, col color.Color) {
	s.Text.Text = text
	s.Text.Color = col
	s.Text.Refresh()
}
