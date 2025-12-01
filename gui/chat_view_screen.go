package gui

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

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
	FilePath     string
}

const (
	PackageTypeHandshake    = 0
	PackageTypeMessage      = 1
	PackageTypeSendFileInfo = 2
	PackageTypeFileReq      = 3
)

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
		msg := messagetype.My_new_text_message(fullAddr, time.Now().Format("02/01/2006 15:04")+": "+a.inputEntry.Text)
		a.updateChatViewChan <- *msg
		tls_communication.SendPackage(fullAddr, mail_data{
			Package_type: 1,
			Username:     contactName,
			FullAddress:  local_net.GetOutboundIP() + ":" + a.settingsScreenWidgets.Port.Text,
			Message:      time.Now().Format("02/01/2006 15:04") + ": " + a.inputEntry.Text,
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
			a.updateChatViewChan <- *messagetype.New_message(fullAddr, "", destPath, "", a.window, a.startFileDownloadingChan)
			tls_communication.SendPackage(fullAddr, mail_data{
				Package_type: PackageTypeSendFileInfo,
				Username:     a.settingsScreenWidgets.Username.Text,
				FullAddress:  local_net.GetOutboundIP() + ":" + a.settingsScreenWidgets.Port.Text,
				FilePath:     destPath,
			})

		}, a.window)
	})

	// attachPhotoButton := widget.NewButton("pin photo", func() {
	// })

	a.inputEntry.ActionItem = sendButton
	sendButton.Importance = widget.LowImportance

	//attachContainer := container.NewVBox(attachFileButton /*, attachPhotoButton*/)
	bottomInputBar := container.NewBorder(
		nil, nil, // top, bottom
		attachFileButton, // left
		nil,              // right
		a.inputEntry,     // center
	)

	a.messageList = widget.NewList(
		func() int {
			a.chatViewMu.RLock()
			defer a.chatViewMu.RUnlock()
			return len(a.temporaryMessagesStorage[fullAddr])
		},
		func() fyne.CanvasObject {
			message := messagetype.New_message("user1", "~~new_text~string~~1125", "", "", a.window, a.startFileDownloadingChan)

			message.Text.Wrapping = fyne.TextWrapWord
			message.Text.Scroll = fyne.ScrollNone
			// message.Text.TextStyle = fyne.TextStyle{Bold: true}

			message.File.Alignment = widget.ButtonAlignLeading
			message.File.Hide()

			return container.NewVBox(message.Text, message.File /*, message.Image*/)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			c := o.(*fyne.Container)
			a.chatViewMu.RLock()
			msgData := a.temporaryMessagesStorage[fullAddr][i]
			a.chatViewMu.RUnlock()

			c.Objects[0].(*widget.Entry).SetText(msgData.Text.Text)
			if msgData.IsMine {
				c.Objects[0].(*widget.Entry).TextStyle = fyne.TextStyle{Italic: true}
			}

			if msgData.Text.Text != "" {
				c.Objects[0].(*widget.Entry).Show()
			} else {
				c.Objects[0].(*widget.Entry).Hide()
			}
			if msgData.File != nil {
				fileTemporary := c.Objects[1].(*messagetype.File_type)
				fileTemporary.CopyFileType(msgData.File)
				fileTemporary.Show()
			} else {
				c.Objects[1].(*messagetype.File_type).Hide()
				c.Objects[1].(*messagetype.File_type).SetText("")
			}
			c.Refresh()
		},
	)

	layout := container.NewBorder(topAppBar, bottomInputBar, nil, nil, a.messageList)

	return a.messageList, layout
}
