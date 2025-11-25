package upd_broadcast

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// TODO: сделать выбор порта возможнго айпи в настройках

const broadcastAddress = "255.255.255.255"

func upd_broadcast_test() {
	broadcastPort := "1337"
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

		message := "TEST"
		for {
			_, err = conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("sending error: %w\n", err)
				conn.Close()
				time.Sleep(5 * time.Second)
				break
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func Start_broadcast_test() {
	go upd_broadcast_test()
}

type mail_data struct {
	Package_type int
	Username     string
	FullAddress  string
	Message      string
}

func Send_connect_data_via_broadcast(username, broadcastLocalAddress, broadcastLocalPort, udpBroadcastPort string) error {
	if username == "" {
		return fmt.Errorf("username missing")
	}

	fullAddr := fmt.Sprintf("%s:%s", broadcastAddress, udpBroadcastPort)

	addr, err := net.ResolveUDPAddr("udp4", fullAddr)
	if err != nil {
		return fmt.Errorf("resolving address error: %w", err)

	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}

	defer conn.Close()

	message := mail_data{
		Package_type: 0,
		Username:     username,
		FullAddress:  fmt.Sprintf("%s:%s", broadcastLocalAddress, broadcastLocalPort),
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshalling error: %w", err)
	}

	_, err = conn.Write([]byte(jsonData))
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}

	return nil
}
