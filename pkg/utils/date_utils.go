package utils

import (
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", nil
	}

	// Парсим исходную дату
	currentDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	// Обрабатываем разные форматы правил повторения
	switch {
	case strings.HasPrefix(repeat, "d "): // Повтор каждые N дней
		return handleDailyRepeat(currentDate, repeat)

	case strings.HasPrefix(repeat, "w "): // Повтор по дням недели
		return handleWeeklyRepeat(currentDate, repeat)

	case strings.HasPrefix(repeat, "m "): // Повтор по дням месяца
		return handleMonthlyRepeat(currentDate, repeat)

	case repeat == "y": // Ежегодно
		return handleYearlyRepeat(currentDate)

	default:
		return "", nil
	}
}

// handleDailyRepeat обрабатывает повторение каждые N дней
func handleDailyRepeat(currentDate time.Time, repeat string) (string, error) {
	daysStr := strings.TrimPrefix(repeat, "d ")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		return "", err
	}
	nextDate := currentDate.AddDate(0, 0, days)
	return nextDate.Format("20060102"), nil
}

// handleWeeklyRepeat обрабатывает повторение по дням недели
func handleWeeklyRepeat(currentDate time.Time, repeat string) (string, error) {
	daysStr := strings.TrimPrefix(repeat, "w ")
	daysList := strings.Split(daysStr, ",")

	if len(daysList) == 0 {
		return "", nil
	}

	// Конвертируем дни недели в time.Weekday
	var weekDays []time.Weekday
	for _, dayStr := range daysList {
		dayNum, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil || dayNum < 1 || dayNum > 7 {
			return "", err
		}

		// Преобразуем: 1=воскресенье, 2=понедельник, ..., 7=суббота
		weekday := time.Weekday(dayNum - 1)
		weekDays = append(weekDays, weekday)
	}

	// Находим следующий подходящий день (в пределах 7 дней)
	for i := 1; i <= 7; i++ {
		nextDate := currentDate.AddDate(0, 0, i)
		for _, targetDay := range weekDays {
			if nextDate.Weekday() == targetDay {
				return nextDate.Format("20060102"), nil
			}
		}
	}

	return "", nil
}

// handleMonthlyRepeat обрабатывает повторение по дням месяца
func handleMonthlyRepeat(currentDate time.Time, repeat string) (string, error) {
	dayStr := strings.TrimPrefix(repeat, "m ")
	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 || day > 31 {
		return "", err
	}

	// Получаем год и месяц текущей даты
	year := currentDate.Year()
	month := currentDate.Month()

	// Пытаемся создать дату в текущем месяце
	testDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	// Если день существует в текущем месяце и дата еще не прошла
	if testDate.Day() == day && testDate.After(currentDate) {
		return testDate.Format("20060102"), nil
	}

	// Если день не существует в текущем месяце, используем последний день текущего месяца
	if testDate.Day() != day {
		lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
		if lastDayOfMonth.After(currentDate) {
			return lastDayOfMonth.Format("20060102"), nil
		}
	}

	// Если текущий месяц уже прошел, переходим к следующему месяцу
	month++
	if month > 12 {
		month = 1
		year++
	}

	// Пытаемся создать дату в следующем месяце
	nextDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	// Если день превышает количество дней в месяце, берем последний день
	if nextDate.Day() != day {
		nextDate = time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	}

	return nextDate.Format("20060102"), nil
}

// handleYearlyRepeat обрабатывает ежегодное повторение
func handleYearlyRepeat(currentDate time.Time) (string, error) {
	originalDay := currentDate.Day()
	originalMonth := currentDate.Month()

	// Для даты 29 февраля ищем следующий високосный год
	if originalMonth == time.February && originalDay == 29 {
		year := currentDate.Year() + 1
		for {
			if isLeapYear(year) {
				return time.Date(year, time.February, 29, 0, 0, 0, 0, time.UTC).Format("20060102"), nil
			}
			year++
			// Защита от бесконечного цикла
			if year > currentDate.Year()+100 {
				return "", nil
			}
		}
	}

	// Для остальных дат просто прибавляем год
	nextYear := currentDate.Year() + 1
	nextDate := time.Date(nextYear, originalMonth, originalDay, 0, 0, 0, 0, time.UTC)

	// Если дата невалидна (например, 31 июня), берем последний день месяца
	if nextDate.Day() != originalDay {
		nextDate = time.Date(nextYear, originalMonth+1, 0, 0, 0, 0, 0, time.UTC)
	}

	return nextDate.Format("20060102"), nil
}

// isLeapYear проверяет, является ли год високосным
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
