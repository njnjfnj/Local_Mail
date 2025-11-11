package upd_broadcast

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	t "github.com/njnjfnj/Local_Mail/internal/local_net/tcp_communication"
)

const listenPort = "1337"

type SettingsWidgets struct {
	Username *widget.Entry
	Port     *widget.Entry
}

func udp_broadcast_reciver(chatList *map[string]string, chatListView *widget.List, chatListMu *sync.RWMutex, s *SettingsWidgets) {
	for {
		addr, err := net.ResolveUDPAddr("udp4", ":"+listenPort)
		if err != nil {
			fmt.Println("resolving address error:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		conn, err := net.ListenUDP("udp4", addr)
		if err != nil {
			//fmt.Println("listen error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		buffer := make([]byte, 1024) // 1024 байт должно хватить для простых сообщений

		for {
			// Читаем данные из соединения
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				conn.Close()
				break
			}

			message := string(buffer[:n])
			fmt.Printf("Получено '%s' от %s\n", message, remoteAddr)

			messageType := int(message[16])

			switch messageType {
			case int('0'):
				var data Connect_data
				json.Unmarshal([]byte(message), &data)

				if data.FullAddress == fmt.Sprint(GetOutboundIP()+":"+s.Port.Text) {
					continue
				}

				chatListMu.Lock()
				(*chatList)[data.FullAddress] = data.Username
				chatListMu.Unlock()

				message := Connect_data{
					Package_type: 0,
					Username:     s.Username.Text,
					FullAddress:  fmt.Sprintf("%s:%s", GetOutboundIP(), s.Port.Text),
				}

				jsonData, err := json.Marshal(message)
				if err != nil {
					continue
				}

				t.SendConnectData(data.FullAddress, jsonData)

				fyne.Do(func() {
					chatListView.Refresh()
				})

			}
		}
	}
}

func Start_udp_broadcast_reciver(chatList *map[string]string, chatListView *widget.List, chatListMu *sync.RWMutex, s *SettingsWidgets) {
	go udp_broadcast_reciver(chatList, chatListView, chatListMu, s)
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
