package gui

import (
	"strings"

	tls_communication "github.com/njnjfnj/Local_Mail/internal/local_net/tls_communication"
)

func (a *AppGUI) startDownloadingFiles(c chan string) {
	go a.downloadFile(c)
}

func (a *AppGUI) downloadFile(ch chan string) {
	for {
		value := <-ch
		value_list := strings.Split(value, "~")

		go tls_communication.DownloadFile(value_list[0], value_list[1], value_list[2])
	}
}
