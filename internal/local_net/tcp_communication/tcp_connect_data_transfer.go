package tcpcommunication

import (
	"fmt"
	"net"
)

func SendConnectData(toAddr string, data []byte) error {
	conn, err := net.Dial("tcp", toAddr)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("sending error: %w", err)
	}

	return nil
}
