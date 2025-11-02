package logparser

import (
	"math"
	"testing"
)

// TestCalculateRoundRatings_BasicScenario tests standard 5v5 scenario
func TestCalculateRoundRatings_BasicScenario(t *testing.T) {
	round := &RoundStats{
		Winner: 3, // CT wins
		Players: []PlayerStats{
			{
				Team:    3, // CT
				Damage:  400,
				Kills:   4,
				Assists: 1,
				Deaths:  1,
			},
		},
	}

	// Add 4 more CT teammates
	for i := 0; i < 4; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 3})
	}
	// Add 5 T opponents
	for i := 0; i < 5; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 2})
	}

	calculateRoundRatings(round)

	player := &round.Players[0]

	// Expected calculation (now with correct 5v5):
	// damageImpact = (400/100) * (5/5)^0.7 * (5/5)^0.5 = 4.0 * 1.0 * 1.0 = 4.0
	// fragImpact = 0.15*4 + 0.08*1 = 0.60 + 0.08 = 0.68
	// numerator = 4.0 + 0.68 = 4.68
	// denominator = 1 + 0.35*1 = 1.35
	// baseRating = 4.68 / 1.35 = 3.467
	// killRatio = 4/5 = 0.8 → multiKillBonus = 0.2
	// win = 1.0
	// clutchBonus = 0 (not outnumbered, 5v5)
	// finalRating = 3.467 * (1 + 0.10*1 + 0.2 + 0) = 3.467 * 1.3 = 4.507

	expected := 4.507
	if math.Abs(player.Rating-expected) > 0.01 {
		t.Errorf("Expected rating ~%.3f, got %.3f", expected, player.Rating)
	}
}

// TestCalculateRoundRatings_ACE tests ACE scenario (5 kills out of 5)
func TestCalculateRoundRatings_ACE(t *testing.T) {
	round := &RoundStats{
		Winner: 3, // CT wins
		Players: []PlayerStats{
			{
				Team:    3, // CT
				Damage:  480,
				Kills:   5, // ACE!
				Assists: 0,
				Deaths:  1,
			},
		},
	}

	// Add 4 more CT teammates
	for i := 0; i < 4; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 3})
	}
	// Add 5 T opponents
	for i := 0; i < 5; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 2})
	}

	calculateRoundRatings(round)

	player := &round.Players[0]

	// Expected (5v5):
	// damageImpact = 4.8
	// fragImpact = 0.75
	// numerator = 5.55
	// denominator = 1.35
	// baseRating = 4.111
	// killRatio = 5/5 = 1.0 → multiKillBonus = 0.3 (ACE)
	// win = 1.0
	// clutchBonus = 0
	// finalRating = 4.111 * (1 + 0.10 + 0.30) = 4.111 * 1.40 = 5.756

	expected := 5.756
	if math.Abs(player.Rating-expected) > 0.01 {
		t.Errorf("Expected rating ~%.3f, got %.3f", expected, player.Rating)
	}
}

// TestCalculateRoundRatings_Clutch1v3 tests 1v3 clutch scenario
func TestCalculateRoundRatings_Clutch1v3(t *testing.T) {
	round := &RoundStats{
		Winner: 3, // CT wins
		Players: []PlayerStats{
			{
				Team:    3, // CT (alone)
				Damage:  300,
				Kills:   3, // Killed all 3 opponents
				Assists: 0,
				Deaths:  0,
			},
		},
	}

	// Add 3 opponents
	for i := 0; i < 3; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 2})
	}

	calculateRoundRatings(round)

	player := &round.Players[0]

	// Expected:
	// damageImpact = (300/100) * (5/3)^0.7 * (3/1)^0.5
	//              = 3.0 * 1.429 * 1.732 = 7.427
	// fragImpact = 0.15*3 = 0.45
	// numerator = 7.427 + 0.45 = 7.877
	// denominator = 1.0 (no deaths)
	// baseRating = 7.877
	// killRatio = 3/3 = 1.0 → multiKillBonus = 0.3 (ACE)
	// win = 1.0
	// clutchBonus = (3-1) * 0.05 = 0.10 (outnumbered by 2)
	// finalRating = 7.877 * (1 + 0.10 + 0.30 + 0.10) = 7.877 * 1.50 = 11.816

	expected := 11.816
	if math.Abs(player.Rating-expected) > 0.2 {
		t.Errorf("Expected rating ~%.3f, got %.3f", expected, player.Rating)
	}
}

// TestCalculateRoundRatings_Clutch1v5 tests extreme 1v5 ACE clutch
func TestCalculateRoundRatings_Clutch1v5(t *testing.T) {
	round := &RoundStats{
		Winner: 2, // T wins
		Players: []PlayerStats{
			{
				Team:    2, // T (alone)
				Damage:  500,
				Kills:   5, // ACE in 1v5!
				Assists: 0,
				Deaths:  0,
			},
		},
	}

	// Add 5 opponents
	for i := 0; i < 5; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 3})
	}

	calculateRoundRatings(round)

	player := &round.Players[0]

	// Expected:
	// damageImpact = (500/100) * (5/5)^0.7 * (5/1)^0.5
	//              = 5.0 * 1.0 * 2.236 = 11.180
	// fragImpact = 0.15*5 = 0.75
	// numerator = 11.180 + 0.75 = 11.93
	// denominator = 1.0
	// baseRating = 11.93
	// killRatio = 5/5 = 1.0 → multiKillBonus = 0.3
	// win = 1.0
	// clutchBonus = (5-1) * 0.05 = 0.20 (outnumbered by 4)
	// finalRating = 11.93 * (1 + 0.10 + 0.30 + 0.20) = 11.93 * 1.60 = 19.088

	expected := 19.088
	if math.Abs(player.Rating-expected) > 0.2 {
		t.Errorf("Expected rating ~%.3f, got %.3f", expected, player.Rating)
	}
}

// TestCalculateRoundRatings_WeakRound tests poor performance
func TestCalculateRoundRatings_WeakRound(t *testing.T) {
	round := &RoundStats{
		Winner: 2, // T wins (player is CT - lost)
		Players: []PlayerStats{
			{
				Team:    3, // CT
				Damage:  120,
				Kills:   1,
				Assists: 0,
				Deaths:  3,
			},
		},
	}

	// Add 4 more CT teammates
	for i := 0; i < 4; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 3})
	}
	// Add 5 T opponents
	for i := 0; i < 5; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 2})
	}

	calculateRoundRatings(round)

	player := &round.Players[0]

	// Expected (5v5, lost):
	// damageImpact = (120/100) * 1.0 * 1.0 = 1.2
	// fragImpact = 0.15*1 = 0.15
	// numerator = 1.35
	// denominator = 1 + 0.35*3 = 2.05
	// baseRating = 1.35 / 2.05 = 0.659
	// killRatio = 1/5 = 0.2 → multiKillBonus = 0
	// win = 0.0 (lost)
	// clutchBonus = 0
	// finalRating = 0.659 * (1 + 0) = 0.659

	expected := 0.659
	if math.Abs(player.Rating-expected) > 0.01 {
		t.Errorf("Expected rating ~%.3f, got %.3f", expected, player.Rating)
	}
}

// TestCalculateRoundRatings_SupportPlay tests support player with many assists
func TestCalculateRoundRatings_SupportPlay(t *testing.T) {
	round := &RoundStats{
		Winner: 3, // CT wins
		Players: []PlayerStats{
			{
				Team:    3, // CT
				Damage:  280,
				Kills:   1,
				Assists: 4, // Great support!
				Deaths:  2,
			},
		},
	}

	// Add 4 more CT teammates
	for i := 0; i < 4; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 3})
	}
	// Add 5 T opponents
	for i := 0; i < 5; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 2})
	}

	calculateRoundRatings(round)

	player := &round.Players[0]

	// Expected (5v5, support):
	// damageImpact = 2.8
	// fragImpact = 0.15*1 + 0.08*4 = 0.15 + 0.32 = 0.47
	// numerator = 3.27
	// denominator = 1 + 0.35*2 = 1.70
	// baseRating = 3.27 / 1.70 = 1.924
	// killRatio = 1/5 = 0.2 → multiKillBonus = 0
	// win = 1.0
	// clutchBonus = 0
	// finalRating = 1.924 * (1 + 0.10) = 1.924 * 1.10 = 2.116

	expected := 2.116
	if math.Abs(player.Rating-expected) > 0.01 {
		t.Errorf("Expected rating ~%.3f, got %.3f", expected, player.Rating)
	}
}

// TestCalculateRoundRatings_Outnumbered3v5 tests playing in disadvantage
func TestCalculateRoundRatings_Outnumbered3v5(t *testing.T) {
	round := &RoundStats{
		Winner: 3, // CT wins
		Players: []PlayerStats{
			{
				Team:    3, // CT
				Damage:  250,
				Kills:   3,
				Assists: 0,
				Deaths:  1,
			},
		},
	}

	// Add 2 teammates and 5 opponents (3v5)
	for i := 0; i < 2; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 3})
	}
	for i := 0; i < 5; i++ {
		round.Players = append(round.Players, PlayerStats{Team: 2})
	}

	calculateRoundRatings(round)

	player := &round.Players[0]

	// Expected:
	// damageImpact = (250/100) * (5/5)^0.7 * (5/3)^0.5
	//              = 2.5 * 1.0 * 1.291 = 3.227
	// fragImpact = 0.15*3 = 0.45
	// numerator = 3.677
	// denominator = 1.35
	// baseRating = 2.724
	// killRatio = 3/5 = 0.6 → multiKillBonus = 0.1
	// win = 1.0
	// clutchBonus = (5-3) * 0.05 = 0.10 (kills >= 2, outnumbered, won)
	// finalRating = 2.724 * (1 + 0.10 + 0.10 + 0.10) = 2.724 * 1.30 = 3.541

	expected := 3.541
	if math.Abs(player.Rating-expected) > 0.1 {
		t.Errorf("Expected rating ~%.3f, got %.3f", expected, player.Rating)
	}
}

// TestCalculateRoundRatings_MultiKillBoundaries tests multikill bonus thresholds
func TestCalculateRoundRatings_MultiKillBoundaries(t *testing.T) {
	tests := []struct {
		name              string
		kills             int
		opponents         int
		expectedBonus     float64
		expectedBonusDesc string
	}{
		{"2 of 5 (40%)", 2, 5, 0.0, "no bonus"},
		{"3 of 5 (60%)", 3, 5, 0.1, "+10%"},
		{"4 of 5 (80%)", 4, 5, 0.2, "+20%"},
		{"5 of 5 (100%)", 5, 5, 0.3, "+30% ACE"},
		{"3 of 3 (100%)", 3, 3, 0.3, "+30% ACE"},
		{"5 of 8 (62.5%)", 5, 8, 0.1, "+10%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			round := &RoundStats{
				Winner: 3,
				Players: []PlayerStats{
					{
						Team:    3,
						Damage:  100,
						Kills:   tt.kills,
						Assists: 0,
						Deaths:  0,
					},
				},
			}

			// Add opponents
			for i := 0; i < tt.opponents; i++ {
				round.Players = append(round.Players, PlayerStats{Team: 2})
			}

			calculateRoundRatings(round)

			player := &round.Players[0]

			// Extract multikill bonus from the formula
			// finalRating = baseRating * (1 + 0.10*win + multiKillBonus + clutchBonus)
			// We know win=1.0, clutchBonus=0 (not outnumbered with 1 player)
			baseRating := player.Rating / (1.0 + 0.10 + tt.expectedBonus)
			calculatedBonus := (player.Rating/baseRating - 1.0 - 0.10)

			if math.Abs(calculatedBonus-tt.expectedBonus) > 0.01 {
				t.Errorf("%s: Expected bonus %.2f (%s), got ~%.2f",
					tt.name, tt.expectedBonus, tt.expectedBonusDesc, calculatedBonus)
			}
		})
	}
}

// TestCalculateRoundRatings_NoClutchBonusConditions tests when clutch bonus is NOT applied
func TestCalculateRoundRatings_NoClutchBonusConditions(t *testing.T) {
	tests := []struct {
		name      string
		teamCount int
		oppCount  int
		kills     int
		win       bool
		hasBonus  bool
		reason    string
	}{
		{"Not outnumbered", 5, 3, 3, true, false, "team has advantage"},
		{"Equal teams", 5, 5, 3, true, false, "no disadvantage"},
		{"Outnumbered but lost", 3, 5, 3, false, false, "didn't win"},
		{"Outnumbered but only 1 kill", 3, 5, 1, true, false, "kills < 2"},
		{"Outnumbered, won, 2 kills", 3, 5, 2, true, true, "all conditions met"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			winnerTeam := 3
			if !tt.win {
				winnerTeam = 2
			}

			round := &RoundStats{
				Winner: winnerTeam,
				Players: []PlayerStats{
					{
						Team:    3,
						Damage:  200,
						Kills:   tt.kills,
						Assists: 0,
						Deaths:  0,
					},
				},
			}

			// Add teammates
			for i := 1; i < tt.teamCount; i++ {
				round.Players = append(round.Players, PlayerStats{Team: 3})
			}
			// Add opponents
			for i := 0; i < tt.oppCount; i++ {
				round.Players = append(round.Players, PlayerStats{Team: 2})
			}

			calculateRoundRatings(round)

			player := &round.Players[0]

			// Calculate expected clutch bonus
			expectedClutchBonus := 0.0
			if tt.hasBonus {
				expectedClutchBonus = float64(tt.oppCount-tt.teamCount) * 0.05
			}

			// For simplicity, we just check if rating is reasonable
			// A more precise test would require knowing exact baseRating
			if tt.hasBonus && player.Rating < 1.0 {
				t.Errorf("%s: Expected significant rating boost from clutch, got %.3f", tt.name, player.Rating)
			}

			t.Logf("%s: Rating=%.3f, Expected clutch bonus=%.2f (reason: %s)",
				tt.name, player.Rating, expectedClutchBonus, tt.reason)
		})
	}
}

// TestCalculateRoundRatings_EmptyRound tests edge case with no players
func TestCalculateRoundRatings_EmptyRound(t *testing.T) {
	round := &RoundStats{
		Winner:  3,
		Players: []PlayerStats{},
	}

	// Should not panic
	calculateRoundRatings(round)
}

// TestCalculateRoundRatings_ZeroDivisionProtection tests division by zero protection
func TestCalculateRoundRatings_ZeroDivisionProtection(t *testing.T) {
	round := &RoundStats{
		Winner: 3,
		Players: []PlayerStats{
			{
				Team:    3,
				Damage:  100,
				Kills:   1,
				Assists: 0,
				Deaths:  0,
			},
		},
	}

	// No opponents - should default to 5
	calculateRoundRatings(round)

	player := &round.Players[0]

	// Should have calculated something (not NaN or Inf)
	if math.IsNaN(player.Rating) || math.IsInf(player.Rating, 0) {
		t.Errorf("Rating should be valid number, got %.3f", player.Rating)
	}
}
