package gui

import (
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	messagetype "github.com/njnjfnj/Local_Mail/gui/message_type"
	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
)

func (a *AppGUI) startUpdateChatList(c chan string, chatListMu *sync.RWMutex) {
	go a.updateChatList(c, chatListMu)
}

func (a *AppGUI) updateChatList(ch chan string, chatListMu *sync.RWMutex) {
	for {
		value := <-ch
		fmt.Println(value)

		value_list := strings.Split(value, "~")

		if _, v := a.chatList[value_list[1]]; value_list[1] == fmt.Sprint(local_net.GetOutboundIP()+":"+a.settingsScreenWidgets.Port.Text) || v {
			fmt.Println("HALLO")
			continue
		}

		chatListMu.Lock()
		a.chatList[value_list[1]] = value_list[0]
		chatListMu.Unlock()

		fyne.Do(func() {
			a.chatListView.Refresh()
		})

	}
}

func (a *AppGUI) startUpdateChatView(c chan messagetype.Message_type, chatViewMu *sync.RWMutex) {
	go a.updateChatView(c, chatViewMu)
}

func (a *AppGUI) updateChatView(ch chan messagetype.Message_type, chatViewMu *sync.RWMutex) {
	for {
		value := <-ch
		fmt.Println("Updater", value.Holdername)

		chatViewMu.Lock()
		a.temporaryMessagesStorage[value.Holdername] = append(a.temporaryMessagesStorage[value.Holdername], value)
		chatViewMu.Unlock()

		fyne.Do(func() {
			if a.messageList != nil {
				a.messageList.Refresh()
				a.messageList.ScrollToBottom()
			}
		})
	}
}
