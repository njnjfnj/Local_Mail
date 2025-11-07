package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"
)

func (a *AppGUI) createChatListScreen() fyne.CanvasObject {
	menuButton := widget.NewButton("≡", func() {
		// Логика открытия главного меню
	})
	searchButton := widget.NewButton("Find user", func() {
		// Логика поиска
	})

	title := widget.NewLabel("Local Mail")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	topAppBar := container.NewBorder(
		nil, nil, // top, bottom
		menuButton,   // left
		searchButton, // right
		title,        // center
	)

	// 2. Список чатов (используем widget.List для производительности)
	chatData := []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Heidi"} // TODO: make dynamic finding via 192.168.0.255

	// сделать так, чтоб при обновлении ползователей с сети обновлялся chatData
	// и потом обновлялся chatList

	chatList := widget.NewList(
		func() int {
			return len(chatData)
		},
		// Функция, создающая шаблон элемента списка
		func() fyne.CanvasObject {
			nameLabel := widget.NewLabel("Contact Name")
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			lastMessageLabel := widget.NewLabel("Last message snippet...")

			// Просто возвращаем VBox с текстом
			return container.NewVBox(nameLabel, lastMessageLabel)
		},
		// Функция, обновляющая данные в шаблоне
		func(i widget.ListItemID, item fyne.CanvasObject) {
			contactName := chatData[i]
			// Теперь 'item' - это напрямую наш VBox
			vbox := item.(*fyne.Container)

			nameLabel := vbox.Objects[0].(*widget.Label)
			lastMessageLabel := vbox.Objects[1].(*widget.Label)

			nameLabel.SetText(contactName)
			lastMessageLabel.SetText(fmt.Sprintf("Last message from %s...", contactName))
			// сделать что-то типа базы данных, возможно даже sqlite
			// для сохранения переписок, файлов, картинок и тд
		},
	)

	chatList.OnSelected = func(id widget.ListItemID) {
		contactName := chatData[id]
		a.navigateToChatView(contactName)
		chatList.Unselect(id)
	}

	return container.NewBorder(topAppBar, nil, nil, nil, chatList)
}
