/*
Реализовать простейший telnet-клиент.
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Требования
1.Программа должна подключаться к указанному хосту
(ip или доменное имя + порт) по протоколу TCP.
После подключения STDIN программы должен записываться в сокет,
а данные полученные из сокета должны выводиться в STDOUT.
2.Опционально в программу можно передать таймаут на подключение к серверу
(через аргумент --timeout, по умолчанию 10s).
3.При нажатии Ctrl+D программа должна закрывать сокет и завершаться.
Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	var timeout string
	flag.StringVar(&timeout, "timeout", "10s", "таймаут на подключение (например, 5s, 1m)")
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: go-telnet [--timeout=10s] host port")
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)

	duration, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: неверный формат таймаута: %v\n", err)
		os.Exit(1)
	}

	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, duration)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка подключения: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Подключено к %s\n", address)
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			fmt.Fprintf(os.Stderr, "Соединение закрыто сервером: %v\n", err)
		}
		done <- struct{}{}
	}()

	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка записи в сокет: %v\n", err)
		}
		done <- struct{}{}
	}()

	<-done
	fmt.Println("\nСоединение закрыто. Завершение программы.")
}
