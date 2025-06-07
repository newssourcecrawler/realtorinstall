package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func ShowError(win fyne.Window, msg string) {
	popup := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabelWithStyle("Error", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel(msg),
			widget.NewButton("OK", func() { popup.Hide() }),
		), win.Canvas(),
	)
	popup.Show()
}
