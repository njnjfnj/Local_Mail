package tcp_communication

import (
	"encoding/json"
	"fmt"
	"net"

	"fyne.io/fyne/v2/widget"
)

type Connect_data struct {
	Package_type int
	Username     string
	FullAddress  string
}

func tcpServer(host string, port *widget.Entry, ch chan string) error {
	addr := host + ":" + port.Text

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen error %s: %w", addr, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn, ch)
	}
}

func StartTCPServer(host string, port *widget.Entry, ch chan string) {
	go tcpServer(host, port, ch)
}

func handleConnection(conn net.Conn, ch chan string) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	message := string(buffer[:n])
	fmt.Printf("Принято сообщение от %s: %s\n", conn.RemoteAddr().String(), message)

	messageType := int(message[16])

	switch messageType {
	case int('0'):
		var data Connect_data
		json.Unmarshal([]byte(message), &data)

		ch <- fmt.Sprintf("%s~%s", data.Username, data.FullAddress)
	}
	// response := []byte("Получено: " + message)
	// conn.Write(response)
}
