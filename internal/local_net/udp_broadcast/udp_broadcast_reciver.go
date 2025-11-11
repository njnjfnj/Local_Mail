package upd_broadcast

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

const listenPort = "1337"

func udp_broadcast_reciver(chatList *map[string]string, chatListView *widget.List, chatListMu *sync.RWMutex) {
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
				var data connect_data
				json.Unmarshal([]byte(message), &data)
				chatListMu.Lock()
				(*chatList)[data.FullAddress] = data.Username
				chatListMu.Unlock()

				fyne.Do(func() {
					chatListView.Refresh()
				})
			}
		}
	}
}

func Start_udp_broadcast_reciver(chatList *map[string]string, chatListView *widget.List, chatListMu *sync.RWMutex) {
	go udp_broadcast_reciver(chatList, chatListView, chatListMu)
}
