package playerrating

import (
	"math"
	"oldfartscounter/internal/player"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	c := &Calculator{}
	t.Run("TEST method ApplyForRange", func(t *testing.T) {
		t.Run("Test With Hardcode Value", func(t *testing.T) {
			mockPlayer := player.Player{}
			timeRange := time.Now()
			hardcodeValue := 7.0

			result := c.ApplyForRange(mockPlayer, timeRange, 0, hardcodeValue)

			assert.Equalf(t, hardcodeValue, result, "Hardcoded value should be returned")
		})

		t.Run("Test With Recent Value", func(t *testing.T) {
			mockPlayer := player.Player{
				Games: []*player.PlayedGame{
					{
						Date: time.Now().Add(-7 * 1 * 24 * time.Hour),
						Maps: []*player.PlayedMap{
							{
								Score: 25.0,
							},
							{
								Score: 2.5,
							},
							{
								Score: 2.5,
							},
						},
					},
					{
						Date: time.Now().Add(-7 * 2 * 24 * time.Hour).Add(-1 * time.Hour),
						Maps: []*player.PlayedMap{
							{
								Score: 25.0,
							},
							{
								Score: 15.0,
							},
							{
								Score: 5.0,
							},
						},
					},
					{
						Date: time.Now().Add(-7 * 3 * 24 * time.Hour),
						Maps: []*player.PlayedMap{
							{
								Score: 15.0,
							},
							{
								Score: 15.0,
							},
							{
								Score: 15.0,
							},
						},
					},
					{
						Date: time.Now().Add(-7 * 4 * 24 * time.Hour),
						Maps: []*player.PlayedMap{
							{
								Score: 15.0,
							},
							{
								Score: 15.0,
							},
							{
								Score: 15.0,
							},
						},
					},
				},
			}

			timeRange := time.Now().Add(-7 * 3 * 24 * time.Hour).Add(-1 * 24 * time.Hour)
			result := c.ApplyForRange(mockPlayer, timeRange, 0)
			expectedScore := 10.0*0.5 + 15.0*0.3 + 15.0*0.2
			assert.InDelta(t, expectedScore, result, 0.0001,
				"Recent games score calculation should be accurate")
		})

		t.Run("Test With No Games", func(t *testing.T) {
			emptyPlayer := player.Player{
				Games: []*player.PlayedGame{},
			}
			timeRange := time.Now()
			defaultValue := 3.0

			// Act
			result := c.ApplyForRange(emptyPlayer, timeRange, defaultValue)

			// Assert
			assert.Equal(t, defaultValue, result,
				"Should return default value when no games exist")
		})

		t.Run("Test With Older Games", func(t *testing.T) {
			mockPlayer := player.Player{
				Games: []*player.PlayedGame{
					{
						Date: time.Now().Add(-7 * 5 * 24 * time.Hour),
						Maps: []*player.PlayedMap{
							{
								Score: 25.0,
							},
							{
								Score: 2.5,
							},
							{
								Score: 2.5,
							},
						},
					},
					{
						Date: time.Now().Add(-7 * 6 * 24 * time.Hour),
						Maps: []*player.PlayedMap{
							{
								Score: 25.0,
							},
							{
								Score: 15.0,
							},
							{
								Score: 5.0,
							},
						},
					},
				},
			}
			timeRange := time.Now().Add(-7 * 3 * 24 * time.Hour).Add(-1 * 24 * time.Hour)
			result := c.ApplyForRange(mockPlayer, timeRange, 0)

			// Assert
			expectedScore := 10.0*math.Exp(-0.1*5) + 15.0*math.Exp(-0.1*6)
			assert.InDelta(t, expectedScore, result, 0.0001,
				"Older games score calculation with degradation should be accurate")
		})
	})
}
