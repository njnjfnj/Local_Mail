package upd_broadcast

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"fyne.io/fyne/v2/widget"

	t "github.com/njnjfnj/Local_Mail/internal/local_net/tcp_communication"
	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
)

const listenPort = "1337"

func udp_broadcast_reciver(Username *widget.Entry, Port *widget.Entry, ch chan string) {
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

				ch <- fmt.Sprintf("%s~%s", data.Username, data.FullAddress)
				// if _, v := (*chatList)[data.FullAddress]; data.FullAddress == fmt.Sprint(local_net.GetOutboundIP()+":"+Port.Text) || !v {
				// 	continue
				// }

				// chatListMu.Lock()
				// (*chatList)[data.FullAddress] = data.Username
				// chatListMu.Unlock()

				message := Connect_data{
					Package_type: 0,
					Username:     Username.Text,
					FullAddress:  fmt.Sprintf("%s:%s", local_net.GetOutboundIP(), Port.Text),
				}

				jsonData, err := json.Marshal(message)
				if err != nil {
					continue
				}

				t.SendConnectData(data.FullAddress, jsonData)

				// fyne.Do(func() {
				// 	chatListView.Refresh()
				// })

			}
		}
	}
}

func Start_udp_broadcast_reciver(Username *widget.Entry, Port *widget.Entry, ch chan string) {
	go udp_broadcast_reciver(Username, Port, ch)
}
