package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	// 1. Добавляем себя в автозагрузку (один раз)
	// Для go run . используем путь к main.go (т.е. текущая папка + main.go)
	addToStartup()

	// 2. Запускаем sing-box
	runSingBox()

	// Держим приложение живым
	fmt.Println("\nWisp работает. sing-box запущен.")
	fmt.Println("Нажмите Ctrl+C для остановки...")
	select {} // бесконечный цикл
}

func addToStartup() {
	// Для go run . берём путь к текущей папке + "main.go"
	// Но в автозагрузку лучше добавлять готовый .exe, поэтому пока оставим как заглушку
	// Позже, когда будешь готов билдить, вернём нормальный код
	fmt.Println("Автозагрузка пока отключена (go run режим).")
	fmt.Println("Когда соберёшь wisp.exe — включим её обратно.")
}

func runSingBox() {
	// Берём текущую рабочую директорию (откуда запущен go run .)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Не удалось получить текущую директорию: %v", err)
	}

	singBoxPath := filepath.Join(wd, "sing-box.exe")
	configPath := filepath.Join(wd, "config.json")

	// Проверяем наличие файлов
	if _, err := os.Stat(singBoxPath); os.IsNotExist(err) {
		log.Fatalf("sing-box.exe не найден в папке проекта!\nПоложите его в %s", wd)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config.json не найден в папке проекта!\nПоложите его в %s", wd)
	}

	fmt.Printf("Запускаю sing-box из: %s\n", singBoxPath)
	fmt.Printf("Конфиг: %s\n", configPath)

	cmd := exec.Command(singBoxPath, "run", "-c", configPath)

	// Скрываем окно консоли sing-box
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Не удалось запустить sing-box: %v", err)
	}

	// Перезапуск, если sing-box упадёт
	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Printf("sing-box завершился с ошибкой: %v", err)
			fmt.Println("Через 10 секунд будет попытка перезапуска...")
			time.Sleep(10 * time.Second)
			runSingBox() // рекурсия
		}
	}()
}
