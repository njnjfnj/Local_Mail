package gui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	tls_communication "github.com/njnjfnj/Local_Mail/internal/local_net/tls_communication"
	udp_broadcast "github.com/njnjfnj/Local_Mail/internal/local_net/udp_broadcast"
	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
)

// AppGUI хранит состояние нашего UI
type AppGUI struct {
	window         fyne.Window
	chatListScreen fyne.CanvasObject
	chatListView   *widget.List
	inputEntry     *widget.Entry

	menuScreen            fyne.CanvasObject
	settingsScreen        fyne.CanvasObject
	settingsScreenWidgets *SettingsWidgets

	chatList           map[string]string
	chatListMu         sync.RWMutex
	updateChatListChan chan string
}

// NewAppGUI создает новый экземпляр нашего UI
func NewAppGUI(w fyne.Window) *AppGUI {
	a := &AppGUI{
		window: w,
	}
	a.chatList = make(map[string]string)
	a.chatListScreen = a.createChatListScreen()
	a.menuScreen = a.createMenuScreen()
	a.settingsScreenWidgets = a.createsettingsWidgets()
	a.settingsScreen = a.createSettingsScreen()
	a.inputEntry = widget.NewEntry()

	a.updateChatListChan = make(chan string)

	udp_broadcast.Start_udp_broadcast_reciver(a.settingsScreenWidgets.Username, a.settingsScreenWidgets.Port, a.settingsScreenWidgets.UdpPort, a.updateChatListChan)
	tls_communication.StartTCPServer(local_net.GetOutboundIP(), a.settingsScreenWidgets.Port, a.updateChatListChan)

	a.startUpdateChatList(a.updateChatListChan, &a.chatListMu)

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
