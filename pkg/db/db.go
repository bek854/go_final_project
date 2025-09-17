package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// SQL-схема для создания таблицы и индекса
const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX idx_date ON scheduler(date);
`

// Init инициализирует базу данных
func Init(dbFile string) error {
	// Проверяем, нужно ли создавать БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем соединение с БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Проверяем соединение
	if err = DB.Ping(); err != nil {
		return err
	}

	// Если БД не существовала, создаем таблицу и индекс
	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close закрывает соединение с БД
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
