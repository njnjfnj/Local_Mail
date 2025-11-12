package local_net

import (
	"fmt"
	"net"
)

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return fmt.Sprintf("Error occurred: %s", err.Error())
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
