package styles

import (
	"fmt"
	"testing"
)

func TestColorGradient(t *testing.T) {
	maxRating := 2.0

	fmt.Println("\nЦветовой градиент (должен совпадать с HTML):")
	fmt.Println("Rating | Normalized | Color")
	fmt.Println("-------|------------|-------")

	for rating := 0.0; rating <= maxRating+0.01; rating += 0.2 {
		color := GetRatingColor(rating, maxRating)
		normalized := rating / maxRating
		fmt.Printf("%.1f    | %.2f       | %s\n", rating, normalized, color)
	}

	// Проверяем граничные случаи
	t.Run("Zero rating", func(t *testing.T) {
		color := GetRatingColor(0, maxRating)
		// При t=0: h=220/360=0.611, s=0.85, l=0.16
		// Должен быть темно-синий
		t.Logf("Zero rating color: %s (should be dark blue)", color)
	})

	t.Run("Max rating", func(t *testing.T) {
		color := GetRatingColor(maxRating, maxRating)
		// При t=1: h=0, s=0.85, l=0.55
		// Должен быть красный
		t.Logf("Max rating color: %s (should be red)", color)
	})

	t.Run("Mid rating", func(t *testing.T) {
		color := GetRatingColor(1.0, maxRating)
		// При t=0.5: h=110/360=0.306, s=0.85, l=0.355
		// Должен быть зеленоватый
		t.Logf("Mid rating color: %s (should be greenish)", color)
	})
}
