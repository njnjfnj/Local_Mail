package gui

import (
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"

	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
)

type SettingsWidgets struct {
	Username *widget.Entry
	Port     *widget.Entry
	UdpPort  *widget.Entry
}

func (a *AppGUI) createsettingsWidgets() *SettingsWidgets {
	portEntry := widget.NewEntry()
	portEntry.SetText("1338")

	udpPortEntry := widget.NewEntry()
	udpPortEntry.SetText("1337")

	return &SettingsWidgets{
		Username: widget.NewEntry(),
		Port:     portEntry,
		UdpPort:  udpPortEntry,
	}
}

type settingsSaving struct {
	Username string
	Port     string
	UdpPort  string
}

func (a *AppGUI) createSettingsScreen() fyne.CanvasObject {
	backButton := widget.NewButton("Back", func() {
		a.navigateToMenu()
	})

	saveButton := widget.NewButton("Save", func() {
		saving, _ := json.Marshal(settingsSaving{
			Username: a.settingsScreenWidgets.Username.Text,
			Port:     a.settingsScreenWidgets.Port.Text,
			UdpPort:  a.settingsScreenWidgets.UdpPort.Text,
		})

		if err := os.WriteFile("settings.json", saving, 0644); err != nil {
			fmt.Println("error settings saving:", err)
		}
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
	formUdpPort := widget.NewFormItem("Udp port: ", a.settingsScreenWidgets.UdpPort)
	formContainer := widget.NewForm(formUsername, formPort, formUdpPort)

	vbox := container.NewVBox(widget.NewLabel(""), formContainer, widget.NewLabel("Your local IP: "+local_net.GetOutboundIP()))

	return container.NewBorder(topAppBar, saveButton, nil, nil, vbox)
}
