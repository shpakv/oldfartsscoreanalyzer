package playerscore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculator(t *testing.T) {
	c := &Calculator{}
	t.Run("Test method Apply", func(t *testing.T) {
		t.Run("TEST Standard", func(t *testing.T) {
			kills, deaths, assists, damage := 17.00, 9.00, 6.00, 2333.00
			opponentCount := 5
			score := c.Apply(kills, deaths, assists, damage, opponentCount)
			assert.InDelta(t, 1876.1773580713254, score, 1e-05)
		})

		t.Run("TEST Zero deaths case", func(t *testing.T) {
			kills, deaths, assists, damage := 17.00, 0.00, 6.00, 2333.00
			opponentCount := 5
			score := c.Apply(kills, deaths, assists, damage, opponentCount)
			assert.InDelta(t, 35872.4, score, 1e-05)
		})
	})

	t.Run("Test method calculateKillDeathsRatio", func(t *testing.T) {
		t.Run("TEST Standard", func(t *testing.T) {
			kills, deaths := 24.00, 12.00
			kdr := c.calculateKillDeathsRatio(kills, deaths)
			assert.InDelta(t, 1.999833347221065, kdr, 1e-05)
		})
		t.Run("TEST Zero deaths case", func(t *testing.T) {
			kills, deaths := 24.00, 0.00
			kdr := c.calculateKillDeathsRatio(kills, deaths)
			assert.InDelta(t, 24000, kdr, 1e-05)
		})
	})

	t.Run("Test method calculateAverageDamageByOpponent", func(t *testing.T) {
		t.Run("TEST Standard", func(t *testing.T) {
			damage := 2000.00
			opponentCount := 5
			ad := c.calculateAverageDamageByOpponent(damage, opponentCount)
			assert.InDelta(t, 400, ad, 0.00)
		})
	})
}
