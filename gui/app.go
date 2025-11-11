package gui

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	udp_broadcast "github.com/njnjfnj/Local_Mail/internal/local_net/udp_broadcast"
)

// AppGUI хранит состояние нашего UI
type AppGUI struct {
	window                fyne.Window
	chatListScreen        fyne.CanvasObject
	chatListView          *widget.List
	menuScreen            fyne.CanvasObject
	settingsScreen        fyne.CanvasObject
	settingsScreenWidgets *udp_broadcast.SettingsWidgets
	chatList              map[string]string
	chatListMu            sync.RWMutex
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

	udp_broadcast.Start_udp_broadcast_reciver(&a.chatList, a.chatListView, &a.chatListMu, a.settingsScreenWidgets)

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
