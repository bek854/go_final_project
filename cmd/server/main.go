package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bek854/go_final_project/internal/api"
	"github.com/bek854/go_final_project/pkg/db"
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
	defer db.Close()

	// Создаем обработчик
	handler := api.NewHandler(db.DB)

	// Настраиваем обработчик для статических файлов
	fileServer := http.FileServer(http.Dir("./web"))

	// Настраиваем маршрутизацию
	http.Handle("/", http.StripPrefix("/", fileServer))
	http.HandleFunc("/api/nextdate", handler.NextDateHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTasks(w, r)
		case http.MethodPost:
			handler.CreateTask(w, r)
		default:
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/task/update", handler.UpdateTask)
	http.HandleFunc("/api/task/delete", handler.DeleteTask)
	http.HandleFunc("/api/task/done", handler.DoneTask)

	// Запускаем сервер
	log.Printf("Сервер запущен на порту :%s", port)
	log.Printf("База данных: %s", dbFile)
	log.Printf("Откройте http://localhost:%s в браузере", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
