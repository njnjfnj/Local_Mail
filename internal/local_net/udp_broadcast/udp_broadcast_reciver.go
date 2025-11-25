package upd_broadcast

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"fyne.io/fyne/v2/widget"

	t "github.com/njnjfnj/Local_Mail/internal/local_net/tls_communication"
	local_net "github.com/njnjfnj/Local_Mail/lib/local_net"
)

// const listenPort = "1337"

func udp_broadcast_reciver(Username *widget.Entry, Port *widget.Entry, UdpPort *widget.Entry, ch chan string) {
	for {
		addr, err := net.ResolveUDPAddr("udp4", ":"+UdpPort.Text)
		if err != nil {
			fmt.Println("resolving address error:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		conn, err := net.ListenUDP("udp4", addr)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		buffer := make([]byte, 1024)

		for {
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				conn.Close()
				break
			}

			message := string(buffer[:n])
			fmt.Printf("Получено '%s' от %s\n", message, remoteAddr)

			var data mail_data
			json.Unmarshal([]byte(message), &data)

			switch data.Package_type {
			case 0:
				ch <- fmt.Sprintf("%s~%s", data.Username, data.FullAddress)

				message := mail_data{
					Package_type: 0,
					Username:     Username.Text,
					FullAddress:  fmt.Sprintf("%s:%s", local_net.GetOutboundIP(), Port.Text),
				}

				jsonData, err := json.Marshal(message)
				if err != nil {
					continue
				}

				t.SendConnectData(data.FullAddress, jsonData)
			}
		}
	}
}

func Start_udp_broadcast_reciver(Username *widget.Entry, Port *widget.Entry, UdpPort *widget.Entry, ch chan string) {
	go udp_broadcast_reciver(Username, Port, UdpPort, ch)
}
