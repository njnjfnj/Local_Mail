package gui

import (
	"encoding/json"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	messagetype "github.com/njnjfnj/Local_Mail/gui/message_type"
	tls_communication "github.com/njnjfnj/Local_Mail/internal/local_net/tls_communication"
	udp_broadcast "github.com/njnjfnj/Local_Mail/internal/local_net/udp_broadcast"
	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
)

const PREF_KEY = "app_configuration_json"

type AppGUI struct {
	app            fyne.App
	window         fyne.Window
	chatListScreen fyne.CanvasObject
	chatListView   *widget.List

	inputEntry               *SubmitEntry
	temporaryMessagesStorage map[string][]messagetype.Message_type
	updateChatViewChan       chan messagetype.Message_type
	startFileDownloadingChan chan string
	messageList              *widget.List
	chatViewMu               sync.RWMutex

	menuScreen            fyne.CanvasObject
	aboutScreen           fyne.CanvasObject
	settingsScreen        fyne.CanvasObject
	settingsScreenWidgets *SettingsWidgets

	chatList           map[string]string
	chatListMu         sync.RWMutex
	updateChatListChan chan string
}

func NewAppGUI(w fyne.Window, appapp fyne.App) *AppGUI {
	a := &AppGUI{
		window: w,
		app:    appapp,
	}
	a.chatList = make(map[string]string)
	a.chatListScreen = a.createChatListScreen()
	a.menuScreen = a.createMenuScreen()
	a.aboutScreen = a.createAboutScreen()
	a.settingsScreenWidgets = a.createsettingsWidgets()
	a.settingsScreen = a.createSettingsScreen()
	a.inputEntry = NewSubmitEntry()
	a.inputEntry.Wrapping = fyne.TextWrapWord
	a.temporaryMessagesStorage = make(map[string][]messagetype.Message_type)

	a.updateChatListChan = make(chan string)
	a.updateChatViewChan = make(chan messagetype.Message_type)
	a.startFileDownloadingChan = make(chan string)

	settings := a.app.Preferences().StringWithFallback(PREF_KEY, "")

	var settingsLoad settingsSaving

	err := json.Unmarshal([]byte(settings), &settingsLoad)
	if err != nil {
		fmt.Println("error decoding JSON:", err)
	}
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
		a.app,
	)
	tls_communication.StartTCPServer(
		local_net.GetOutboundIP(),
		a.settingsScreenWidgets.Port,
		a.updateChatListChan,
		a.startFileDownloadingChan,
		a.updateChatViewChan,
		a.window,
		a.app,
	)

	a.startUpdateChatList(a.updateChatListChan, &a.chatListMu)
	a.startUpdateChatView(a.updateChatViewChan, &a.chatViewMu)
	a.startDownloadingFiles(a.startFileDownloadingChan)

	return a
}

func (a *AppGUI) CreateMainLayout() fyne.CanvasObject {
	return a.chatListScreen
}

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

func (a *AppGUI) navigateToAbout() {
	a.window.SetContent(a.aboutScreen)
}
