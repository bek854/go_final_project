package utils

import (
	"testing"
	"time"
)

func TestNextDate(t *testing.T) {
	now := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		date     string
		repeat   string
		expected string
		wantErr  bool
	}{
		// Daily repeats
		{"d 3", "20231201", "d 3", "20231204", false},
		{"d 1", "20231201", "d 1", "20231202", false},
		{"Неверное daily правило", "20231201", "d 0", "", true},
		{"Неверный формат daily", "20231201", "d abc", "", true},

		// Weekly repeats
		{"w 1 (воскресенье)", "20231201", "w 1", "20231203", false},
		{"w 2 (понедельник)", "20231201", "w 2", "20231204", false},
		{"w 1,3,5", "20231201", "w 1,3,5", "20231203", false},
		{"Неверный день недели", "20231201", "w 8", "", true},
		{"Неверный формат weekly", "20231201", "w abc", "", true},

		// Monthly repeats
		{"m 15", "20231201", "m 15", "20231215", false},
		{"m 31 в феврале", "20230131", "m 31", "20230228", false},
		{"m 31 в апреле", "20230401", "m 31", "20230430", false},
		{"Неверный день месяца", "20231201", "m 32", "", true},

		// Yearly repeats
		{"y 29 февраля 2020", "20200229", "y", "20240229", false},
		{"y 29 февраля 2020 -> 2024", "20200229", "y", "20240229", false},

		// Error cases
		{"Пустое правило", "20231201", "", "", true},
		{"Неверный формат даты", "2023121", "d 1", "", true},
		{"Неизвестное правило", "20231201", "x 1", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NextDate(now, tt.date, tt.repeat)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("NextDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}
