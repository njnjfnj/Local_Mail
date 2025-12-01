package main

import (
	"github.com/njnjfnj/Local_Mail/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()

	w := a.NewWindow("Local Mail")

	w.Resize(fyne.NewSize(450, 740))

	appGUI := gui.NewAppGUI(w)

	w.SetContent(appGUI.CreateMainLayout())

	w.ShowAndRun()
}
