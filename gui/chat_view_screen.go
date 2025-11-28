package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	messagetype "github.com/njnjfnj/Local_Mail/gui/message_type"
	tls_communication "github.com/njnjfnj/Local_Mail/internal/local_net/tls_communication"
)

type mail_data struct {
	Package_type int
	Username     string
	FullAddress  string
	Message      string
}

func (a *AppGUI) createChatViewScreen(contactName, fullAddr string) (*widget.List, fyne.CanvasObject) {
	backButton := widget.NewButton("Back", func() {
		a.navigateBackToList() // Навигация назад
	})
	// infoButton := widget.NewButton("Info", func() {
	// })
	title := widget.NewLabel(contactName)
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	topAppBar := container.NewBorder(
		nil, nil, // top, bottom
		backButton, // left
		nil,        //infoButton, // right
		title,      // center
	)

	a.inputEntry.SetPlaceHolder("Message...")

	sendButton := widget.NewButton("Send", func() {
		tls_communication.SendPackage(fullAddr, mail_data{
			Package_type: 1,
			Username:     contactName,
			FullAddress:  fullAddr,
			Message:      a.inputEntry.Text,
		})
		a.inputEntry.SetText("")
	})

	attachFileButton := widget.NewButton("pin file", func() {
	})

	attachPhotoButton := widget.NewButton("pin photo", func() {
	})

	attachContainer := container.NewVBox(attachFileButton, attachPhotoButton)

	bottomInputBar := container.NewBorder(
		nil, nil, // top, bottom
		attachContainer, // left
		sendButton,      // right
		a.inputEntry,    // center
	)

	// messages := []string{ // взять с бд
	// 	"Hi!",
	// 	"Hey, how are you?",
	// 	"I'm good, thanks! Fyne is cool.",
	// 	"Totally agree with you, Fyne is great for cross-platform UI. This is a longer message to test wrapping.",
	// 	"That's a very long message indeed, let's see how it looks. It should wrap nicely.",
	// 	"Short one.",
	// }

	a.messageList = nil
	a.messageList = widget.NewList(
		func() int {
			a.chatViewMu.RLock()
			defer a.chatViewMu.RUnlock()
			return len(a.temporaryMessagesStorage[fullAddr])
		},
		func() fyne.CanvasObject {
			message := messagetype.New_message("user1", "new_text string", "", "")

			return container.NewVBox(message.Text, message.Image, message.File)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			a.chatViewMu.RLock()
			defer a.chatViewMu.RUnlock()
			c.Objects[0].(*widget.Label).SetText(a.temporaryMessagesStorage[fullAddr][i].Text.Text)
		},
	)

	layout := container.NewBorder(topAppBar, bottomInputBar, nil, nil, a.messageList)

	return a.messageList, layout
}
