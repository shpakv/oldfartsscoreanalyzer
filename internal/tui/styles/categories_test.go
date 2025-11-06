package styles

import (
	"fmt"
	"testing"
)

func TestGetRatingCategory(t *testing.T) {
	averageMu := 2.541 // Средний рейтинг из repository.go

	tests := []struct {
		rating   float64
		expected string
	}{
		{1.0, CategoryWeak},    // < 2.160
		{2.0, CategoryWeak},    // < 2.160
		{2.2, CategoryAverage}, // >= 2.160, < 2.668
		{2.5, CategoryAverage}, // >= 2.160, < 2.668
		{2.7, CategoryGood},    // >= 2.668, < 3.176
		{3.0, CategoryGood},    // >= 2.668, < 3.176
		{3.2, CategoryMonster}, // >= 3.176
		{4.5, CategoryMonster}, // >= 3.176
	}

	fmt.Println("\nТестирование категорий рейтинга (реальные значения Score):")
	fmt.Println("Рейтинг | Категория    | Ожидаемая категория")
	fmt.Println("--------|--------------|--------------------")

	for _, tt := range tests {
		category, color := GetRatingCategory(tt.rating, averageMu)
		match := "✅"
		if category != tt.expected {
			match = "❌"
			t.Errorf("GetRatingCategory(%.1f, %.3f) = %s, expected %s",
				tt.rating, averageMu, category, tt.expected)
		}
		fmt.Printf("%.1f     | %-12s | %-19s %s (цвет: %s)\n",
			tt.rating, category, tt.expected, match, color)
	}
}

func TestCategoryThresholds(t *testing.T) {
	averageMu := 2.541

	// Проверяем точные пороги из HTML
	weak := averageMu * 0.85    // 2.160
	average := averageMu * 1.05 // 2.668
	monster := averageMu * 1.25 // 3.176

	fmt.Println("\nПороги категорий (как в HTML):")
	fmt.Printf("Подпивас (Mil-Spec):   < %.3f\n", weak)
	fmt.Printf("Пердун (Restricted):   %.3f - %.3f\n", weak, average)
	fmt.Printf("Ебака (Classified):    %.3f - %.3f\n", average, monster)
	fmt.Printf("Гиперебака (Covert):   >= %.3f\n", monster)

	// Проверяем граничные случаи
	t.Run("Граница Подпивас/Пердун", func(t *testing.T) {
		rating1 := weak - 0.01 // 2.159
		rating2 := weak        // 2.160

		cat1, _ := GetRatingCategory(rating1, averageMu)
		cat2, _ := GetRatingCategory(rating2, averageMu)

		if cat1 != CategoryWeak {
			t.Errorf("%.3f должен быть %s, получен %s", rating1, CategoryWeak, cat1)
		}
		if cat2 != CategoryAverage {
			t.Errorf("%.3f должен быть %s, получен %s", rating2, CategoryAverage, cat2)
		}
	})

	t.Run("Граница Пердун/Ебака", func(t *testing.T) {
		rating1 := average - 0.01 // 2.667
		rating2 := average        // 2.668

		cat1, _ := GetRatingCategory(rating1, averageMu)
		cat2, _ := GetRatingCategory(rating2, averageMu)

		if cat1 != CategoryAverage {
			t.Errorf("%.3f должен быть %s, получен %s", rating1, CategoryAverage, cat1)
		}
		if cat2 != CategoryGood {
			t.Errorf("%.3f должен быть %s, получен %s", rating2, CategoryGood, cat2)
		}
	})

	t.Run("Граница Ебака/Гиперебака", func(t *testing.T) {
		rating1 := monster - 0.01 // 3.175
		rating2 := monster        // 3.176

		cat1, _ := GetRatingCategory(rating1, averageMu)
		cat2, _ := GetRatingCategory(rating2, averageMu)

		if cat1 != CategoryGood {
			t.Errorf("%.3f должен быть %s, получен %s", rating1, CategoryGood, cat1)
		}
		if cat2 != CategoryMonster {
			t.Errorf("%.3f должен быть %s, получен %s", rating2, CategoryMonster, cat2)
		}
	})
}
