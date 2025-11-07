package upd_broadcast

import (
	"fmt"
	"net"
	"time"
)

const listenPort = "1337"

func udp_broadcast_reciver() {
	for {
		fmt.Printf("Запуск UDP-слушателя на порту :%s\n", listenPort)

		// Мы используем "udp4", чтобы явно указать IPv4.
		// Адрес ":12345" означает "слушать на всех IP-адресах (0.0.0.0) на порту 12345".
		addr, err := net.ResolveUDPAddr("udp4", ":"+listenPort)
		if err != nil {
			fmt.Println("resolving address error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Создаем "слушателя" пакетов
		conn, err := net.ListenUDP("udp4", addr)
		if err != nil {
			fmt.Println("listen error:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		//defer conn.Close()

		// Создаем буфер для чтения данных
		buffer := make([]byte, 1024) // 1024 байт должно хватить для простых сообщений

		fmt.Println("Ожидание broadcast-сообщений...")

		// Бесконечный цикл для чтения пакетов
		for {
			// Читаем данные из соединения
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Ошибка чтения:", err)
				continue // Не выходим из цикла, просто ждем следующий пакет
			}

			// Выводим полученное сообщение и адрес отправителя
			message := string(buffer[:n])
			fmt.Printf("Получено '%s' от %s\n", message, remoteAddr)
		}
	}
}

func Start_udp_broadcast_reciver() {
	go udp_broadcast_reciver()
}
