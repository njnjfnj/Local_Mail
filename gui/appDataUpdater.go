package gui

import (
	"fmt"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
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
