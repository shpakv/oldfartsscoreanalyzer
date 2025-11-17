package tournament

import (
	"encoding/json"
	"os"
)

// Config представляет конфигурацию турнира
type Config struct {
	Name          string   `json:"name"`
	Date          string   `json:"date"`
	DurationHours int      `json:"duration_hours"`
	StartTime     string   `json:"start_time"`
	Participants  []string `json:"participants"`
	Prizes        struct {
		FirstPlace  string `json:"1st_place"`
		SecondPlace string `json:"2nd_place"`
		ThirdPlace  string `json:"3rd_place"`
	} `json:"prizes"`
	Program struct {
		Stage1 Stage `json:"stage_1"`
		Break  Break `json:"break"`
		Stage2 Stage `json:"stage_2"`
	} `json:"program"`
	Format struct {
		Type           string `json:"type"`
		Description    string `json:"description"`
		MinMapsPerTeam int    `json:"min_maps_per_team"`
		MaxMapsPerTeam int    `json:"max_maps_per_team"`
	} `json:"format"`
	Maps       []string   `json:"maps"`
	Teams      TeamConfig `json:"teams"`
	Organizers struct {
		Main                  string `json:"main"`
		Servers               string `json:"servers"`
		MapSelectionModerator string `json:"map_selection_moderator"`
	} `json:"organizers"`
}

// Stage представляет этап турнира
type Stage struct {
	Name    string  `json:"name"`
	Format  string  `json:"format"`
	Matches []Match `json:"matches"`
	Note    string  `json:"note"`
}

// Match представляет матч
type Match struct {
	Match string   `json:"match"`
	Teams []string `json:"teams"`
}

// Break представляет перерыв
type Break struct {
	DurationMinutes int    `json:"duration_minutes"`
	Purpose         string `json:"purpose"`
}

// TeamConfig представляет конфигурацию команд
type TeamConfig struct {
	Count int      `json:"count"`
	Size  int      `json:"size"`
	Names []string `json:"names"`
}

// LoadConfig загружает конфигурацию турнира из JSON файла
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is controlled by application code
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// FormatDate форматирует дату в читаемый вид (например, "20 декабря 2025")
func (c *Config) FormatDate() string {
	// Можно добавить более сложную логику форматирования
	// Пока просто возвращаем как есть
	return c.Date
}
