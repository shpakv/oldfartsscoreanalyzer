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
				{NickName: "povidlo boy", Score: 4.675},
				{NickName: "maslina420", Score: 4.276},
				{NickName: "Looka", Score: 3.811},
				{NickName: "Astracore", Score: 3.334},
				{NickName: "C.C.Capwell", Score: 3.309},
				{NickName: "jojo", Score: 3.238},
				{NickName: "Pyatka", Score: 3.148},
				{NickName: "d3msk", Score: 3.088},
				{NickName: "Rezec", Score: 3.026},
				{NickName: "whereispie", Score: 2.995},
				{NickName: "Mr. Titspervert", Score: 2.989},
				{NickName: "Fitz [BadCom]", Score: 2.876},
				{NickName: "çruşş", Score: 2.785},
				{NickName: "T1TAN", Score: 2.667},
				{NickName: "Крыса Сплинтер", Score: 2.631},
				{NickName: "Баба Валя", Score: 2.606},
				{NickName: "Chu [BadCom]", Score: 2.578},
				{NickName: "Boberto", Score: 2.570},
				{NickName: "cyberhawk2000n", Score: 2.508},
				{NickName: "petya_vpered", Score: 2.433},
				{NickName: "Gharb", Score: 2.402},
				{NickName: "Djafar-AGA", Score: 1.929},
				{NickName: "atlas", Score: 1.655},

				// Не играли по стате с сентября 2025.
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
