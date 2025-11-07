package upd_broadcast

import (
	"fmt"
	"net"
	"time"
)

// TODO: сделать выбор порта возможнго айпи в настройках
const broadcastPort = "1337"
const broadcastAddress = "255.255.255.255"

func upd_broadcast() {
start_end_point:
	for {
		fullAddr := fmt.Sprintf("%s:%s", broadcastAddress, broadcastPort)

		addr, err := net.ResolveUDPAddr("udp4", fullAddr)
		if err != nil {
			fmt.Printf("resolving address error: %w\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		conn, err := net.DialUDP("udp4", nil, addr)
		if err != nil {
			fmt.Printf("connect error: %w\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		message := "HIIIII MY NIGGAS"
		for {
			_, err = conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("sending error: %w\n", err)
				conn.Close()
				time.Sleep(5 * time.Second)
				continue start_end_point
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func Start_broadcast() {
	go upd_broadcast()
}
