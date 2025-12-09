package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"
)

func (a *AppGUI) createMenuScreen() fyne.CanvasObject {
	backButton := widget.NewButton("Back", func() {
		a.navigateBackToList()
	})

	a.window.SetOnClosed(func() {
		a.navigateBackToList()
	})

	settingsButton := widget.NewButton("Settings", func() {
		a.navigateToSettings()
	})

	aboutButton := widget.NewButton("About", func() {
		a.navigateToAbout()
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

func (a *AppGUI) createAboutScreen() fyne.CanvasObject {
	backButton := widget.NewButton("Back", func() {
		a.navigateToMenu()
	})
	a.window.SetOnClosed(func() {
		a.navigateToMenu()
	})

	title := widget.NewLabel("About")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	topAppBar := container.NewBorder(
		nil, nil, // top, bottom
		backButton, // left
		nil,        // right
		title,      // center
	)
	img := canvas.NewImageFromFile("icons/CreatorPhoto.jpg")
	img.FillMode = canvas.ImageFillContain
	img.SetMinSize(fyne.NewSize(300, 300))
	content := container.NewMax(img)

	text := widget.NewLabel(`
	Local mail was created by Ilhin Serhii, student of National Technical University of Ukraine "Igor Sikorsky Kyiv Polytechnic Institute"

	Institute of Applied System Analysis
	
	Ukraine Cherkassy 
	2025
	`)
	text.Wrapping = fyne.TextWrapWord
	text.Alignment = fyne.TextAlignCenter

	vbox := container.NewVBox(content, text)

	lbl1 := widget.NewLabel("powered by golang & framework fyne")
	lbl1.Alignment = fyne.TextAlignCenter

	return container.NewBorder(topAppBar, lbl1, nil, nil, vbox)
}
