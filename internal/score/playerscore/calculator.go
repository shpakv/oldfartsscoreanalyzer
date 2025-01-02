package playerscore

var (
	W1 = 4.0
	W2 = 2.0
	W3 = 1.0
)

type Calculator struct {
}

func (c *Calculator) Apply(kills, deaths, assists, damage float64, opponentCount int) float64 {
	ad := c.calculateAverageDamageByOpponent(damage, opponentCount)
	kdr := c.calculateKillDeathsRatio(kills, deaths)

	return c.calculateWithWeight(ad, W1) +
		c.calculateWithWeight(kdr, W2) +
		c.calculateWithWeight(assists, W3)
}

func (c *Calculator) calculateKillDeathsRatio(kills, deaths float64) float64 {
	return kills / (deaths + 0.001) // avoid zero delimiter
}

func (c *Calculator) calculateAverageDamageByOpponent(damage float64, opponentCount int) float64 {
	return damage / float64(opponentCount)
}

func (c *Calculator) calculateWithWeight(number, weight float64) float64 {
	return number * weight
}
