package tls_communication

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	messagetype "github.com/njnjfnj/Local_Mail/gui/message_type"
)

func tcpServer(host string, port *widget.Entry, ch, startFileDownloadingChan chan string, ch2 chan messagetype.Message_type, a fyne.Window) error {
	addr := host + ":" + port.Text

	cert, err := GetOrGenerateCertificate("cert.pem", "key.pem")
	if err != nil {
		return fmt.Errorf("cettificate error %s: %w", addr, err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAnyClientCert,
	}

	listener, err := tls.Listen("tcp", ":"+port.Text, config)
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

		go handleConnection(conn, ch, startFileDownloadingChan, ch2, a)
	}
}

func StartTCPServer(host string, port *widget.Entry, ch, startFileDownloadingChan chan string, ch2 chan messagetype.Message_type, a fyne.Window) {
	go tcpServer(host, port, ch, startFileDownloadingChan, ch2, a)
}

func handleConnection(conn net.Conn, ch, startFileDownloadingChan chan string, ch2 chan messagetype.Message_type, a fyne.Window) {
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

	// message := string(buffer[:n])
	// fmt.Printf("Принято сообщение от %s: %s\n", conn.RemoteAddr().String(), message)

	var data mail_data
	if err := json.NewDecoder(conn).Decode(&data); err != nil {
		log.Println("can not decode json")
		return
	}

	switch data.Package_type {
	case PackageTypeHandshake:
		ch <- fmt.Sprintf("%s~%s", data.Username, data.FullAddress)
	case PackageTypeMessage:
		ch2 <- *messagetype.New_text_message(data.FullAddress, data.Message)
	case PackageTypeSendFileInfo:
		ch2 <- *messagetype.New_message(data.FullAddress, "", data.FilePath, "", a, startFileDownloadingChan)
	case PackageTypeFileReq:
		safePath := filepath.Join("Shared", filepath.Base(data.FilePath))

		file, err := os.Open(safePath)
		if err != nil {
			log.Println("file is not found:", safePath)
		}
		defer file.Close()

		_, err = io.Copy(conn, file)
		if err != nil {
			log.Println("Ошибка отправки файла:", err)
		}

	}
}
