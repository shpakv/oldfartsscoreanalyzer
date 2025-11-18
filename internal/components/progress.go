package components

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"oldfartscounter/internal/stats"
)

// ProgressTabComponent отвечает за таб "Прогресс игроков"
type ProgressTabComponent struct{}

// NewProgressTab создает новый компонент таба прогресса
func NewProgressTab() *ProgressTabComponent {
	return &ProgressTabComponent{}
}

// TimeSlotStats статистика по временному слоту
type TimeSlotStats struct {
	Label        string  `json:"label"`
	RoundsPlayed int     `json:"rounds_played"`
	Wins         int     `json:"wins"`
	WinRate      float64 `json:"win_rate"`
	TotalKills   int     `json:"total_kills"`
	TotalDeaths  int     `json:"total_deaths"`
	AvgKD        float64 `json:"avg_kd"`
	TotalDamage  int     `json:"total_damage"`
	AvgADR       float64 `json:"avg_adr"`
}

// PlayerPairStats статистика для пары игроков
type PlayerPairStats struct {
	Player1        string  `json:"player1"`
	Player2        string  `json:"player2"`
	RoundsTogether int     `json:"rounds_together"`
	Wins           int     `json:"wins"`
	WinRate        float64 `json:"win_rate"`
	AvgKD          float64 `json:"avg_kd"`
	TotalKills     int     `json:"total_kills"`
	TotalDeaths    int     `json:"total_deaths"`
}

// MapStats статистика по одной карте
type MapStats struct {
	MapName     string           `json:"map_name"`
	TotalRounds int              `json:"total_rounds"`
	TWins       int              `json:"t_wins"`
	CTWins      int              `json:"ct_wins"`
	TWinRate    float64          `json:"t_win_rate"`
	CTWinRate   float64          `json:"ct_win_rate"`
	PlayerStats []PlayerMapStats `json:"player_stats"`
	TotalKills  int              `json:"total_kills"`
	TotalDeaths int              `json:"total_deaths"`
	AvgKD       float64          `json:"avg_kd"`
}

// PlayerMapStats статистика игрока на карте
type PlayerMapStats struct {
	MapName      string  `json:"map_name"`
	RoundsPlayed int     `json:"rounds_played"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	Assists      int     `json:"assists"`
	Damage       int     `json:"damage"`
	WinRounds    int     `json:"win_rounds"`
	KD           float64 `json:"kd"`
	ADR          float64 `json:"adr"`
	WinRate      float64 `json:"win_rate"`
}

// PlayerSideStats статистика игрока на одной стороне (T или CT)
type PlayerSideStats struct {
	RoundsPlayed int     `json:"rounds_played"`
	Kills        int     `json:"kills"`
	Deaths       int     `json:"deaths"`
	Assists      int     `json:"assists"`
	Damage       int     `json:"damage"`
	WinRounds    int     `json:"win_rounds"`
	KD           float64 `json:"kd"`
	ADR          float64 `json:"adr"`
	WinRate      float64 `json:"win_rate"`
}

// PlayerTvsCTStats статистика игрока T vs CT
type PlayerTvsCTStats struct {
	AccountID     int64           `json:"account_id"`
	Name          string          `json:"name"`
	TStats        PlayerSideStats `json:"t_stats"`
	CTStats       PlayerSideStats `json:"ct_stats"`
	TotalRounds   int             `json:"total_rounds"`
	PreferredSide string          `json:"preferred_side"` // "T", "CT", или "Balanced"
	KDDiff        float64         `json:"kd_diff"`        // T K/D - CT K/D
}

// PlayerProgress представляет данные прогресса игрока по датам
type PlayerProgress struct {
	AccountID     int64              `json:"account_id"`
	Name          string             `json:"name"`
	Daily         []DailyPlayerStats `json:"daily"`
	Totals        PlayerTotalStats   `json:"totals"`
	ByHour        []TimeSlotStats    `json:"by_hour"`
	ByDayOfWeek   []TimeSlotStats    `json:"by_day_of_week"`
	BestTimeSlot  string             `json:"best_time_slot"`
	WorstTimeSlot string             `json:"worst_time_slot"`
	TopPartners   []PlayerPairStats  `json:"top_partners"`
	MapStats      []PlayerMapStats   `json:"map_stats"`   // Статистика игрока по картам
	TvsCTStats    *PlayerTvsCTStats  `json:"tvsct_stats"` // T vs CT статистика игрока
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

// ProgressData общие данные для всего таба
type ProgressData struct {
	Players       []PlayerProgress   `json:"players"`
	OverallByHour []TimeSlotStats    `json:"overall_by_hour"`
	OverallByDay  []TimeSlotStats    `json:"overall_by_day"`
	TopPairs      []PlayerPairStats  `json:"top_pairs"`
	MapStats      []MapStats         `json:"map_stats"`   // Общая статистика по картам
	TvsCTStats    []PlayerTvsCTStats `json:"tvsct_stats"` // T vs CT статистика всех игроков
}

// GenerateHTML генерирует HTML для таба прогресса
func (p *ProgressTabComponent) GenerateHTML() string {
	return `
<!-- PLAYER PROGRESS & STATS -->
<div id="tab-progress" class="view">
  <div class="toolbar">
    <label style="display:flex;align-items:center;gap:8px">
      <span style="color:var(--text);font-weight:600;">Выберите игрока:</span>
      <select id="playerSelect" style="background:var(--panel);color:var(--text);border:1px solid rgba(124,92,255,0.3);border-radius:6px;padding:8px 12px;font-size:14px;cursor:pointer;min-width:200px;">
        <option value="">Все игроки (общая статистика)</option>
      </select>
    </label>
    <span class="small" style="margin-left:auto;" id="progressInfo">Статистика и прогресс</span>
  </div>

  <!-- Общая статистика (все игроки) -->
  <div id="overallStats" style="padding:20px;">
    <div style="background:var(--panel);padding:20px;border-radius:12px;margin-bottom:24px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">🕐 Активность по времени суток</h3>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;">
        <div>
          <canvas id="overallHourChart" style="max-height:300px;"></canvas>
        </div>
        <div>
          <canvas id="overallHourWRChart" style="max-height:300px;"></canvas>
        </div>
      </div>
    </div>

    <div style="background:var(--panel);padding:20px;border-radius:12px;margin-bottom:24px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">📅 Активность по дням недели</h3>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;">
        <div>
          <canvas id="overallDayChart" style="max-height:300px;"></canvas>
        </div>
        <div>
          <canvas id="overallDayWRChart" style="max-height:300px;"></canvas>
        </div>
      </div>
    </div>

    <div style="background:var(--panel);padding:20px;border-radius:12px;margin-bottom:24px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">🤝 Топ пары игроков</h3>
      <div id="topPairsContent"></div>
    </div>

    <div style="background:var(--panel);padding:20px;border-radius:12px;margin-bottom:24px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">🗺️ Статистика по картам</h3>
      <div id="mapStatsContent"></div>
    </div>

    <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">⚔️ T vs CT Статистика</h3>
      <div id="tvsctStatsContent"></div>
    </div>
  </div>

  <!-- Статистика конкретного игрока -->
  <div id="playerStats" style="display:none;padding:20px;">
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

    <!-- Графики прогресса -->
    <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(500px,1fr));gap:20px;margin-bottom:24px;">
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">📈 EPI по датам</h4>
        <canvas id="chartEPI" style="max-height:300px;"></canvas>
      </div>
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">⚔️ K/D Ratio по датам</h4>
        <canvas id="chartKD" style="max-height:300px;"></canvas>
      </div>
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">💥 Урон по датам</h4>
        <canvas id="chartDamage" style="max-height:300px;"></canvas>
      </div>
      <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
        <h4 style="margin:0 0 16px;color:var(--accent);font-size:16px;">🏆 Win Rate по датам</h4>
        <canvas id="chartWinRate" style="max-height:300px;"></canvas>
      </div>
    </div>

    <!-- Лучшее время -->
    <div style="background:var(--panel);padding:20px;border-radius:12px;margin-bottom:24px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">⏰ Лучшее время для игры</h3>
      <div style="display:grid;grid-template-columns:1fr 1fr;gap:24px;margin-bottom:24px;">
        <div style="padding:16px;background:linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);border-radius:8px;">
          <div style="font-size:12px;color:var(--muted);margin-bottom:8px;">🌟 Лучшее время</div>
          <div id="playerBestTime" style="font-size:20px;font-weight:bold;color:#22c55e;">-</div>
        </div>
        <div style="padding:16px;background:linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);border-radius:8px;">
          <div style="font-size:12px;color:var(--muted);margin-bottom:8px;">⚠️ Худшее время</div>
          <div id="playerWorstTime" style="font-size:20px;font-weight:bold;color:#ef4444;">-</div>
        </div>
      </div>
      <div style="display:grid;grid-template-columns:1fr;gap:20px;">
        <canvas id="playerHourChart" style="max-height:250px;"></canvas>
        <canvas id="playerDayChart" style="max-height:250px;"></canvas>
      </div>
    </div>

    <!-- Лучшие партнеры -->
    <div style="background:var(--panel);padding:20px;border-radius:12px;margin-bottom:24px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">🤝 Лучшие партнеры</h3>
      <div id="playerPartnersContent"></div>
    </div>

    <!-- Статистика по картам -->
    <div style="background:var(--panel);padding:20px;border-radius:12px;margin-bottom:24px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">🗺️ Статистика по картам</h3>
      <div id="playerMapStatsContent"></div>
    </div>

    <!-- T vs CT -->
    <div style="background:var(--panel);padding:20px;border-radius:12px;border:1px solid rgba(124,92,255,0.1);">
      <h3 style="margin:0 0 16px;color:var(--accent);font-size:18px;">⚔️ T vs CT</h3>
      <div id="playerTvsCTContent"></div>
    </div>
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
// Init: Player Progress with Synergy and Time Analysis
(function() {
  const data = %s;
  const playerSelect = document.getElementById('playerSelect');
  const overallStatsDiv = document.getElementById('overallStats');
  const playerStatsDiv = document.getElementById('playerStats');

  let charts = {};

  // Заполняем список игроков
  data.players.forEach(player => {
    const option = document.createElement('option');
    option.value = player.account_id;
    option.textContent = player.name;
    playerSelect.appendChild(option);
  });

  // Обработчик выбора игрока
  playerSelect.addEventListener('change', function() {
    const accountId = parseInt(this.value);
    if (!accountId) {
      // Показываем общую статистику
      showOverallStats();
      return;
    }

    const player = data.players.find(p => p.account_id === accountId);
    if (!player) return;

    // Показываем статистику игрока
    showPlayerStats(player);
  });

  function showOverallStats() {
    overallStatsDiv.style.display = 'block';
    playerStatsDiv.style.display = 'none';

    // Уничтожаем старые графики
    Object.values(charts).forEach(chart => chart.destroy());
    charts = {};

    // График активности по часам
    charts.overallHour = new Chart(document.getElementById('overallHourChart'), {
      type: 'bar',
      data: {
        labels: data.overall_by_hour.map(slot => slot.label),
        datasets: [{
          label: 'Раундов',
          data: data.overall_by_hour.map(slot => slot.rounds_played),
          backgroundColor: 'rgba(124, 92, 255, 0.6)',
          borderColor: '#7c5cff',
          borderWidth: 1
        }]
      },
      options: getChartOptions('Активность по часам')
    });

    // График винрейта по часам
    const hourWROptions = getChartOptions('Win Rate по часам');
    hourWROptions.scales.y.max = 100;
    charts.overallHourWR = new Chart(document.getElementById('overallHourWRChart'), {
      type: 'line',
      data: {
        labels: data.overall_by_hour.map(slot => slot.label),
        datasets: [{
          label: 'Win Rate',
          data: data.overall_by_hour.map(slot => slot.win_rate),
          borderColor: '#22c55e',
          backgroundColor: 'rgba(34, 197, 94, 0.1)',
          borderWidth: 2,
          tension: 0.3,
          fill: true
        }]
      },
      options: hourWROptions
    });

    // График активности по дням недели
    charts.overallDay = new Chart(document.getElementById('overallDayChart'), {
      type: 'bar',
      data: {
        labels: data.overall_by_day.map(slot => slot.label),
        datasets: [{
          label: 'Раундов',
          data: data.overall_by_day.map(slot => slot.rounds_played),
          backgroundColor: 'rgba(59, 130, 246, 0.6)',
          borderColor: '#3b82f6',
          borderWidth: 1
        }]
      },
      options: getChartOptions('Активность по дням недели')
    });

    // График винрейта по дням недели
    const dayWROptions = getChartOptions('Win Rate по дням недели');
    dayWROptions.scales.y.max = 100;
    charts.overallDayWR = new Chart(document.getElementById('overallDayWRChart'), {
      type: 'line',
      data: {
        labels: data.overall_by_day.map(slot => slot.label),
        datasets: [{
          label: 'Win Rate',
          data: data.overall_by_day.map(slot => slot.win_rate),
          borderColor: '#8b5cf6',
          backgroundColor: 'rgba(139, 92, 246, 0.1)',
          borderWidth: 2,
          tension: 0.3,
          fill: true
        }]
      },
      options: dayWROptions
    });

    // Топ пары
    renderTopPairs();

    // Статистика по картам
    renderMapStats();

    // T vs CT статистика
    renderTvsCTStats();
  }

  function showPlayerStats(player) {
    overallStatsDiv.style.display = 'none';
    playerStatsDiv.style.display = 'block';

    // Уничтожаем старые графики
    Object.values(charts).forEach(chart => chart.destroy());
    charts = {};

    // Обновляем общую статистику
    document.getElementById('stat-rounds').textContent = player.totals.rounds_played;
    document.getElementById('stat-epi').textContent = player.totals.avg_epi.toFixed(2);
    document.getElementById('stat-kd').textContent = player.totals.avg_kd.toFixed(2);
    document.getElementById('stat-adr').textContent = Math.round(player.totals.avg_adr);
    document.getElementById('stat-winrate').textContent = player.totals.win_rate.toFixed(1) + '%%';

    const dates = player.daily.map(d => d.date);

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
      options: getChartOptions('Эффективность (EPI)')
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
      options: getChartOptions('K/D Ratio')
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
      options: getChartOptions('Средний урон за раунд')
    });

    // График Win Rate
    const winRateOptions = getChartOptions('Процент побед');
    winRateOptions.scales.y.max = 100;
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
      options: winRateOptions
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
      options: getChartOptions('Активность')
    });

    // График K/D/A
    const kdaOptions = getChartOptions('Kills / Deaths / Assists');
    kdaOptions.plugins.legend = { display: true, labels: { color: '#e5e5e5' } };
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
      options: kdaOptions
    });

    // График времени игрока (по часам)
    if (player.by_hour && player.by_hour.length > 0) {
      const playerHourOptions = getChartOptions('Ваш Win Rate по часам');
      playerHourOptions.scales.y.max = 100;
      charts.playerHour = new Chart(document.getElementById('playerHourChart'), {
        type: 'bar',
        data: {
          labels: player.by_hour.map(slot => slot.label),
          datasets: [{
            label: 'Win Rate',
            data: player.by_hour.map(slot => slot.win_rate),
            backgroundColor: player.by_hour.map(slot =>
              slot.win_rate >= 60 ? 'rgba(34, 197, 94, 0.7)' :
              slot.win_rate >= 50 ? 'rgba(245, 158, 11, 0.7)' : 'rgba(239, 68, 68, 0.7)'
            ),
            borderWidth: 1
          }]
        },
        options: playerHourOptions
      });
    }

    // График времени игрока (по дням недели)
    if (player.by_day_of_week && player.by_day_of_week.length > 0) {
      const playerDayOptions = getChartOptions('Ваш Win Rate по дням недели');
      playerDayOptions.scales.y.max = 100;
      charts.playerDay = new Chart(document.getElementById('playerDayChart'), {
        type: 'bar',
        data: {
          labels: player.by_day_of_week.map(slot => slot.label),
          datasets: [{
            label: 'Win Rate',
            data: player.by_day_of_week.map(slot => slot.win_rate),
            backgroundColor: player.by_day_of_week.map(slot =>
              slot.win_rate >= 60 ? 'rgba(34, 197, 94, 0.7)' :
              slot.win_rate >= 50 ? 'rgba(245, 158, 11, 0.7)' : 'rgba(239, 68, 68, 0.7)'
            ),
            borderWidth: 1
          }]
        },
        options: playerDayOptions
      });
    }

    // Лучшие партнеры
    renderPlayerPartners(player);

    // Статистика игрока по картам
    renderPlayerMapStats(player);

    // T vs CT статистика игрока
    renderPlayerTvsCT(player);
  }

  function renderTopPairs() {
    const div = document.getElementById('topPairsContent');
    if (!data.top_pairs || data.top_pairs.length === 0) {
      div.innerHTML = '<div style="text-align:center;padding:40px;color:var(--muted);">Недостаточно данных</div>';
      return;
    }

    let html = '<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:16px;">';
    data.top_pairs.slice(0, 10).forEach((pair, index) => {
      const medal = index === 0 ? '🥇' : index === 1 ? '🥈' : index === 2 ? '🥉' : (index + 1) + '.';
      const winRateColor = pair.win_rate >= 60 ? '#22c55e' : pair.win_rate >= 50 ? '#fde047' : '#ef4444';

      html += '<div style="padding:16px;background:linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);border-radius:8px;border:1px solid rgba(124,92,255,0.2);">' +
        '<div style="display:flex;align-items:center;gap:8px;margin-bottom:12px;">' +
          '<span style="font-size:20px;">' + medal + '</span>' +
          '<div style="flex:1;">' +
            '<div style="font-size:13px;font-weight:bold;color:#e5e5e5;">' + pair.player1 + ' + ' + pair.player2 + '</div>' +
            '<div style="font-size:10px;color:var(--muted);">' + pair.rounds_together + ' раундов вместе</div>' +
          '</div>' +
        '</div>' +
        '<div style="display:grid;grid-template-columns:1fr 1fr;gap:8px;">' +
          '<div style="text-align:center;padding:8px;background:rgba(0,0,0,0.3);border-radius:6px;">' +
            '<div style="font-size:16px;font-weight:bold;color:' + winRateColor + '">' + pair.win_rate.toFixed(1) + '%%</div>' +
            '<div style="font-size:9px;color:var(--muted);">Win Rate</div>' +
          '</div>' +
          '<div style="text-align:center;padding:8px;background:rgba(0,0,0,0.3);border-radius:6px;">' +
            '<div style="font-size:16px;font-weight:bold;color:' + (pair.avg_kd >= 1 ? '#22c55e' : '#ef4444') + '">' + pair.avg_kd.toFixed(2) + '</div>' +
            '<div style="font-size:9px;color:var(--muted);">K/D</div>' +
          '</div>' +
        '</div>' +
      '</div>';
    });
    html += '</div>';
    div.innerHTML = html;
  }

  function renderPlayerPartners(player) {
    const div = document.getElementById('playerPartnersContent');
    if (!player.top_partners || player.top_partners.length === 0) {
      div.innerHTML = '<div style="text-align:center;padding:20px;color:var(--muted);">Недостаточно данных</div>';
      return;
    }

    let html = '<table style="width:100%;"><thead><tr><th>Партнер</th><th>Раундов</th><th>Win Rate</th><th>K/D</th></tr></thead><tbody>';
    player.top_partners.forEach(pair => {
      const partner = pair.player1 === player.name ? pair.player2 : pair.player1;
      const winRateColor = pair.win_rate >= 60 ? '#22c55e' : pair.win_rate >= 50 ? '#fde047' : '#ef4444';
      html += '<tr>' +
        '<td>' + partner + '</td>' +
        '<td>' + pair.rounds_together + '</td>' +
        '<td style="color:' + winRateColor + ';font-weight:bold;">' + pair.win_rate.toFixed(1) + '%%</td>' +
        '<td style="color:' + (pair.avg_kd >= 1 ? '#22c55e' : '#ef4444') + '">' + pair.avg_kd.toFixed(2) + '</td>' +
      '</tr>';
    });
    html += '</tbody></table>';
    div.innerHTML = html;
  }

  function renderMapStats() {
    const div = document.getElementById('mapStatsContent');
    if (!data.map_stats || data.map_stats.length === 0) {
      div.innerHTML = '<div style="text-align:center;padding:40px;color:var(--muted);">Недостаточно данных</div>';
      return;
    }

    let html = '<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:16px;">';
    data.map_stats.forEach(mapStat => {
      const tWinColor = mapStat.t_win_rate >= 55 ? '#22c55e' : mapStat.t_win_rate >= 45 ? '#fde047' : '#ef4444';
      const ctWinColor = mapStat.ct_win_rate >= 55 ? '#22c55e' : mapStat.ct_win_rate >= 45 ? '#fde047' : '#ef4444';

      html += '<div style="padding:20px;background:linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);border-radius:8px;border:1px solid rgba(124,92,255,0.2);">' +
        '<div style="font-size:16px;font-weight:bold;color:#e5e5e5;margin-bottom:12px;text-align:center;">' + mapStat.map_name + '</div>' +
        '<div style="font-size:12px;color:var(--muted);text-align:center;margin-bottom:16px;">' + mapStat.total_rounds + ' раундов</div>' +
        '<div style="display:grid;grid-template-columns:1fr 1fr;gap:12px;">' +
          '<div style="text-align:center;padding:12px;background:rgba(0,0,0,0.3);border-radius:6px;">' +
            '<div style="font-size:20px;font-weight:bold;color:' + tWinColor + '">' + mapStat.t_win_rate.toFixed(1) + '%%</div>' +
            '<div style="font-size:10px;color:var(--muted);margin-top:4px;">T Win Rate</div>' +
          '</div>' +
          '<div style="text-align:center;padding:12px;background:rgba(0,0,0,0.3);border-radius:6px;">' +
            '<div style="font-size:20px;font-weight:bold;color:' + ctWinColor + '">' + mapStat.ct_win_rate.toFixed(1) + '%%</div>' +
            '<div style="font-size:10px;color:var(--muted);margin-top:4px;">CT Win Rate</div>' +
          '</div>' +
        '</div>' +
      '</div>';
    });
    html += '</div>';
    div.innerHTML = html;
  }

  function renderTvsCTStats() {
    const div = document.getElementById('tvsctStatsContent');
    if (!data.tvsct_stats || data.tvsct_stats.length === 0) {
      div.innerHTML = '<div style="text-align:center;padding:40px;color:var(--muted);">Недостаточно данных</div>';
      return;
    }

    let html = '<table style="width:100%;"><thead><tr>' +
      '<th>Игрок</th>' +
      '<th>Предпочтение</th>' +
      '<th>T Раундов</th>' +
      '<th>T K/D</th>' +
      '<th>T WR%%</th>' +
      '<th>CT Раундов</th>' +
      '<th>CT K/D</th>' +
      '<th>CT WR%%</th>' +
    '</tr></thead><tbody>';

    data.tvsct_stats.forEach(tvs => {
      const preferredColor = tvs.preferred_side === 'T' ? '#f59e0b' : tvs.preferred_side === 'CT' ? '#3b82f6' : '#8b5cf6';
      const tKDColor = tvs.t_stats.kd >= 1 ? '#22c55e' : '#ef4444';
      const ctKDColor = tvs.ct_stats.kd >= 1 ? '#22c55e' : '#ef4444';
      const tWRColor = tvs.t_stats.win_rate >= 50 ? '#22c55e' : '#ef4444';
      const ctWRColor = tvs.ct_stats.win_rate >= 50 ? '#22c55e' : '#ef4444';

      html += '<tr>' +
        '<td>' + tvs.name + '</td>' +
        '<td style="color:' + preferredColor + ';font-weight:bold;">' + tvs.preferred_side + '</td>' +
        '<td>' + tvs.t_stats.rounds_played + '</td>' +
        '<td style="color:' + tKDColor + '">' + tvs.t_stats.kd.toFixed(2) + '</td>' +
        '<td style="color:' + tWRColor + '">' + tvs.t_stats.win_rate.toFixed(1) + '%%</td>' +
        '<td>' + tvs.ct_stats.rounds_played + '</td>' +
        '<td style="color:' + ctKDColor + '">' + tvs.ct_stats.kd.toFixed(2) + '</td>' +
        '<td style="color:' + ctWRColor + '">' + tvs.ct_stats.win_rate.toFixed(1) + '%%</td>' +
      '</tr>';
    });
    html += '</tbody></table>';
    div.innerHTML = html;
  }

  function renderPlayerMapStats(player) {
    const div = document.getElementById('playerMapStatsContent');
    if (!player.map_stats || player.map_stats.length === 0) {
      div.innerHTML = '<div style="text-align:center;padding:20px;color:var(--muted);">Недостаточно данных</div>';
      return;
    }

    let html = '<table style="width:100%;"><thead><tr>' +
      '<th>Карта</th>' +
      '<th>Раундов</th>' +
      '<th>K/D</th>' +
      '<th>ADR</th>' +
      '<th>Win Rate</th>' +
    '</tr></thead><tbody>';

    player.map_stats.forEach(mapStat => {
      const kdColor = mapStat.kd >= 1 ? '#22c55e' : '#ef4444';
      const wrColor = mapStat.win_rate >= 50 ? '#22c55e' : '#ef4444';

      html += '<tr>' +
        '<td>' + mapStat.map_name + '</td>' +
        '<td>' + mapStat.rounds_played + '</td>' +
        '<td style="color:' + kdColor + ';font-weight:bold;">' + mapStat.kd.toFixed(2) + '</td>' +
        '<td>' + Math.round(mapStat.adr) + '</td>' +
        '<td style="color:' + wrColor + ';font-weight:bold;">' + mapStat.win_rate.toFixed(1) + '%%</td>' +
      '</tr>';
    });
    html += '</tbody></table>';
    div.innerHTML = html;
  }

  function renderPlayerTvsCT(player) {
    const div = document.getElementById('playerTvsCTContent');
    if (!player.tvsct_stats) {
      div.innerHTML = '<div style="text-align:center;padding:20px;color:var(--muted);">Недостаточно данных</div>';
      return;
    }

    const tvs = player.tvsct_stats;
    const preferredColor = tvs.preferred_side === 'T' ? '#f59e0b' : tvs.preferred_side === 'CT' ? '#3b82f6' : '#8b5cf6';

    let html = '<div style="margin-bottom:20px;text-align:center;">' +
      '<div style="font-size:14px;color:var(--muted);margin-bottom:8px;">Предпочтительная сторона</div>' +
      '<div style="font-size:24px;font-weight:bold;color:' + preferredColor + ';">' + tvs.preferred_side + '</div>' +
      '<div style="font-size:12px;color:var(--muted);margin-top:4px;">Разница K/D: ' + tvs.kd_diff.toFixed(2) + '</div>' +
    '</div>';

    html += '<div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;">';

    // T side
    const tKDColor = tvs.t_stats.kd >= 1 ? '#22c55e' : '#ef4444';
    const tWRColor = tvs.t_stats.win_rate >= 50 ? '#22c55e' : '#ef4444';
    html += '<div style="padding:20px;background:linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);border-radius:8px;border:1px solid rgba(245,158,11,0.3);">' +
      '<div style="font-size:18px;font-weight:bold;color:#f59e0b;margin-bottom:16px;text-align:center;">Terrorists</div>' +
      '<div style="display:grid;gap:12px;">' +
        '<div style="display:flex;justify-content:space-between;">' +
          '<span style="color:var(--muted);">Раундов:</span>' +
          '<span style="color:#e5e5e5;font-weight:bold;">' + tvs.t_stats.rounds_played + '</span>' +
        '</div>' +
        '<div style="display:flex;justify-content:space-between;">' +
          '<span style="color:var(--muted);">K/D:</span>' +
          '<span style="color:' + tKDColor + ';font-weight:bold;">' + tvs.t_stats.kd.toFixed(2) + '</span>' +
        '</div>' +
        '<div style="display:flex;justify-content:space-between;">' +
          '<span style="color:var(--muted);">ADR:</span>' +
          '<span style="color:#e5e5e5;font-weight:bold;">' + Math.round(tvs.t_stats.adr) + '</span>' +
        '</div>' +
        '<div style="display:flex;justify-content:space-between;">' +
          '<span style="color:var(--muted);">Win Rate:</span>' +
          '<span style="color:' + tWRColor + ';font-weight:bold;">' + tvs.t_stats.win_rate.toFixed(1) + '%%</span>' +
        '</div>' +
      '</div>' +
    '</div>';

    // CT side
    const ctKDColor = tvs.ct_stats.kd >= 1 ? '#22c55e' : '#ef4444';
    const ctWRColor = tvs.ct_stats.win_rate >= 50 ? '#22c55e' : '#ef4444';
    html += '<div style="padding:20px;background:linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);border-radius:8px;border:1px solid rgba(59,130,246,0.3);">' +
      '<div style="font-size:18px;font-weight:bold;color:#3b82f6;margin-bottom:16px;text-align:center;">Counter-Terrorists</div>' +
      '<div style="display:grid;gap:12px;">' +
        '<div style="display:flex;justify-content:space-between;">' +
          '<span style="color:var(--muted);">Раундов:</span>' +
          '<span style="color:#e5e5e5;font-weight:bold;">' + tvs.ct_stats.rounds_played + '</span>' +
        '</div>' +
        '<div style="display:flex;justify-content:space-between;">' +
          '<span style="color:var(--muted);">K/D:</span>' +
          '<span style="color:' + ctKDColor + ';font-weight:bold;">' + tvs.ct_stats.kd.toFixed(2) + '</span>' +
        '</div>' +
        '<div style="display:flex;justify-content:space-between;">' +
          '<span style="color:var(--muted);">ADR:</span>' +
          '<span style="color:#e5e5e5;font-weight:bold;">' + Math.round(tvs.ct_stats.adr) + '</span>' +
        '</div>' +
        '<div style="display:flex;justify-content:space-between;">' +
          '<span style="color:var(--muted);">Win Rate:</span>' +
          '<span style="color:' + ctWRColor + ';font-weight:bold;">' + tvs.ct_stats.win_rate.toFixed(1) + '%%</span>' +
        '</div>' +
      '</div>' +
    '</div>';

    html += '</div>';
    div.innerHTML = html;
  }

  function getChartOptions(title) {
    return {
      responsive: true,
      maintainAspectRatio: true,
      plugins: {
        title: title ? { display: true, text: title, color: '#e5e5e5' } : undefined,
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
  }

  // Показываем общую статистику по умолчанию
  showOverallStats();
})();
`, string(jsonData))
}

// buildProgressData собирает данные прогресса для всех игроков + синергия + временная аналитика
func (p *ProgressTabComponent) buildProgressData(data *stats.StatsData) *ProgressData {
	result := &ProgressData{
		OverallByHour: make([]TimeSlotStats, 24),
		OverallByDay:  make([]TimeSlotStats, 7),
	}

	// Инициализируем слоты времени
	for i := 0; i < 24; i++ {
		result.OverallByHour[i].Label = fmt.Sprintf("%02d:00", i)
	}
	dayNames := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	for i := 0; i < 7; i++ {
		result.OverallByDay[i].Label = dayNames[i]
	}

	// Создаем карты для сбора данных
	playerMap := make(map[int64]*PlayerProgress)
	pairStatsMap := make(map[string]*PlayerPairStats)
	playerNames := make(map[int64]string)

	// Заполняем имена игроков
	for _, rating := range data.PlayerRatings {
		playerNames[rating.AccountID] = rating.Name
		// Инициализируем структуру для каждого игрока
		playerMap[rating.AccountID] = &PlayerProgress{
			AccountID:   rating.AccountID,
			Name:        rating.Name,
			ByHour:      make([]TimeSlotStats, 24),
			ByDayOfWeek: make([]TimeSlotStats, 7),
		}
		for i := 0; i < 24; i++ {
			playerMap[rating.AccountID].ByHour[i].Label = fmt.Sprintf("%02d:00", i)
		}
		for i := 0; i < 7; i++ {
			playerMap[rating.AccountID].ByDayOfWeek[i].Label = dayNames[i]
		}
	}

	// Создаем карту дневной статистики: accountID -> date -> stats
	dailyMap := make(map[int64]map[string]*DailyPlayerStats)

	// Обрабатываем раунды по датам
	for date, rounds := range data.DailyRounds {
		// Парсим дату для получения дня недели
		parsedDate, err := time.Parse("2006-01-02", date)
		dayOfWeek := 0
		if err == nil {
			// Конвертируем в Monday = 0
			dow := int(parsedDate.Weekday())
			if dow == 0 {
				dayOfWeek = 6
			} else {
				dayOfWeek = dow - 1
			}
		}
		hour := 12 // Используем фиксированный час, так как нет точного времени

		for _, round := range rounds {
			// Группируем игроков по командам для синергии
			teams := make(map[int][]int64) // team -> []accountID
			playerData := make(map[int64]struct {
				kills  int
				deaths int
				damage int
			})

			// Обрабатываем каждого игрока в раунде
			for _, playerStat := range round.Players {
				accountID := playerStat.AccountID

				// Пропускаем игроков без имен
				if playerNames[accountID] == "" {
					continue
				}

				// Добавляем в команду для синергии
				teams[playerStat.Team] = append(teams[playerStat.Team], accountID)
				playerData[accountID] = struct {
					kills  int
					deaths int
					damage int
				}{
					kills:  playerStat.Kills,
					deaths: playerStat.Deaths,
					damage: playerStat.Damage,
				}

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
				isWin := round.Winner == playerStat.Team
				if isWin {
					daily.WinRounds++
				}

				// Обновляем временные слоты игрока
				player := playerMap[accountID]
				if player != nil {
					// По часам
					player.ByHour[hour].RoundsPlayed++
					player.ByHour[hour].TotalKills += playerStat.Kills
					player.ByHour[hour].TotalDeaths += playerStat.Deaths
					player.ByHour[hour].TotalDamage += playerStat.Damage
					if isWin {
						player.ByHour[hour].Wins++
					}

					// По дням недели
					player.ByDayOfWeek[dayOfWeek].RoundsPlayed++
					player.ByDayOfWeek[dayOfWeek].TotalKills += playerStat.Kills
					player.ByDayOfWeek[dayOfWeek].TotalDeaths += playerStat.Deaths
					player.ByDayOfWeek[dayOfWeek].TotalDamage += playerStat.Damage
					if isWin {
						player.ByDayOfWeek[dayOfWeek].Wins++
					}
				}

				// Обновляем общую статистику
				result.OverallByHour[hour].RoundsPlayed++
				result.OverallByHour[hour].TotalKills += playerStat.Kills
				result.OverallByHour[hour].TotalDeaths += playerStat.Deaths
				result.OverallByHour[hour].TotalDamage += playerStat.Damage
				if isWin {
					result.OverallByHour[hour].Wins++
				}

				result.OverallByDay[dayOfWeek].RoundsPlayed++
				result.OverallByDay[dayOfWeek].TotalKills += playerStat.Kills
				result.OverallByDay[dayOfWeek].TotalDeaths += playerStat.Deaths
				result.OverallByDay[dayOfWeek].TotalDamage += playerStat.Damage
				if isWin {
					result.OverallByDay[dayOfWeek].Wins++
				}
			}

			// Обрабатываем синергию - создаем пары из игроков одной команды
			for team, players := range teams {
				if len(players) < 2 {
					continue
				}

				isWin := round.Winner == team

				// Создаем все возможные пары
				for i := 0; i < len(players); i++ {
					for j := i + 1; j < len(players); j++ {
						player1ID := players[i]
						player2ID := players[j]

						player1Name := playerNames[player1ID]
						player2Name := playerNames[player2ID]

						// Создаем уникальный ключ
						var pairKey string
						var p1, p2 string
						if player1Name < player2Name {
							pairKey = player1Name + "|" + player2Name
							p1, p2 = player1Name, player2Name
						} else {
							pairKey = player2Name + "|" + player1Name
							p1, p2 = player2Name, player1Name
						}

						// Инициализируем пару
						if pairStatsMap[pairKey] == nil {
							pairStatsMap[pairKey] = &PlayerPairStats{
								Player1: p1,
								Player2: p2,
							}
						}

						pair := pairStatsMap[pairKey]
						pair.RoundsTogether++
						if isWin {
							pair.Wins++
						}
						pair.TotalKills += playerData[player1ID].kills + playerData[player2ID].kills
						pair.TotalDeaths += playerData[player1ID].deaths + playerData[player2ID].deaths
					}
				}
			}
		}
	}

	// Вычисляем метрики для общей статистики по времени
	for i := range result.OverallByHour {
		slot := &result.OverallByHour[i]
		if slot.RoundsPlayed > 0 {
			slot.WinRate = (float64(slot.Wins) / float64(slot.RoundsPlayed)) * 100
			if slot.TotalDeaths > 0 {
				slot.AvgKD = float64(slot.TotalKills) / float64(slot.TotalDeaths)
			}
			slot.AvgADR = float64(slot.TotalDamage) / float64(slot.RoundsPlayed)
		}
	}

	for i := range result.OverallByDay {
		slot := &result.OverallByDay[i]
		if slot.RoundsPlayed > 0 {
			slot.WinRate = (float64(slot.Wins) / float64(slot.RoundsPlayed)) * 100
			if slot.TotalDeaths > 0 {
				slot.AvgKD = float64(slot.TotalKills) / float64(slot.TotalDeaths)
			}
			slot.AvgADR = float64(slot.TotalDamage) / float64(slot.RoundsPlayed)
		}
	}

	// Вычисляем метрики для пар и выбираем топ
	for _, pair := range pairStatsMap {
		if pair.RoundsTogether > 0 {
			pair.WinRate = (float64(pair.Wins) / float64(pair.RoundsTogether)) * 100
			if pair.TotalDeaths > 0 {
				pair.AvgKD = float64(pair.TotalKills) / float64(pair.TotalDeaths)
			}
		}
		result.TopPairs = append(result.TopPairs, *pair)
	}

	// Сортируем пары по win rate (минимум 10 раундов)
	sort.Slice(result.TopPairs, func(i, j int) bool {
		if result.TopPairs[i].RoundsTogether < 10 {
			return false
		}
		if result.TopPairs[j].RoundsTogether < 10 {
			return true
		}
		return result.TopPairs[i].WinRate > result.TopPairs[j].WinRate
	})

	// Оставляем топ 10
	if len(result.TopPairs) > 10 {
		result.TopPairs = result.TopPairs[:10]
	}

	// Обрабатываем каждого игрока
	for accountID, dateMap := range dailyMap {
		player := playerMap[accountID]
		if player == nil {
			continue
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

				// Упрощенный расчет EPI
				daily.EPI = (float64(daily.Kills)*1.5 - float64(daily.Deaths) + float64(daily.Assists)*0.5 + float64(daily.Damage)/100.0) / float64(daily.RoundsPlayed)
			}

			dailyStats = append(dailyStats, *daily)

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

		player.Daily = dailyStats
		player.Totals = PlayerTotalStats{
			RoundsPlayed: totalRounds,
			TotalKills:   totalKills,
			TotalDeaths:  totalDeaths,
			AvgEPI:       avgEPI,
			AvgKD:        avgKD,
			AvgADR:       avgADR,
			WinRate:      winRate,
		}

		// Вычисляем метрики для временных слотов игрока
		bestWR := 0.0
		worstWR := 100.0
		bestSlot := ""
		worstSlot := ""

		for i := range player.ByHour {
			slot := &player.ByHour[i]
			if slot.RoundsPlayed > 0 {
				slot.WinRate = (float64(slot.Wins) / float64(slot.RoundsPlayed)) * 100
				if slot.TotalDeaths > 0 {
					slot.AvgKD = float64(slot.TotalKills) / float64(slot.TotalDeaths)
				}
				slot.AvgADR = float64(slot.TotalDamage) / float64(slot.RoundsPlayed)

				if slot.RoundsPlayed >= 5 {
					if slot.WinRate > bestWR {
						bestWR = slot.WinRate
						bestSlot = slot.Label + fmt.Sprintf(" (%.1f%%)", slot.WinRate)
					}
					if slot.WinRate < worstWR {
						worstWR = slot.WinRate
						worstSlot = slot.Label + fmt.Sprintf(" (%.1f%%)", slot.WinRate)
					}
				}
			}
		}

		for i := range player.ByDayOfWeek {
			slot := &player.ByDayOfWeek[i]
			if slot.RoundsPlayed > 0 {
				slot.WinRate = (float64(slot.Wins) / float64(slot.RoundsPlayed)) * 100
				if slot.TotalDeaths > 0 {
					slot.AvgKD = float64(slot.TotalKills) / float64(slot.TotalDeaths)
				}
				slot.AvgADR = float64(slot.TotalDamage) / float64(slot.RoundsPlayed)
			}
		}

		if bestSlot != "" {
			player.BestTimeSlot = bestSlot
		} else {
			player.BestTimeSlot = "Недостаточно данных"
		}

		if worstSlot != "" {
			player.WorstTimeSlot = worstSlot
		} else {
			player.WorstTimeSlot = "Недостаточно данных"
		}

		// Находим топ партнеров для этого игрока
		var partners []PlayerPairStats
		for _, pair := range pairStatsMap {
			if (pair.Player1 == player.Name || pair.Player2 == player.Name) && pair.RoundsTogether >= 5 {
				partners = append(partners, *pair)
			}
		}

		sort.Slice(partners, func(i, j int) bool {
			return partners[i].WinRate > partners[j].WinRate
		})

		if len(partners) > 5 {
			partners = partners[:5]
		}

		player.TopPartners = partners
	}

	// Собираем данные по картам и T vs CT
	mapStatsMap := make(map[string]*MapStats)
	tvsctMap := make(map[int64]*PlayerTvsCTStats)
	playerMapStatsMap := make(map[int64]map[string]*PlayerMapStats) // accountID -> mapName -> stats

	for _, rounds := range data.DailyRounds {
		for _, round := range rounds {
			mapName := round.Map
			if mapName == "" {
				continue
			}

			// Инициализируем статистику карты
			if mapStatsMap[mapName] == nil {
				mapStatsMap[mapName] = &MapStats{
					MapName:     mapName,
					PlayerStats: []PlayerMapStats{},
				}
			}
			mapStat := mapStatsMap[mapName]
			mapStat.TotalRounds++

			// Подсчет побед T и CT
			if round.Winner == 2 {
				mapStat.TWins++
			} else if round.Winner == 3 {
				mapStat.CTWins++
			}

			// Обрабатываем игроков
			for _, player := range round.Players {
				if playerNames[player.AccountID] == "" {
					continue
				}

				// T vs CT статистика
				if tvsctMap[player.AccountID] == nil {
					tvsctMap[player.AccountID] = &PlayerTvsCTStats{
						AccountID: player.AccountID,
						Name:      playerNames[player.AccountID],
					}
				}
				tvs := tvsctMap[player.AccountID]
				tvs.TotalRounds++

				if player.Team == 2 { // T side
					tvs.TStats.RoundsPlayed++
					tvs.TStats.Kills += player.Kills
					tvs.TStats.Deaths += player.Deaths
					tvs.TStats.Assists += player.Assists
					tvs.TStats.Damage += player.Damage
					if round.Winner == 2 {
						tvs.TStats.WinRounds++
					}
				} else if player.Team == 3 { // CT side
					tvs.CTStats.RoundsPlayed++
					tvs.CTStats.Kills += player.Kills
					tvs.CTStats.Deaths += player.Deaths
					tvs.CTStats.Assists += player.Assists
					tvs.CTStats.Damage += player.Damage
					if round.Winner == 3 {
						tvs.CTStats.WinRounds++
					}
				}

				// Статистика игрока по картам
				if playerMapStatsMap[player.AccountID] == nil {
					playerMapStatsMap[player.AccountID] = make(map[string]*PlayerMapStats)
				}
				if playerMapStatsMap[player.AccountID][mapName] == nil {
					playerMapStatsMap[player.AccountID][mapName] = &PlayerMapStats{
						MapName: mapName,
					}
				}
				pms := playerMapStatsMap[player.AccountID][mapName]
				pms.RoundsPlayed++
				pms.Kills += player.Kills
				pms.Deaths += player.Deaths
				pms.Assists += player.Assists
				pms.Damage += player.Damage
				if round.Winner == player.Team {
					pms.WinRounds++
				}
			}
		}
	}

	// Вычисляем метрики для T vs CT
	for _, tvs := range tvsctMap {
		// T side
		if tvs.TStats.RoundsPlayed > 0 {
			if tvs.TStats.Deaths > 0 {
				tvs.TStats.KD = float64(tvs.TStats.Kills) / float64(tvs.TStats.Deaths)
			} else if tvs.TStats.Kills > 0 {
				tvs.TStats.KD = float64(tvs.TStats.Kills)
			}
			tvs.TStats.ADR = float64(tvs.TStats.Damage) / float64(tvs.TStats.RoundsPlayed)
			tvs.TStats.WinRate = (float64(tvs.TStats.WinRounds) / float64(tvs.TStats.RoundsPlayed)) * 100
		}

		// CT side
		if tvs.CTStats.RoundsPlayed > 0 {
			if tvs.CTStats.Deaths > 0 {
				tvs.CTStats.KD = float64(tvs.CTStats.Kills) / float64(tvs.CTStats.Deaths)
			} else if tvs.CTStats.Kills > 0 {
				tvs.CTStats.KD = float64(tvs.CTStats.Kills)
			}
			tvs.CTStats.ADR = float64(tvs.CTStats.Damage) / float64(tvs.CTStats.RoundsPlayed)
			tvs.CTStats.WinRate = (float64(tvs.CTStats.WinRounds) / float64(tvs.CTStats.RoundsPlayed)) * 100
		}

		// Определяем предпочтительную сторону
		tvs.KDDiff = tvs.TStats.KD - tvs.CTStats.KD
		if tvs.TStats.RoundsPlayed > int(float64(tvs.CTStats.RoundsPlayed)*1.2) {
			tvs.PreferredSide = "T"
		} else if tvs.CTStats.RoundsPlayed > int(float64(tvs.TStats.RoundsPlayed)*1.2) {
			tvs.PreferredSide = "CT"
		} else {
			tvs.PreferredSide = "Balanced"
		}

		result.TvsCTStats = append(result.TvsCTStats, *tvs)

		// Добавляем T vs CT статистику к игроку
		if player := playerMap[tvs.AccountID]; player != nil {
			player.TvsCTStats = tvs
		}
	}

	// Вычисляем метрики для статистики игрока по картам и привязываем к игрокам
	for accountID, mapsData := range playerMapStatsMap {
		player := playerMap[accountID]
		if player == nil {
			continue
		}

		var mapStats []PlayerMapStats
		for _, pms := range mapsData {
			if pms.RoundsPlayed > 0 {
				if pms.Deaths > 0 {
					pms.KD = float64(pms.Kills) / float64(pms.Deaths)
				} else if pms.Kills > 0 {
					pms.KD = float64(pms.Kills)
				}
				pms.ADR = float64(pms.Damage) / float64(pms.RoundsPlayed)
				pms.WinRate = (float64(pms.WinRounds) / float64(pms.RoundsPlayed)) * 100
			}
			mapStats = append(mapStats, *pms)
		}

		// Сортируем по количеству раундов
		sort.Slice(mapStats, func(i, j int) bool {
			return mapStats[i].RoundsPlayed > mapStats[j].RoundsPlayed
		})

		player.MapStats = mapStats
	}

	// Вычисляем метрики для карт
	for _, mapStat := range mapStatsMap {
		if mapStat.TotalRounds > 0 {
			mapStat.TWinRate = (float64(mapStat.TWins) / float64(mapStat.TotalRounds)) * 100
			mapStat.CTWinRate = (float64(mapStat.CTWins) / float64(mapStat.TotalRounds)) * 100
		}
		result.MapStats = append(result.MapStats, *mapStat)
	}

	// Сортируем карты по количеству раундов
	sort.Slice(result.MapStats, func(i, j int) bool {
		return result.MapStats[i].TotalRounds > result.MapStats[j].TotalRounds
	})

	// Сортируем T vs CT по имени
	sort.Slice(result.TvsCTStats, func(i, j int) bool {
		return result.TvsCTStats[i].Name < result.TvsCTStats[j].Name
	})

	// Преобразуем в массив и сортируем по имени
	for _, progress := range playerMap {
		result.Players = append(result.Players, *progress)
	}

	sort.Slice(result.Players, func(i, j int) bool {
		return result.Players[i].Name < result.Players[j].Name
	})

	return result
}
