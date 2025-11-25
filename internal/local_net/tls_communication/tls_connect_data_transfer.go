package tls_communication

import (
	"crypto/tls"
	"fmt"
)

func SendConnectData(toAddr string, data []byte) error {
	cert, err := GetOrGenerateCertificate("cert.pem", "key.pem")
	if err != nil {
		return fmt.Errorf("certificate error: %w", err)
	}

	conf := &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	conn, err := tls.Dial("tcp", toAddr, conf)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, string(data)+"\n")

	return nil
}
