package tests

import (
	"testing"
	
	"github.com/bek854/go_final_project/pkg/db"
)

func TestDBConnection(t *testing.T) {
	err := db.Init("test.db")
	if err != nil {
		t.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close()
	
	// Проверяем, что таблица существует
	_, err = db.DB.Exec("SELECT count(*) FROM scheduler")
	if err != nil {
		t.Errorf("Ошибка доступа к таблице: %v", err)
	}
}
