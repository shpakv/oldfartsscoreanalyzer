package components

import (
	"encoding/json"
	"fmt"
	"sort"

	"oldfartscounter/internal/stats"
)

// ProgressTabComponent отвечает за таб "Прогресс игроков"
type ProgressTabComponent struct{}

// NewProgressTab создает новый компонент таба прогресса
func NewProgressTab() *ProgressTabComponent {
	return &ProgressTabComponent{}
}

// PlayerProgress представляет данные прогресса игрока по датам
type PlayerProgress struct {
	AccountID int64              `json:"account_id"`
	Name      string             `json:"name"`
	Daily     []DailyPlayerStats `json:"daily"`
	Totals    PlayerTotalStats   `json:"totals"`
}

// DailyPlayerStats статистика игрока за один день
type DailyPlayerStats struct {
	Date         string  `json:"date"`
	RoundsPlayed int     `json:"rounds_played"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	Assists      int     `json:"assists"`
	Damage       int     `json:"damage"`
	WinRounds    int     `json:"win_rounds"`
	EPI          float64 `json:"epi"`
	KD           float64 `json:"kd"`
	ADR          float64 `json:"adr"`
	WinRate      float64 `json:"win_rate"`
}

// PlayerTotalStats общая статистика игрока
type PlayerTotalStats struct {
	RoundsPlayed int     `json:"rounds_played"`
	TotalKills   int     `json:"total_kills"`
	TotalDeaths  int     `json:"total_deaths"`
	AvgEPI       float64 `json:"avg_epi"`
	AvgKD        float64 `json:"avg_kd"`
	AvgADR       float64 `json:"avg_adr"`
	WinRate      float64 `json:"win_rate"`
}

// GenerateHTML генерирует HTML для таба прогресса
func (p *ProgressTabComponent) GenerateHTML() string {
	return `
<!-- PLAYER PROGRESS -->
<div id="tab-progress" class="view">
  <div class="toolbar">
    <label style="display:flex;align-items:center;gap:8px">
      <span style="color:var(--text);font-weight:600;">Выберите игрока:</span>
      <select id="playerSelect" style="background:var(--panel);color:var(--text);border:1px solid rgba(124,92,255,0.3);border-radius:6px;padding:8px 12px;font-size:14px;cursor:pointer;min-width:200px;">
        <option value="">-- Выберите игрока --</option>
      </select>
    </label>
    <span class="small" style="margin-left:auto;" id="progressInfo">Прогресс игрока</span>
  </div>

  <div id="playerProgressContent" style="display:none;padding:20px;">
    <!-- Карточка с общей статистикой -->
    <div style="background:linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);padding:24px;border-radius:12px;margin-bottom:24px;border:1px solid rgba(124,92,255,0.2);">
      <h3 style="margin:0 0 16px;color:#7c5cff;font-size:18px;">📊 Общая статистика</h3>
      <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(150px,1fr));gap:16px;">
        <div style="text-align:center;">
          <div style="font-size:28px;font-weight:bold;color:#22c55e;" id="stat-rounds">0</div>
          <div style="font-size:12px;color:var(--muted);margin-top:4px;">Раундов сыграно</div>
        </div>
        <div style="text-align:center;">
          <div style="font-size:28px;font-weight:bold;color:#3b82f6;" id="stat-epi">0.00</div>
          <div style="font-size:12px;color:var(--muted);margin-top:4px;">Средний EPI</div>
        </div>
        <div style="text-align:center;">
          <div style="font-size:28px;font-weight:bold;color:#f59e0b;" id="stat-kd">0.00</div>
          <div style="font-size:12px;color:var(--muted);margin-top:4px;">K/D Ratio</div>
        </div>
        <div style="text-align:center;">
          <div style="font-size:28px;font-weight:bold;color:#ef4444;" id="stat-adr">0</div>
          <div style="font-size:12px;color:var(--muted);margin-top:4px;">ADR</div>
        </div>
        <div style="text-align:center;">
          <div style="font-size:28px;font-weight:bold;color:#8b5cf6;" id="stat-winrate">0%</div>
          <div style="font-size:12px;color:var(--muted);margin-top:4px;">Win Rate</div>
        </div>
      </div>
    </div>

    <!-- Графики -->
    <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(500px,1fr));gap:20px;">
      <!-- EPI Progress -->
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">📈 EPI по датам</h4>
        <canvas id="chartEPI" style="max-height:300px;"></canvas>
      </div>

      <!-- K/D Progress -->
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">⚔️ K/D Ratio по датам</h4>
        <canvas id="chartKD" style="max-height:300px;"></canvas>
      </div>

      <!-- Damage Progress -->
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">💥 Урон по датам</h4>
        <canvas id="chartDamage" style="max-height:300px;"></canvas>
      </div>

      <!-- Win Rate Progress -->
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">🏆 Win Rate по датам</h4>
        <canvas id="chartWinRate" style="max-height:300px;"></canvas>
      </div>

      <!-- Activity -->
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">📅 Активность (раунды)</h4>
        <canvas id="chartActivity" style="max-height:300px;"></canvas>
      </div>

      <!-- Kills/Deaths/Assists -->
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">🎯 Kills/Deaths/Assists</h4>
        <canvas id="chartKDA" style="max-height:300px;"></canvas>
      </div>
    </div>
  </div>

  <div id="noPlayerSelected" style="text-align:center;padding:60px 20px;color:var(--muted);">
    <div style="font-size:64px;margin-bottom:16px;">📊</div>
    <div style="font-size:18px;margin-bottom:8px;">Выберите игрока</div>
    <div style="font-size:14px;">Выберите игрока из списка выше, чтобы увидеть его прогресс</div>
  </div>
</div>`
}

// GenerateJS генерирует JavaScript для таба прогресса
func (p *ProgressTabComponent) GenerateJS(data *stats.StatsData) string {
	// Собираем данные прогресса для каждого игрока
	progressData := p.buildProgressData(data)

	// Сериализуем в JSON
	jsonData, err := json.Marshal(progressData)
	if err != nil {
		return fmt.Sprintf(`console.error('Failed to generate progress data: %s');`, err.Error())
	}

	return fmt.Sprintf(`
// Init: Player Progress
(function() {
  const progressData = %s;
  const playerSelect = document.getElementById('playerSelect');
  const contentDiv = document.getElementById('playerProgressContent');
  const noPlayerDiv = document.getElementById('noPlayerSelected');

  let charts = {};

  // Заполняем список игроков
  progressData.forEach(player => {
    const option = document.createElement('option');
    option.value = player.account_id;
    option.textContent = player.name;
    playerSelect.appendChild(option);
  });

  // Обработчик выбора игрока
  playerSelect.addEventListener('change', function() {
    const accountId = parseInt(this.value);
    if (!accountId) {
      contentDiv.style.display = 'none';
      noPlayerDiv.style.display = 'block';
      return;
    }

    const player = progressData.find(p => p.account_id === accountId);
    if (!player) return;

    noPlayerDiv.style.display = 'none';
    contentDiv.style.display = 'block';

    renderPlayerProgress(player);
  });

  function renderPlayerProgress(player) {
    // Обновляем общую статистику
    document.getElementById('stat-rounds').textContent = player.totals.rounds_played;
    document.getElementById('stat-epi').textContent = player.totals.avg_epi.toFixed(2);
    document.getElementById('stat-kd').textContent = player.totals.avg_kd.toFixed(2);
    document.getElementById('stat-adr').textContent = Math.round(player.totals.avg_adr);
    document.getElementById('stat-winrate').textContent = player.totals.win_rate.toFixed(1) + '%%';

    // Уничтожаем старые графики
    Object.values(charts).forEach(chart => chart.destroy());
    charts = {};

    const dates = player.daily.map(d => d.date);
    const chartOptions = {
      responsive: true,
      maintainAspectRatio: true,
      plugins: {
        legend: { display: false },
        tooltip: {
          backgroundColor: 'rgba(26, 26, 30, 0.95)',
          titleColor: '#7c5cff',
          bodyColor: '#e5e5e5',
          borderColor: 'rgba(124, 92, 255, 0.3)',
          borderWidth: 1,
          padding: 12,
          displayColors: false
        }
      },
      scales: {
        x: {
          grid: { color: 'rgba(124, 92, 255, 0.1)' },
          ticks: { color: '#888' }
        },
        y: {
          grid: { color: 'rgba(124, 92, 255, 0.1)' },
          ticks: { color: '#888' },
          beginAtZero: true
        }
      }
    };

    // График EPI
    charts.epi = new Chart(document.getElementById('chartEPI'), {
      type: 'line',
      data: {
        labels: dates,
        datasets: [{
          label: 'EPI',
          data: player.daily.map(d => d.epi),
          borderColor: '#3b82f6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderWidth: 2,
          tension: 0.3,
          fill: true
        }]
      },
      options: chartOptions
    });

    // График K/D
    charts.kd = new Chart(document.getElementById('chartKD'), {
      type: 'line',
      data: {
        labels: dates,
        datasets: [{
          label: 'K/D',
          data: player.daily.map(d => d.kd),
          borderColor: '#f59e0b',
          backgroundColor: 'rgba(245, 158, 11, 0.1)',
          borderWidth: 2,
          tension: 0.3,
          fill: true
        }]
      },
      options: chartOptions
    });

    // График урона
    charts.damage = new Chart(document.getElementById('chartDamage'), {
      type: 'bar',
      data: {
        labels: dates,
        datasets: [{
          label: 'ADR',
          data: player.daily.map(d => d.adr),
          backgroundColor: 'rgba(239, 68, 68, 0.7)',
          borderColor: '#ef4444',
          borderWidth: 1
        }]
      },
      options: chartOptions
    });

    // График Win Rate
    charts.winrate = new Chart(document.getElementById('chartWinRate'), {
      type: 'line',
      data: {
        labels: dates,
        datasets: [{
          label: 'Win Rate',
          data: player.daily.map(d => d.win_rate),
          borderColor: '#8b5cf6',
          backgroundColor: 'rgba(139, 92, 246, 0.2)',
          borderWidth: 2,
          tension: 0.3,
          fill: true
        }]
      },
      options: {
        ...chartOptions,
        scales: {
          ...chartOptions.scales,
          y: {
            ...chartOptions.scales.y,
            max: 100
          }
        }
      }
    });

    // График активности
    charts.activity = new Chart(document.getElementById('chartActivity'), {
      type: 'bar',
      data: {
        labels: dates,
        datasets: [{
          label: 'Раундов',
          data: player.daily.map(d => d.rounds_played),
          backgroundColor: 'rgba(34, 197, 94, 0.7)',
          borderColor: '#22c55e',
          borderWidth: 1
        }]
      },
      options: chartOptions
    });

    // График K/D/A
    charts.kda = new Chart(document.getElementById('chartKDA'), {
      type: 'line',
      data: {
        labels: dates,
        datasets: [
          {
            label: 'Kills',
            data: player.daily.map(d => d.kills),
            borderColor: '#22c55e',
            backgroundColor: 'rgba(34, 197, 94, 0.1)',
            borderWidth: 2,
            tension: 0.3
          },
          {
            label: 'Deaths',
            data: player.daily.map(d => d.deaths),
            borderColor: '#ef4444',
            backgroundColor: 'rgba(239, 68, 68, 0.1)',
            borderWidth: 2,
            tension: 0.3
          },
          {
            label: 'Assists',
            data: player.daily.map(d => d.assists),
            borderColor: '#3b82f6',
            backgroundColor: 'rgba(59, 130, 246, 0.1)',
            borderWidth: 2,
            tension: 0.3
          }
        ]
      },
      options: {
        ...chartOptions,
        plugins: {
          ...chartOptions.plugins,
          legend: { display: true, labels: { color: '#e5e5e5' } }
        }
      }
    });
  }
})();
`, string(jsonData))
}

// buildProgressData собирает данные прогресса для всех игроков
func (p *ProgressTabComponent) buildProgressData(data *stats.StatsData) []PlayerProgress {
	// Создаем карту игроков
	playerMap := make(map[int64]*PlayerProgress)

	// Создаем карту дневной статистики: accountID -> date -> stats
	dailyMap := make(map[int64]map[string]*DailyPlayerStats)

	// Обрабатываем раунды по датам
	for date, rounds := range data.DailyRounds {
		for _, round := range rounds {
			// Обрабатываем каждого игрока в раунде
			for _, playerStat := range round.Players {
				accountID := playerStat.AccountID

				// Инициализируем карты если нужно
				if dailyMap[accountID] == nil {
					dailyMap[accountID] = make(map[string]*DailyPlayerStats)
				}
				if dailyMap[accountID][date] == nil {
					dailyMap[accountID][date] = &DailyPlayerStats{
						Date: date,
					}
				}

				daily := dailyMap[accountID][date]
				daily.RoundsPlayed++
				daily.Kills += playerStat.Kills
				daily.Deaths += playerStat.Deaths
				daily.Assists += playerStat.Assists
				daily.Damage += playerStat.Damage

				// Считаем победы
				if round.Winner == playerStat.Team {
					daily.WinRounds++
				}
			}
		}
	}

	// Вычисляем агрегированные метрики и создаем PlayerProgress
	for accountID, dateMap := range dailyMap {
		// Получаем имя игрока из рейтингов
		var playerName string
		for _, rating := range data.PlayerRatings {
			if rating.AccountID == accountID {
				playerName = rating.Name
				break
			}
		}
		if playerName == "" {
			continue // Пропускаем игроков без имени
		}

		// Сортируем даты
		var dates []string
		for date := range dateMap {
			dates = append(dates, date)
		}
		sort.Strings(dates)

		// Создаем массив дневной статистики
		var dailyStats []DailyPlayerStats
		totalRounds := 0
		totalKills := 0
		totalDeaths := 0
		totalDamage := 0
		totalWins := 0

		for _, date := range dates {
			daily := dateMap[date]

			// Вычисляем средние значения
			if daily.RoundsPlayed > 0 {
				daily.ADR = float64(daily.Damage) / float64(daily.RoundsPlayed)
				daily.WinRate = (float64(daily.WinRounds) / float64(daily.RoundsPlayed)) * 100

				if daily.Deaths > 0 {
					daily.KD = float64(daily.Kills) / float64(daily.Deaths)
				} else if daily.Kills > 0 {
					daily.KD = float64(daily.Kills)
				}

				// Упрощенный расчет EPI на основе метрик (приблизительный)
				// Формула приблизительная: (Kills * 1.5 - Deaths * 1.0 + Assists * 0.5 + Damage/100) / Rounds
				daily.EPI = (float64(daily.Kills)*1.5 - float64(daily.Deaths) + float64(daily.Assists)*0.5 + float64(daily.Damage)/100.0) / float64(daily.RoundsPlayed)
			}

			dailyStats = append(dailyStats, *daily)

			// Собираем общую статистику
			totalRounds += daily.RoundsPlayed
			totalKills += daily.Kills
			totalDeaths += daily.Deaths
			totalDamage += daily.Damage
			totalWins += daily.WinRounds
		}

		// Вычисляем общую статистику
		avgEPI := 0.0
		avgKD := 0.0
		avgADR := 0.0
		winRate := 0.0

		if totalRounds > 0 {
			// EPI из агрегированных данных
			for _, rating := range data.PlayerRatings {
				if rating.AccountID == accountID {
					avgEPI = rating.AverageEPI
					break
				}
			}

			if totalDeaths > 0 {
				avgKD = float64(totalKills) / float64(totalDeaths)
			} else if totalKills > 0 {
				avgKD = float64(totalKills)
			}

			avgADR = float64(totalDamage) / float64(totalRounds)
			winRate = (float64(totalWins) / float64(totalRounds)) * 100
		}

		playerMap[accountID] = &PlayerProgress{
			AccountID: accountID,
			Name:      playerName,
			Daily:     dailyStats,
			Totals: PlayerTotalStats{
				RoundsPlayed: totalRounds,
				TotalKills:   totalKills,
				TotalDeaths:  totalDeaths,
				AvgEPI:       avgEPI,
				AvgKD:        avgKD,
				AvgADR:       avgADR,
				WinRate:      winRate,
			},
		}
	}

	// Преобразуем в массив и сортируем по имени
	var result []PlayerProgress
	for _, progress := range playerMap {
		result = append(result, *progress)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}
