package gui

import (
	"fmt"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"

	udp_broadcast "github.com/njnjfnj/Local_Mail/internal/local_net/udp_broadcast"
)

// type settingsWidgets struct {
// 	Username *widget.Entry
// 	Port     *widget.Entry
// }

func (a *AppGUI) createsettingsWidgets() *udp_broadcast.SettingsWidgets {
	portEntry := widget.NewEntry()
	portEntry.SetText("1338")

	return &udp_broadcast.SettingsWidgets{
		Username: widget.NewEntry(),
		Port:     portEntry,
	}
}

func (a *AppGUI) createSettingsScreen() fyne.CanvasObject {
	backButton := widget.NewButton("Back", func() {
		a.navigateToMenu()
	})

	saveButton := widget.NewButton("Save", func() {

	})

	title := widget.NewLabel("Settings")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	topAppBar := container.NewBorder(
		nil, nil, // top, bottom
		backButton, // left
		nil,        // right
		title,      // center
	)
	formUsername := widget.NewFormItem("Username: ", a.settingsScreenWidgets.Username)
	formPort := widget.NewFormItem("Port: ", a.settingsScreenWidgets.Port)
	formContainer := widget.NewForm(formUsername, formPort)

	vbox := container.NewVBox(widget.NewLabel(""), formContainer, widget.NewLabel("Your local IP: "+GetOutboundIP()))

	return container.NewBorder(topAppBar, saveButton, nil, nil, vbox)
}

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return fmt.Sprintf("Error occurred: %s", err.Error())
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
