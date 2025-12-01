package tls_communication

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func SendPackage(toAddr string, data interface{}) error {
	cert, err := GetOrGenerateCertificate("cert.pem", "key.pem")
	if err != nil {
		return fmt.Errorf("certificate error: %w", err)
	}

	conf := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert}, // Не забываем сертификат!
	}

	conn, err := tls.Dial("tcp", toAddr, conf)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer conn.Close()

	err = json.NewEncoder(conn).Encode(data)
	if err != nil {
		return fmt.Errorf("encode error: %w", err)
	}

	return nil
}

func DownloadFile(peerAddr string, filePath string, savePath string) error {
	cert, err := GetOrGenerateCertificate("cert.pem", "key.pem")
	if err != nil {
		return err
	}

	conf := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	conn, err := tls.Dial("tcp", peerAddr, conf)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer conn.Close()

	req := mail_data{
		Package_type: 2,
		FilePath:     filePath,
	}

	if err := json.NewEncoder(conn).Encode(req); err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	outFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, conn)
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}

	fmt.Printf("Downloaded %d bytes\n", n)
	return nil
}
