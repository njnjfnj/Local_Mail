package gui

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"

	udp_broadcas "github.com/njnjfnj/Local_Mail/internal/local_net/udp_broadcast"
	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
)

func (a *AppGUI) createChatListScreen() fyne.CanvasObject {
	var vbox1 *fyne.Container

	menuButton := widget.NewButton("≡", func() {
		a.navigateToMenu()
	})
	// searchButton := widget.NewButton("Find user", func() {
	//
	// })

	refreshButton := widget.NewButton("Refresh users", func() {
		vbox1.Objects[0].(*widget.Label).SetText("")
		clear(a.chatList)
		a.chatListView.Refresh()
		if err := udp_broadcas.Send_connect_data_via_broadcast(a.settingsScreenWidgets.Username.Text,
			local_net.GetOutboundIP(),
			a.settingsScreenWidgets.Port.Text,
			a.settingsScreenWidgets.UdpPort.Text); err != nil {
			if a.settingsScreenWidgets.Username.Text == "" || a.settingsScreenWidgets.Port.Text == "" {
				vbox1.Objects[0].(*widget.Label).SetText("CAN NOT REFRESH!!\nGo to menu -> settings -> fill Username & Port")
				return
			}
			vbox1.Objects[0].(*widget.Label).SetText("Error: " + err.Error())
		}
		vbox1.Objects[0].(*widget.Label).SetText("Users have been refreshed!")
	})

	title := widget.NewLabel("Local Mail")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	topAppBar := container.NewBorder(
		nil, nil, // top, bottom
		menuButton, // left
		nil,        // right
		title,      // center
	)

	a.chatListView = widget.NewList(
		func() int {
			a.chatListMu.RLock()
			defer a.chatListMu.RUnlock()
			return len(a.chatList)
		},
		func() fyne.CanvasObject {
			//a.chatListMu.RLock()
			//defer a.chatListMu.Unlock()
			nameLabel := widget.NewLabel("Contact Name")
			nameLabel.TextStyle = fyne.TextStyle{Bold: true}
			userFullAddr := widget.NewLabel("full addr snippet")

			return container.NewVBox(nameLabel, userFullAddr)
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			a.chatListMu.RLock()
			fullAddr, username := GetItemFromMapByIndex(a.chatList, i)
			vbox := item.(*fyne.Container)
			a.chatListMu.RUnlock()

			nameLabel := vbox.Objects[0].(*widget.Label)
			userFullAddr := vbox.Objects[1].(*widget.Label)

			nameLabel.SetText(username)
			userFullAddr.SetText(fullAddr)
			// сделать что-то типа базы данных, возможно даже sqlite
			// для сохранения переписок, файлов, картинок и тд
		},
	)

	a.chatListView.OnSelected = func(id widget.ListItemID) {
		fullAddr, username := GetItemFromMapByIndex(a.chatList, id)
		a.navigateToChatView(username, fullAddr) // добавить контакт фул адрес
		a.chatListView.Unselect(id)
	}

	vbox1 = container.NewVBox(widget.NewLabel(""), refreshButton)

	return container.NewBorder(topAppBar, vbox1, nil, nil, a.chatListView)
}

func GetItemFromMapByIndex(m map[string]string, i int) (string, string) {
	s := make([]string, 0)

	for k := range m {
		s = append(s, k)
	}

	sort.Strings(s)

	var s1 string
	var s2 string

	s1 = s[i]
	s2 = m[s1]

	return s1, s2
}
