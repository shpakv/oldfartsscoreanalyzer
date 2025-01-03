package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"oldfartscounter/internal/teambuilder"
	"os"
)

var score *bool

func main() {
	score = flag.Bool("s", false, "show player score")
	filePath := flag.String("f", "", "Path to the configuration.json file")

	flag.Parse()
	if *filePath == "" {
		log.Fatal("Please provide a file path using the -f flag")
	}

	if _, err := os.Stat(*filePath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", *filePath)
	}

	content, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var config teambuilder.TeamConfiguration
	if err = json.Unmarshal(content, &config); err != nil {
		log.Fatalf("Invalid JSON format: %v", err)
	}

	b := &teambuilder.TeamBuilder{}

	team1, team2 := b.Build(&config)

	totalWeight1 := teamWeight(team1, "Team 1")
	totalWeight2 := teamWeight(team2, "Team 2")

	fmt.Printf("*** Score Difference: %.2f ***\n", math.Abs(totalWeight1-totalWeight2))
}

func teamWeight(team teambuilder.Team, name string) float64 {
	fmt.Println(name)
	totalWeight := 0.0
	for _, player := range team {
		if *score {
			fmt.Printf("%s (%.2f)\n", player.Name, player.Score)
		} else {
			fmt.Printf("%s \n", player.Name)
		}
		totalWeight += player.Score
	}
	fmt.Printf("Total Score: %.2f\n\n", totalWeight)

	return totalWeight
}
