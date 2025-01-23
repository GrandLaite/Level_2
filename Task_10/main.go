/* Реализовать утилиту wget с возможностью скачивать сайты целиком. */

package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	var outputDir string
	flag.StringVar(&outputDir, "output", ".", "директория для сохранения страницы")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: wget [--output=directory] URL")
		os.Exit(1)
	}

	startURL := flag.Arg(0)

	err := downloadMainPage(startURL, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки страницы: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Страница успешно загружена!")
}

func downloadMainPage(pageURL, outputDir string) error {
	resp, err := http.Get(pageURL)
	if err != nil {
		return fmt.Errorf("ошибка подключения к %s: %v", pageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("не удалось загрузить %s: статус %d", pageURL, resp.StatusCode)
	}

	u, err := url.Parse(pageURL)
	if err != nil {
		return fmt.Errorf("не удалось разобрать URL %s: %v", pageURL, err)
	}
	fileName := filepath.Join(outputDir, filepath.Base(u.Path))
	if fileName == outputDir || filepath.Ext(fileName) == "" {
		fileName = filepath.Join(outputDir, "index.html")
	}

	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("ошибка создания файла %s: %v", fileName, err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл %s: %v", fileName, err)
	}

	fmt.Printf("Страница загружена: %s -> %s\n", pageURL, fileName)
	return nil
}
