package gui

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	messagetype "github.com/njnjfnj/Local_Mail/gui/message_type"
	tls_communication "github.com/njnjfnj/Local_Mail/internal/local_net/tls_communication"
	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
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
		a.updateChatViewChan <- *messagetype.New_text_message(fullAddr, "ME: "+a.inputEntry.Text)
		tls_communication.SendPackage(fullAddr, mail_data{
			Package_type: 1,
			Username:     contactName,
			FullAddress:  local_net.GetOutboundIP() + ":" + a.settingsScreenWidgets.Port.Text,
			Message:      a.inputEntry.Text,
		})
		a.inputEntry.SetText("")
	})

	attachFileButton := widget.NewButton("send file", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			filePath := reader.URI().Path()
			sourceFile, _ := os.Open(filePath)
			defer sourceFile.Close()

			os.MkdirAll("Shared", 0755)
			destPath := filepath.Join("Shared", filepath.Base(filePath))

			destFile, err := os.Create(destPath)
			if err != nil {
				log.Println("error os.Create: ", err.Error())
			}
			defer destFile.Close()

			if _, err = io.Copy(destFile, sourceFile); err != nil {
				log.Println("error io.Copy: ", err.Error())
			}
			a.updateChatViewChan <- *messagetype.New_message(fullAddr, "ME: ", destPath, "", a.window, a.startFileDownloadingChan)

		}, a.window)
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

	a.messageList = widget.NewList(
		func() int {
			a.chatViewMu.RLock()
			defer a.chatViewMu.RUnlock()
			return len(a.temporaryMessagesStorage[fullAddr])
		},
		func() fyne.CanvasObject {
			message := messagetype.New_message("user1", "new_text string", "", "", a.window, a.startFileDownloadingChan)

			return container.NewVBox(message.Text, message.Image, message.File)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			a.chatViewMu.RLock()
			defer a.chatViewMu.RUnlock()
			c.Objects[0].(*widget.Label).SetText(a.temporaryMessagesStorage[fullAddr][i].Text.Text)
			if a.temporaryMessagesStorage[fullAddr][i].File.Visible() {
				c.Objects[2] = messagetype.CopyFileType(a.temporaryMessagesStorage[fullAddr][i].File)
			}
			fmt.Println("Chat: ", fullAddr)
		},
	)

	layout := container.NewBorder(topAppBar, bottomInputBar, nil, nil, a.messageList)

	return a.messageList, layout
}
