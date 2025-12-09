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
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	messagetype "github.com/njnjfnj/Local_Mail/gui/message_type"
)

const MaxHeaderBytes = 1024 * 1024

var connectionSemaphore = make(chan struct{}, 100)

func tcpServer(host string, port *widget.Entry, ch, startFileDownloadingChan chan string, ch2 chan messagetype.Message_type, a fyne.Window, app fyne.App) error {
	addr := host + ":" + port.Text

	cert, err := GetOrGenerateCertificate(app)
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

		select {
		case connectionSemaphore <- struct{}{}:
			go func() {
				defer func() { <-connectionSemaphore }()
				handleConnection(conn, ch, startFileDownloadingChan, ch2, a)
			}()
		default:
			log.Println("Server busy, rejecting connection")
			conn.Close()
		}
	}
}

func StartTCPServer(host string, port *widget.Entry, ch, startFileDownloadingChan chan string, ch2 chan messagetype.Message_type, a fyne.Window, app fyne.App) {
	go tcpServer(host, port, ch, startFileDownloadingChan, ch2, a, app)
}

func handleConnection(conn net.Conn, ch, startFileDownloadingChan chan string, ch2 chan messagetype.Message_type, a fyne.Window) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Println("it is not tls connection")
		return
	}

	err := tlsConn.Handshake()
	if err != nil {
		log.Println("handshake error:", err)
		return
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		log.Println("target did not send a certificate")
		return
	}

	peerCert := state.PeerCertificates[0]

	hash := sha256.Sum256(peerCert.Raw)
	fingerprint := hex.EncodeToString(hash[:])

	fmt.Printf("connected user ID: %s\n", fingerprint)

	limitReader := io.LimitReader(conn, MaxHeaderBytes)

	var data mail_data
	if err := json.NewDecoder(limitReader).Decode(&data); err != nil {
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
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		conn.SetWriteDeadline(time.Now().Add(20 * time.Minute)) // TODO: make this changeable in settings
		safePath := filepath.Join("Shared", filepath.Base(data.FilePath))
		file, err := os.Open(safePath)
		if err != nil {
			log.Println("file is not found:", safePath)
			return
		}
		defer file.Close()

		_, err = io.Copy(conn, file)
		if err != nil {
			log.Println("file sending error:", err)
			return
		}

	}
}
