package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"
)

func (a *AppGUI) createMenuScreen() fyne.CanvasObject {
	backButton := widget.NewButton("Back", func() {
		a.navigateBackToList()
	})

	settingsButton := widget.NewButton("Settings", func() {
		a.navigateToSettings()
	})

	aboutButton := widget.NewButton("About", func() {

	})

	title := widget.NewLabel("Menu")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	topAppBar := container.NewBorder(
		nil, nil, // top, bottom
		backButton, // left
		nil,        // right
		title,      // center
	)

	vbox := container.NewVBox(widget.NewLabel(""), settingsButton, aboutButton)

	return container.NewBorder(topAppBar, nil, nil, nil, vbox)
}
