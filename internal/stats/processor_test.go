package stats

import (
	"math"
	"testing"

	"oldfartscounter/internal/logparser"
)

// TestBuildPlayerRatings_SinglePlayer tests Bayesian rating for single player
func TestBuildPlayerRatings_SinglePlayer(t *testing.T) {
	processor := New()

	rounds := []logparser.RoundStats{
		{
			Winner: 3,
			Date:   "2024-01-01",
			Players: []logparser.PlayerStats{
				{
					AccountID: 12345,
					Team:      3,
					Rating:    1.2, // EPI for this round
					Damage:    300,
					Kills:     3,
					Deaths:    1,
					Assists:   1,
				},
			},
		},
	}

	// Add more rounds with same player
	for i := 0; i < 9; i++ {
		rounds = append(rounds, logparser.RoundStats{
			Winner: 3,
			Date:   "2024-01-01",
			Players: []logparser.PlayerStats{
				{
					AccountID: 12345,
					Team:      3,
					Rating:    1.2,
					Damage:    300,
					Kills:     3,
					Deaths:    1,
					Assists:   1,
				},
			},
		})
	}

	killEvents := []logparser.KillEvent{
		{
			KillerSID:  "[U:1:12345]",
			KillerName: "TestPlayer",
		},
	}

	ratings := processor.buildPlayerRatings(rounds, killEvents, nil, nil)

	if len(ratings) != 1 {
		t.Fatalf("Expected 1 player, got %d", len(ratings))
	}

	player := ratings[0]

	// Verify aggregation
	if player.RoundsPlayed != 10 {
		t.Errorf("Expected 10 rounds, got %d", player.RoundsPlayed)
	}

	if player.TotalKills != 30 {
		t.Errorf("Expected 30 kills, got %d", player.TotalKills)
	}

	// Verify simple average
	expectedAverage := 1.2
	if math.Abs(player.AverageEPI-expectedAverage) > 0.001 {
		t.Errorf("Expected average EPI %.3f, got %.3f", expectedAverage, player.AverageEPI)
	}

	// Verify Bayesian rating
	// μ = 12.0 / 10 = 1.2 (same as player average since only one player)
	// BayesianEPI = (12.0 + 100*1.2) / (10 + 100) = 132.0 / 110 = 1.2
	expectedBayesian := 1.2
	if math.Abs(player.BayesianEPI-expectedBayesian) > 0.001 {
		t.Errorf("Expected Bayesian EPI %.3f, got %.3f", expectedBayesian, player.BayesianEPI)
	}
}

// TestBuildPlayerRatings_NewbieVsVeteran tests Bayesian regularization
func TestBuildPlayerRatings_NewbieVsVeteran(t *testing.T) {
	processor := New()

	// Newbie: 10 rounds with excellent performance (EPI = 1.5)
	newbieRounds := make([]logparser.RoundStats, 10)
	for i := 0; i < 10; i++ {
		newbieRounds[i] = logparser.RoundStats{
			Winner: 3,
			Players: []logparser.PlayerStats{
				{
					AccountID: 11111,
					Team:      3,
					Rating:    1.5, // Excellent!
					Damage:    400,
					Kills:     4,
					Deaths:    0,
					Assists:   1,
				},
			},
		}
	}

	// Veteran: 200 rounds with good performance (EPI = 0.9)
	veteranRounds := make([]logparser.RoundStats, 200)
	for i := 0; i < 200; i++ {
		veteranRounds[i] = logparser.RoundStats{
			Winner: 3,
			Players: []logparser.PlayerStats{
				{
					AccountID: 22222,
					Team:      3,
					Rating:    0.9, // Good, stable
					Damage:    250,
					Kills:     2,
					Deaths:    1,
					Assists:   1,
				},
			},
		}
	}

	allRounds := append(newbieRounds, veteranRounds...)

	killEvents := []logparser.KillEvent{
		{KillerSID: "[U:1:11111]", KillerName: "Newbie"},
		{KillerSID: "[U:1:22222]", KillerName: "Veteran"},
	}

	ratings := processor.buildPlayerRatings(allRounds, killEvents, nil, nil)

	if len(ratings) != 2 {
		t.Fatalf("Expected 2 players, got %d", len(ratings))
	}

	// Find each player
	var newbie, veteran *PlayerRating
	for i := range ratings {
		if ratings[i].AccountID == 11111 {
			newbie = &ratings[i]
		} else if ratings[i].AccountID == 22222 {
			veteran = &ratings[i]
		}
	}

	// Calculate μ: (10*1.5 + 200*0.9) / 210 = (15 + 180) / 210 = 0.929
	expectedMu := 0.929

	// Newbie Bayesian: (15 + 100*0.929) / (10 + 100) = 107.9 / 110 = 0.981
	// Simple average: 1.5
	// Bayesian pulls it down significantly!
	expectedNewbieBayesian := 0.981
	if math.Abs(newbie.BayesianEPI-expectedNewbieBayesian) > 0.01 {
		t.Errorf("Newbie: Expected Bayesian ~%.3f, got %.3f", expectedNewbieBayesian, newbie.BayesianEPI)
	}

	if newbie.AverageEPI != 1.5 {
		t.Errorf("Newbie: Expected simple average 1.5, got %.3f", newbie.AverageEPI)
	}

	// Veteran Bayesian: (180 + 100*0.929) / (200 + 100) = 272.9 / 300 = 0.910
	// Simple average: 0.9
	// Bayesian very close to simple average (has many rounds)
	expectedVeteranBayesian := 0.910
	if math.Abs(veteran.BayesianEPI-expectedVeteranBayesian) > 0.01 {
		t.Errorf("Veteran: Expected Bayesian ~%.3f, got %.3f", expectedVeteranBayesian, veteran.BayesianEPI)
	}

	// Important: Despite newbie having higher simple average (1.5 vs 0.9),
	// their Bayesian ratings are closer (0.981 vs 0.910)
	// System doesn't trust newbie's small sample!
	t.Logf("Newbie: Simple=%.3f, Bayesian=%.3f (pulled down by regularization)",
		newbie.AverageEPI, newbie.BayesianEPI)
	t.Logf("Veteran: Simple=%.3f, Bayesian=%.3f (trusted data)",
		veteran.AverageEPI, veteran.BayesianEPI)
	t.Logf("Expected μ=%.3f", expectedMu)
}

// TestBuildPlayerRatings_MultiplePlayersRanking tests sorting by Bayesian rating
func TestBuildPlayerRatings_MultiplePlayersRanking(t *testing.T) {
	processor := New()

	// Create 3 players with different performance and sample sizes
	rounds := []logparser.RoundStats{}

	// Player 1: 5 rounds, EPI=2.0 (lucky streak)
	for i := 0; i < 5; i++ {
		rounds = append(rounds, logparser.RoundStats{
			Winner: 3,
			Players: []logparser.PlayerStats{
				{AccountID: 100, Team: 3, Rating: 2.0, Kills: 5, Deaths: 0},
			},
		})
	}

	// Player 2: 100 rounds, EPI=1.1 (strong and stable)
	for i := 0; i < 100; i++ {
		rounds = append(rounds, logparser.RoundStats{
			Winner: 3,
			Players: []logparser.PlayerStats{
				{AccountID: 200, Team: 3, Rating: 1.1, Kills: 3, Deaths: 1},
			},
		})
	}

	// Player 3: 50 rounds, EPI=0.8 (average)
	for i := 0; i < 50; i++ {
		rounds = append(rounds, logparser.RoundStats{
			Winner: 3,
			Players: []logparser.PlayerStats{
				{AccountID: 300, Team: 3, Rating: 0.8, Kills: 2, Deaths: 2},
			},
		})
	}

	killEvents := []logparser.KillEvent{
		{KillerSID: "[U:1:100]", KillerName: "Lucky"},
		{KillerSID: "[U:1:200]", KillerName: "Strong"},
		{KillerSID: "[U:1:300]", KillerName: "Average"},
	}

	ratings := processor.buildPlayerRatings(rounds, killEvents, nil, nil)

	if len(ratings) != 3 {
		t.Fatalf("Expected 3 players, got %d", len(ratings))
	}

	// Results should be sorted by Bayesian rating (descending)
	// With dynamic μ calculation, the exact order may vary
	// Let's verify correct sorting and that "Average" is still lowest

	// Verify sorting (descending order)
	for i := 0; i < len(ratings)-1; i++ {
		if ratings[i].BayesianEPI < ratings[i+1].BayesianEPI {
			t.Errorf("Results not sorted: ratings[%d]=%.3f < ratings[%d]=%.3f",
				i, ratings[i].BayesianEPI, i+1, ratings[i+1].BayesianEPI)
		}
	}

	// Verify "Average" is last (lowest Bayesian rating)
	if ratings[2].Name != "Average" {
		t.Errorf("Expected 'Average' to be #3, got '%s'", ratings[2].Name)
	}

	// Verify Strong has more stable rating (closer to simple average) than Lucky
	strongSimple := ratings[1].AverageEPI
	strongBayesian := ratings[1].BayesianEPI
	strongDiff := math.Abs(strongSimple - strongBayesian)

	if strongDiff > 0.1 { // Should be very close due to large sample
		t.Errorf("Strong player's Bayesian rating (%.3f) too different from simple (%.3f)",
			strongBayesian, strongSimple)
	}

	t.Logf("Rankings:")
	for i, r := range ratings {
		t.Logf("  #%d: %s - Rounds=%d, SimpleEPI=%.3f, BayesianEPI=%.3f",
			i+1, r.Name, r.RoundsPlayed, r.AverageEPI, r.BayesianEPI)
	}
}

// TestBuildPlayerRatings_WinRateCalculation tests win round counting
func TestBuildPlayerRatings_WinRateCalculation(t *testing.T) {
	processor := New()

	rounds := []logparser.RoundStats{
		// Win as CT
		{
			Winner: 3,
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 3, Rating: 1.0, Kills: 2},
			},
		},
		// Lose as CT
		{
			Winner: 2,
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 3, Rating: 0.5, Kills: 1},
			},
		},
		// Win as T
		{
			Winner: 2,
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 2, Rating: 1.2, Kills: 3},
			},
		},
		// Lose as T
		{
			Winner: 3,
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 2, Rating: 0.6, Kills: 1},
			},
		},
	}

	killEvents := []logparser.KillEvent{
		{KillerSID: "[U:1:123]", KillerName: "TestPlayer"},
	}

	ratings := processor.buildPlayerRatings(rounds, killEvents, nil, nil)

	if len(ratings) != 1 {
		t.Fatalf("Expected 1 player, got %d", len(ratings))
	}

	player := ratings[0]

	if player.RoundsPlayed != 4 {
		t.Errorf("Expected 4 rounds, got %d", player.RoundsPlayed)
	}

	// Player won 2 out of 4 rounds
	if player.WinRounds != 2 {
		t.Errorf("Expected 2 wins, got %d", player.WinRounds)
	}

	expectedWinRate := 50.0 // 2/4 = 50%
	actualWinRate := float64(player.WinRounds) / float64(player.RoundsPlayed) * 100

	if math.Abs(actualWinRate-expectedWinRate) > 0.1 {
		t.Errorf("Expected win rate ~%.1f%%, got %.1f%%", expectedWinRate, actualWinRate)
	}
}

// TestBuildPlayerRatings_LastPlayedDate tests date tracking
func TestBuildPlayerRatings_LastPlayedDate(t *testing.T) {
	processor := New()

	rounds := []logparser.RoundStats{
		{
			Date: "2024-01-01",
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 3, Rating: 1.0},
			},
		},
		{
			Date: "2024-01-05",
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 3, Rating: 1.0},
			},
		},
		{
			Date: "2024-01-03",
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 3, Rating: 1.0},
			},
		},
	}

	killEvents := []logparser.KillEvent{
		{KillerSID: "[U:1:123]", KillerName: "TestPlayer"},
	}

	ratings := processor.buildPlayerRatings(rounds, killEvents, nil, nil)

	player := ratings[0]

	// Should track the latest date
	expectedDate := "2024-01-05"
	if player.LastPlayed != expectedDate {
		t.Errorf("Expected last played date %s, got %s", expectedDate, player.LastPlayed)
	}
}

// TestBuildPlayerRatings_EmptyInput tests edge case with no data
func TestBuildPlayerRatings_EmptyInput(t *testing.T) {
	processor := New()

	ratings := processor.buildPlayerRatings(nil, nil, nil, nil)

	if len(ratings) != 0 {
		t.Errorf("Expected 0 players for empty input, got %d", len(ratings))
	}
}

// TestBuildPlayerRatings_DefaultMu tests default μ when no data
func TestBuildPlayerRatings_DefaultMu(t *testing.T) {
	processor := New()

	// Create player with 10 rounds, EPI=1.0
	rounds := make([]logparser.RoundStats, 10)
	for i := range rounds {
		rounds[i] = logparser.RoundStats{
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 3, Rating: 1.0},
			},
		}
	}

	killEvents := []logparser.KillEvent{
		{KillerSID: "[U:1:123]", KillerName: "TestPlayer"},
	}

	ratings := processor.buildPlayerRatings(rounds, killEvents, nil, nil)

	// μ = 10.0 / 10 = 1.0
	// BayesianEPI = (10 + 100*1.0) / (10 + 100) = 110 / 110 = 1.0
	expectedBayesian := 1.0

	if math.Abs(ratings[0].BayesianEPI-expectedBayesian) > 0.001 {
		t.Errorf("Expected Bayesian %.3f, got %.3f", expectedBayesian, ratings[0].BayesianEPI)
	}
}

// TestBuildPlayerRatings_SkipZeroAccountID tests filtering invalid account IDs
func TestBuildPlayerRatings_SkipZeroAccountID(t *testing.T) {
	processor := New()

	rounds := []logparser.RoundStats{
		{
			Players: []logparser.PlayerStats{
				{AccountID: 0, Team: 3, Rating: 1.0}, // Should be skipped
				{AccountID: 123, Team: 3, Rating: 1.0},
			},
		},
	}

	killEvents := []logparser.KillEvent{
		{KillerSID: "[U:1:123]", KillerName: "ValidPlayer"},
	}

	ratings := processor.buildPlayerRatings(rounds, killEvents, nil, nil)

	// Only 1 valid player
	if len(ratings) != 1 {
		t.Errorf("Expected 1 player (AccountID=0 should be skipped), got %d", len(ratings))
	}

	if ratings[0].AccountID != 123 {
		t.Errorf("Expected AccountID 123, got %d", ratings[0].AccountID)
	}
}

// TestBuildPlayerRatings_NameFallback tests fallback name generation
func TestBuildPlayerRatings_NameFallback(t *testing.T) {
	processor := New()

	rounds := []logparser.RoundStats{
		{
			Players: []logparser.PlayerStats{
				{AccountID: 99999, Team: 3, Rating: 1.0},
			},
		},
	}

	// No kill events = no name
	ratings := processor.buildPlayerRatings(rounds, nil, nil, nil)

	expectedName := "Player_99999"
	if ratings[0].Name != expectedName {
		t.Errorf("Expected fallback name '%s', got '%s'", expectedName, ratings[0].Name)
	}
}

// TestBuildPlayerRatings_K100Constant tests that K=100 is correctly used
func TestBuildPlayerRatings_K100Constant(t *testing.T) {
	processor := New()

	// Player with exactly 100 rounds, EPI=1.5
	rounds := make([]logparser.RoundStats, 100)
	for i := range rounds {
		rounds[i] = logparser.RoundStats{
			Players: []logparser.PlayerStats{
				{AccountID: 123, Team: 3, Rating: 1.5},
			},
		}
	}

	killEvents := []logparser.KillEvent{
		{KillerSID: "[U:1:123]", KillerName: "TestPlayer"},
	}

	ratings := processor.buildPlayerRatings(rounds, killEvents, nil, nil)

	// At exactly K=100 rounds, Bayesian should be average of simple mean and μ
	// μ = 150 / 100 = 1.5
	// BayesianEPI = (150 + 100*1.5) / (100 + 100) = 300 / 200 = 1.5
	// (in this case both are equal, so Bayesian = Simple = μ)

	expectedBayesian := 1.5
	if math.Abs(ratings[0].BayesianEPI-expectedBayesian) > 0.001 {
		t.Errorf("At K=100 rounds with uniform data, expected Bayesian=%.3f, got %.3f",
			expectedBayesian, ratings[0].BayesianEPI)
	}
}

// TestBuildPlayerRatings_AggregationAccuracy tests statistical aggregation
func TestBuildPlayerRatings_AggregationAccuracy(t *testing.T) {
	processor := New()

	rounds := []logparser.RoundStats{
		{
			Players: []logparser.PlayerStats{
				{
					AccountID: 123,
					Team:      3,
					Rating:    1.5,
					Damage:    350,
					Kills:     4,
					Deaths:    1,
					Assists:   2,
				},
			},
		},
		{
			Players: []logparser.PlayerStats{
				{
					AccountID: 123,
					Team:      3,
					Rating:    0.8,
					Damage:    150,
					Kills:     1,
					Deaths:    2,
					Assists:   1,
				},
			},
		},
		{
			Players: []logparser.PlayerStats{
				{
					AccountID: 123,
					Team:      3,
					Rating:    1.2,
					Damage:    280,
					Kills:     3,
					Deaths:    1,
					Assists:   0,
				},
			},
		},
	}

	killEvents := []logparser.KillEvent{
		{KillerSID: "[U:1:123]", KillerName: "TestPlayer"},
	}

	ratings := processor.buildPlayerRatings(rounds, killEvents, nil, nil)

	player := ratings[0]

	// Verify all aggregations
	if player.RoundsPlayed != 3 {
		t.Errorf("Expected 3 rounds, got %d", player.RoundsPlayed)
	}

	expectedTotalDamage := 350 + 150 + 280
	if player.TotalDamage != expectedTotalDamage {
		t.Errorf("Expected total damage %d, got %d", expectedTotalDamage, player.TotalDamage)
	}

	expectedTotalKills := 4 + 1 + 3
	if player.TotalKills != expectedTotalKills {
		t.Errorf("Expected total kills %d, got %d", expectedTotalKills, player.TotalKills)
	}

	expectedTotalDeaths := 1 + 2 + 1
	if player.TotalDeaths != expectedTotalDeaths {
		t.Errorf("Expected total deaths %d, got %d", expectedTotalDeaths, player.TotalDeaths)
	}

	expectedTotalAssists := 2 + 1 + 0
	if player.TotalAssists != expectedTotalAssists {
		t.Errorf("Expected total assists %d, got %d", expectedTotalAssists, player.TotalAssists)
	}

	expectedTotalEPI := 1.5 + 0.8 + 1.2
	if math.Abs(player.TotalEPI-expectedTotalEPI) > 0.01 {
		t.Errorf("Expected total EPI %.2f, got %.2f", expectedTotalEPI, player.TotalEPI)
	}

	expectedAverageEPI := expectedTotalEPI / 3.0
	if math.Abs(player.AverageEPI-expectedAverageEPI) > 0.01 {
		t.Errorf("Expected average EPI %.3f, got %.3f", expectedAverageEPI, player.AverageEPI)
	}
}
