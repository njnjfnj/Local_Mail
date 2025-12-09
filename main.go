package main

import (
	"github.com/njnjfnj/Local_Mail/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.NewWithID("com.localmail.app")
	a := app.New()

	w := a.NewWindow("Local Mail")

	w.Resize(fyne.NewSize(450, 740))

	appGUI := gui.NewAppGUI(w, myApp)

	w.SetContent(appGUI.CreateMainLayout())

	w.ShowAndRun()
}
