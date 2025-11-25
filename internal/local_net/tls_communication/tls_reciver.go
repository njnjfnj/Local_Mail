package tls_communication

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"fyne.io/fyne/v2/widget"
)

func tcpServer(host string, port *widget.Entry, ch chan string) error {
	addr := host + ":" + port.Text

	cert, err := GetOrGenerateCertificate("server.crt", "server.key")
	if err != nil {
		return fmt.Errorf("cettificate error %s: %w", addr, err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
	}

	listener, err := tls.Listen("tcp", port.Text, config)
	if err != nil {
		return fmt.Errorf("listen error %s: %w", addr, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
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

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("Это не TLS соединение!")
		return
	}

	err := tlsConn.Handshake()
	if err != nil {
		log.Println("Ошибка рукопожатия:", err)
		return
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		log.Println("Собеседник не прислал сертификат")
		return
	}

	peerCert := state.PeerCertificates[0]

	hash := sha256.Sum256(peerCert.Raw)
	fingerprint := hex.EncodeToString(hash[:])

	fmt.Printf("Подключился пользователь с ID: %s\n", fingerprint)

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	message := string(buffer[:n])
	fmt.Printf("Принято сообщение от %s: %s\n", conn.RemoteAddr().String(), message)

	var data mail_data
	json.Unmarshal([]byte(message), &data)

	switch data.Package_type {
	case 0:
		ch <- fmt.Sprintf("%s~%s", data.Username, data.FullAddress)
	}
	// response := []byte("Получено: " + message)
	// conn.Write(response)
}
