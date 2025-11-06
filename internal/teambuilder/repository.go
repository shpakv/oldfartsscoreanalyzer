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
	GetAverageMu() float64 // Средний EPI (μ) для расчета категорий рейтинга
}

// singleton реализация
type playerRepository struct {
	data      []Player
	averageMu float64 // Средний EPI (μ) из логов для расчета категорий
}

var (
	instance *playerRepository
	once     sync.Once
)

func NewPlayerRepository() PlayerRepository {
	once.Do(func() {
		instance = &playerRepository{
			averageMu: 3.027, // Средний EPI (μ) из HTML (обновлено 2025-10-31)
			data: []Player{
				{NickName: "povidlo boy", Score: 4.561},
				{NickName: "maslina420", Score: 4.260},
				{NickName: "Looka", Score: 3.758},
				{NickName: "C.C.Capwell", Score: 3.348},
				{NickName: "Astracore", Score: 3.297},
				{NickName: "jojo", Score: 3.194},
				{NickName: "Pyatka", Score: 3.192},
				{NickName: "whereispie", Score: 3.002},
				{NickName: "d3msk", Score: 2.971},
				{NickName: "Mr. Titspervert", Score: 2.917},
				{NickName: "Fitz [BadCom]", Score: 2.838},
				{NickName: "Rezec", Score: 2.812},
				{NickName: "ℭŗυşş", Score: 2.772},
				{NickName: "Крыса Сплинтер", Score: 2.589},
				{NickName: "Баба Валя", Score: 2.515},
				{NickName: "cyberhawk2000n", Score: 2.481},
				{NickName: "petya_vpered", Score: 2.407},
				{NickName: "Boberto", Score: 2.391},
				{NickName: "Gharb", Score: 2.382},
				{NickName: "T1TAN", Score: 2.221},
				{NickName: "Chu [BadCom]", Score: 1.959},
				{NickName: "Djafar-AGA", Score: 1.866},
				{NickName: "atlas", Score: 1.632},

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

// GetAverageMu — возвращает средний EPI (μ) для расчета категорий рейтинга
func (r *playerRepository) GetAverageMu() float64 {
	return r.averageMu
}
