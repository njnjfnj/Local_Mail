package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	messagetype "github.com/njnjfnj/Local_Mail/gui/message_type"
)

func (a *AppGUI) createChatViewScreen(contactName string) (*widget.List, fyne.CanvasObject) {
	// 1. Верхняя панель (App Bar)
	backButton := widget.NewButton("Back", func() {
		a.navigateBackToList() // Навигация назад
	})
	infoButton := widget.NewButton("Info", func() {
		// Логика 'Инфо о контакте'
	})
	title := widget.NewLabel(contactName)
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	topAppBar := container.NewBorder(
		nil, nil, // top, bottom
		backButton, // left
		infoButton, // right
		title,      // center
	)

	// 2. "Подвал" (Footer) с полем ввода
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Message...")

	sendButton := widget.NewButton("Send", func() {
		// Логика отправки сообщения
		inputEntry.SetText("")
	})

	attachFileButton := widget.NewButton("pin file", func() {
		// Логика прикрепления файла
	})

	attachPhotoButton := widget.NewButton("pin photo", func() {

	})

	attachContainer := container.NewVBox(attachFileButton, attachPhotoButton)

	bottomInputBar := container.NewBorder(
		nil, nil, // top, bottom
		attachContainer, // left
		sendButton,      // right
		inputEntry,      // center
	)

	// 3. Список сообщений
	// 3. Список сообщений
	messages := []string{ // взять с бд
		"Hi!",
		"Hey, how are you?",
		"I'm good, thanks! Fyne is cool.",
		"Totally agree with you, Fyne is great for cross-platform UI. This is a longer message to test wrapping.",
		"That's a very long message indeed, let's see how it looks. It should wrap nicely.",
		"Short one.",
	}

	messageList := widget.NewList(
		func() int {
			return len(messages)
		},
		func() fyne.CanvasObject {
			message := messagetype.New_message("user1", "new_text string", "", "")

			return container.NewVBox(message.Text, message.Image, message.File)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			c.Objects[0].(*widget.Label).SetText(messages[i])
		},
	)

	// 4. Собираем экран
	layout := container.NewBorder(topAppBar, bottomInputBar, nil, nil, messageList)

	// ИЗМЕНЕНИЕ: Возвращаем и список, и layout
	return messageList, layout
}
