package score

var (
	W1 = 4.0
	W2 = 2.0
	W3 = 1.0
)

type RatingCalculator struct {
}

func (c *RatingCalculator) CalculateScore(kills, deaths, assists, damage float64, opponentCount int) float64 {
	kdRatio := kills / (deaths + 0.001) // avoid zero delimiter

	return (damage/float64(opponentCount))*W1 +
		kdRatio*W2 +
		assists*W3
}
