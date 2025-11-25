package tls_communication

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
)

func SendConnectData(toAddr string, data interface{}) error {
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
