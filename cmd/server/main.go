package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bek854/go_final_project/pkg/db" // Импорт вашего пакета!
)

func main() {
	// Получаем порт из переменной окружения или используем по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Получаем путь к БД из переменной окружения или используем по умолчанию
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализируем базу данных
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close() // Закрываем соединение при выходе

	// Настраиваем обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("./web"))
	http.Handle("/", http.StripPrefix("/", fileServer))

	// Запускаем сервер
	log.Printf("Сервер запущен на порту :%s", port)
	log.Printf("База данных: %s", dbFile)
	log.Printf("Откройте http://localhost:%s в браузере", port)
	
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

