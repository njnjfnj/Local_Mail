package gui

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	messagetype "github.com/njnjfnj/Local_Mail/gui/message_type"
	tls_communication "github.com/njnjfnj/Local_Mail/internal/local_net/tls_communication"
	udp_broadcast "github.com/njnjfnj/Local_Mail/internal/local_net/udp_broadcast"
	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
)

// AppGUI хранит состояние нашего UI
type AppGUI struct {
	window         fyne.Window
	chatListScreen fyne.CanvasObject
	chatListView   *widget.List

	inputEntry               *widget.Entry
	temporaryMessagesStorage map[string][]messagetype.Message_type
	updateChatViewChan       chan messagetype.Message_type
	startFileDownloadingChan chan string
	messageList              *widget.List
	chatViewMu               sync.RWMutex

	menuScreen            fyne.CanvasObject
	settingsScreen        fyne.CanvasObject
	settingsScreenWidgets *SettingsWidgets

	chatList           map[string]string
	chatListMu         sync.RWMutex
	updateChatListChan chan string
}

// NewAppGUI создает новый экземпляр нашего UI
func NewAppGUI(w fyne.Window) *AppGUI {
	settings, err := os.ReadFile("settings.json")
	if err != nil {
		fmt.Println("error loading savings:", err)
	}

	var settingsLoad settingsSaving

	err = json.Unmarshal(settings, &settingsLoad)
	if err != nil {
		fmt.Println("error decoding JSON:", err)
	}

	a := &AppGUI{
		window: w,
	}
	a.chatList = make(map[string]string)
	a.chatListScreen = a.createChatListScreen()
	a.menuScreen = a.createMenuScreen()
	a.settingsScreenWidgets = a.createsettingsWidgets()
	a.settingsScreen = a.createSettingsScreen()
	a.inputEntry = widget.NewEntry()
	a.temporaryMessagesStorage = make(map[string][]messagetype.Message_type)

	a.updateChatListChan = make(chan string)
	a.updateChatViewChan = make(chan messagetype.Message_type)
	a.startFileDownloadingChan = make(chan string)

	if string(settings) != "" {
		a.settingsScreenWidgets.Username.SetText(settingsLoad.Username)
		a.settingsScreenWidgets.Port.SetText(settingsLoad.Port)
		a.settingsScreenWidgets.UdpPort.SetText(settingsLoad.UdpPort)
	}

	udp_broadcast.Start_udp_broadcast_reciver(
		a.settingsScreenWidgets.Username,
		a.settingsScreenWidgets.Port,
		a.settingsScreenWidgets.UdpPort,
		a.updateChatListChan,
	)
	tls_communication.StartTCPServer(
		local_net.GetOutboundIP(),
		a.settingsScreenWidgets.Port,
		a.updateChatListChan,
		a.startFileDownloadingChan,
		a.updateChatViewChan,
		a.window,
	)

	a.startUpdateChatList(a.updateChatListChan, &a.chatListMu)
	a.startUpdateChatView(a.updateChatViewChan, &a.chatViewMu)

	return a
}

// CreateMainLayout возвращает стартовый экран приложения
func (a *AppGUI) CreateMainLayout() fyne.CanvasObject {
	return a.chatListScreen
}

// ИЗМЕНЕНИЕ: Эта функция теперь управляет прокруткой
func (a *AppGUI) navigateToChatView(contactName, fullAddr string) {
	messageList, chatScreen := a.createChatViewScreen(contactName, fullAddr)

	a.window.SetContent(chatScreen)

	messageList.ScrollToBottom()
}

func (a *AppGUI) navigateBackToList() {
	a.window.SetContent(a.chatListScreen)
}

func (a *AppGUI) navigateToMenu() {
	a.window.SetContent(a.menuScreen)
}

func (a *AppGUI) navigateToSettings() {
	a.window.SetContent(a.settingsScreen)
}
