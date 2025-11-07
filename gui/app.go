package gui

import (
	"fyne.io/fyne/v2"
	udp_broadcast "github.com/njnjfnj/Local_Mail/internal/local_net/udp_broadcast"
)

// AppGUI хранит состояние нашего UI
type AppGUI struct {
	window         fyne.Window
	chatListScreen fyne.CanvasObject
}

// NewAppGUI создает новый экземпляр нашего UI
func NewAppGUI(w fyne.Window) *AppGUI {
	a := &AppGUI{
		window: w,
	}
	a.chatListScreen = a.createChatListScreen()

	udp_broadcast.Start_broadcast()
	udp_broadcast.Start_udp_broadcast_reciver()

	return a
}

// CreateMainLayout возвращает стартовый экран приложения
func (a *AppGUI) CreateMainLayout() fyne.CanvasObject {
	return a.chatListScreen
}

// ИЗМЕНЕНИЕ: Эта функция теперь управляет прокруткой
func (a *AppGUI) navigateToChatView(contactName string) {
	messageList, chatScreen := a.createChatViewScreen(contactName)

	a.window.SetContent(chatScreen)

	messageList.ScrollToBottom()
}

func (a *AppGUI) navigateBackToList() {
	a.window.SetContent(a.chatListScreen)
}
