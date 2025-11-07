package main

import (
	"github.com/njnjfnj/Local_Mail/gui" // Импортируем наш пакет gui

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	// Создаем новое Fyne-приложение
	a := app.New()

	// Создаем главное окно
	w := a.NewWindow("Local Mail")

	// Устанавливаем размер, типичный для мобильного телефона (для удобства отладки на ПК)
	w.Resize(fyne.NewSize(360, 740))

	// Создаем наш кастомный GUI и передаем ему окно
	appGUI := gui.NewAppGUI(w)

	// Устанавливаем главный экран (список чатов) как контент окна
	w.SetContent(appGUI.CreateMainLayout())

	// Запускаем приложение
	w.ShowAndRun()
}
