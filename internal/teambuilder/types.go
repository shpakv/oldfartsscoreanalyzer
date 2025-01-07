package teambuilder

// TeamPlayer представляет игрока с его индивидуальными характеристиками.
// Содержит имя и вес (или рейтинг) игрока.
type TeamPlayer struct {
	// NickName - уникальный идентификатор игрока.
	// Используется для применения ограничений и идентификации.
	NickName string `json:"nickName"`

	// Score - количественная характеристика игрока.
	// Может представлять собой рейтинг, навыки или другой числовой показатель.
	// Используется для балансировки команд.
	Score float64 `json:"score"`
}

// Team представляет собой коллекцию игроков.
// Является типом среза TeamPlayer с возможностью расширения функциональности.
type Team []TeamPlayer

// ConstrainType определяет тип ограничения при распределении игроков.
// Используется для задания правил формирования команд.
type ConstrainType string

const (
	// ConstraintTogether указывает, что определенные игроки.
	// ДОЛЖНЫ быть распределены в одну команду.
	// Применяется, когда важно сохранить взаимодействие между конкретными игроками.
	ConstraintTogether ConstrainType = "together"

	// ConstraintSeparate указывает, что определенные игроки.
	// ДОЛЖНЫ быть распределены в разные команды.
	// Используется для предотвращения конфликтов или создания баланса.
	ConstraintSeparate ConstrainType = "separate"
)

// Constraint описывает правило распределения игроков.
// Позволяет создавать сложные сценарии формирования команд.
type Constraint struct {
	// Type определяет тип ограничения.
	// Может принимать значения ConstraintTogether или ConstraintSeparate.
	Type ConstrainType `json:"type"`

	// Player1 и Player2 - имена игроков, к которым применяется ограничение.
	// Используются для идентификации конкретных участников правила.
	//
	// Примеры использования:
	//   - Держать сильных игроков вместе
	//   - Разделять игроков с конфликтной историей
	//   - Создавать сбалансированные составы
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
}

type Constraints []Constraint

type TeamConfiguration struct {
	Players     Team        `json:"players"`
	Constraints Constraints `json:"constraints"`
}

func (t Team) Score() float64 {
	score := 0.0
	for _, player := range t {
		score += player.Score
	}
	return score
}
