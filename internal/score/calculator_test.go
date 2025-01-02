package score

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRatingCalculator_CalculateScore(t *testing.T) {
	c := &RatingCalculator{}
	t.Run("TEST Standard", func(t *testing.T) {
		kills, deaths, assists, damage := 17.00, 9.00, 6.00, 2333.00
		opponentCount := 5
		score := c.CalculateScore(kills, deaths, assists, damage, opponentCount)
		assert.Equal(t, score, 2365)
	})

	t.Run("TEST Zero deaths case", func(t *testing.T) {
		kills, deaths, assists, damage := 17.00, 0.00, 6.00, 2333.00
		opponentCount := 5
		score := c.CalculateScore(kills, deaths, assists, damage, opponentCount)
		assert.Equal(t, score, 2365)
	})

}
