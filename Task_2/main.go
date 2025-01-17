/*
Создать программу печатающую точное время с использованием NTP-библиотеки.
Инициализировать как go module. Использовать библиотеку github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.
Требования
1)Программа должна быть оформлена как go module.
2)Программа должна корректно обрабатывать ошибки библиотеки: выводить их в STDERR и
возвращать ненулевой код выхода в OS.
*/

package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

func main() {
	currentTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка получения времени: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Точное время:", currentTime)
}
