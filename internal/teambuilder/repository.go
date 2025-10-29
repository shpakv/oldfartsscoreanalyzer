package teambuilder

import "sync"

// Player — модель игрока
type Player struct {
	NickName string  `json:"nickName"`
	Score    float64 `json:"score"`
}

// PlayerRepository — интерфейс репозитория
type PlayerRepository interface {
	GetAll() []Player
	GetTop(n int) []Player
	FindByName(nick string) *Player
}

// singleton реализация
type playerRepository struct {
	data []Player
}

var (
	instance *playerRepository
	once     sync.Once
)

func NewPlayerRepository() PlayerRepository {
	once.Do(func() {
		instance = &playerRepository{
			data: []Player{
				{NickName: "povidlo boy", Score: 3.825},
				{NickName: "maslina420", Score: 3.338},
				{NickName: "Looka", Score: 2.904},
				{NickName: "Astracore", Score: 2.648},
				{NickName: "C.C.Capwell", Score: 2.641},
				{NickName: "jojo", Score: 2.592},
				{NickName: "d3msk", Score: 2.501},
				{NickName: "whereispie", Score: 2.461},
				{NickName: "Rezec", Score: 2.417},
				{NickName: "Mr. Titspervert", Score: 2.403},
				{NickName: "Pyatka", Score: 2.377},
				{NickName: "Fitz [BadCom]", Score: 2.309},
				{NickName: "çruşş", Score: 2.253},
				{NickName: "T1TAN", Score: 2.144},
				{NickName: "Chu [BadCom]", Score: 2.112},
				{NickName: "Крыса Сплинтер", Score: 2.108},
				{NickName: "Boberto", Score: 2.081},
				{NickName: "Баба Валя", Score: 2.079},
				{NickName: "cyberhawk2000n", Score: 2.016},
				{NickName: "petya_vpered", Score: 1.983},
				{NickName: "Gharb", Score: 1.953},
				{NickName: "atlas", Score: 1.598},
				{NickName: "Djafar-AGA", Score: 1.565},
				{NickName: "BingoBongo", Score: 1.150},
				{NickName: "Mr. Your mom", Score: 1.150},
				{NickName: "Mirshdd", Score: 1.150},
				{NickName: "Extain", Score: 1.0},
				{NickName: "Lexus", Score: 1.0},
				{NickName: "Station77", Score: 0.7},
			},
		}
	})
	return instance
}

// GetAll — возвращает всех игроков
func (r *playerRepository) GetAll() []Player {
	return r.data
}

// GetTop — возвращает n лучших игроков
func (r *playerRepository) GetTop(n int) []Player {
	if n > len(r.data) {
		n = len(r.data)
	}
	return r.data[:n]
}

// FindByName — поиск игрока по нику
func (r *playerRepository) FindByName(nick string) *Player {
	for _, p := range r.data {
		if p.NickName == nick {
			return &p
		}
	}
	return nil
}
