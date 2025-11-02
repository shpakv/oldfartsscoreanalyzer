package components

import (
	"encoding/json"
	"fmt"

	"oldfartscounter/internal/stats"
)

// RoundsTabComponent отвечает за таб "Раунды"
type RoundsTabComponent struct{}

// NewRoundsTab создает новый компонент таба раундов
func NewRoundsTab() *RoundsTabComponent {
	return &RoundsTabComponent{}
}

// GenerateHTML генерирует HTML для таба раундов
func (r *RoundsTabComponent) GenerateHTML() string {
	return `
<!-- ROUNDS -->
<div id="tab-rounds" class="view">
  <div style="padding:12px;margin-bottom:16px;background:rgba(255,165,0,0.15);border-left:4px solid #ffa500;border-radius:8px;">
    <div style="display:flex;align-items:center;gap:8px;margin-bottom:6px;">
      <span style="font-size:18px;">⚠️</span>
      <strong style="color:#ffa500;">Раздел в разработке</strong>
    </div>
    <div class="small" style="color:var(--muted);">
      Данный раздел находится в активной разработке и может содержать баги и неточности. Информация отображается на основе экспериментального парсинга логов.
    </div>
  </div>

  <div class="toolbar">
    <span class="small" style="margin-left: auto;" id="roundsCount">Раундов: 0</span>
  </div>

  <div id="roundsList" style="margin-top: 20px;">
    <!-- Здесь будут отображаться раунды -->
  </div>
</div>`
}

// GenerateJS генерирует JavaScript для таба раундов
func (r *RoundsTabComponent) GenerateJS(data *stats.StatsData) string {
	// Маппинг AccountID -> Player Name
	playerMap := make(map[int64]string)

	// Парсим из событий убийств
	for _, event := range data.KillEvents {
		// Парсим SteamID для получения AccountID
		// Формат: [U:1:26840160] -> 26840160
		if len(event.KillerSID) > 6 {
			var accountID int64
			fmt.Sscanf(event.KillerSID[5:len(event.KillerSID)-1], "%d", &accountID)
			if accountID > 0 {
				playerMap[accountID] = event.KillerName
			}
		}
		if len(event.VictimSID) > 6 {
			var accountID int64
			fmt.Sscanf(event.VictimSID[5:len(event.VictimSID)-1], "%d", &accountID)
			if accountID > 0 {
				playerMap[accountID] = event.VictimName
			}
		}
	}

	// Парсим из событий флешек
	for _, event := range data.FlashEvents {
		if len(event.FlasherSID) > 6 {
			var accountID int64
			fmt.Sscanf(event.FlasherSID[5:len(event.FlasherSID)-1], "%d", &accountID)
			if accountID > 0 && playerMap[accountID] == "" {
				playerMap[accountID] = event.FlasherName
			}
		}
		if len(event.VictimSID) > 6 {
			var accountID int64
			fmt.Sscanf(event.VictimSID[5:len(event.VictimSID)-1], "%d", &accountID)
			if accountID > 0 && playerMap[accountID] == "" {
				playerMap[accountID] = event.VictimName
			}
		}
	}

	// Парсим из событий дефьюза
	for _, event := range data.DefuseEvents {
		if len(event.PlayerSID) > 6 {
			var accountID int64
			fmt.Sscanf(event.PlayerSID[5:len(event.PlayerSID)-1], "%d", &accountID)
			if accountID > 0 && playerMap[accountID] == "" {
				playerMap[accountID] = event.PlayerName
			}
		}
	}

	jRounds, _ := json.Marshal(data.RoundStats)
	jPlayerMap, _ := json.Marshal(playerMap)

	return fmt.Sprintf(`
// Init: Раунды
(function() {
  const allRounds = %s;
  const playerNames = %s;

  // Определяем состав команды Team 1 для каждой даты
  // Team 1 - команда, которая играет за CT в первом раунде первой карты дня
  const dateT1Teams = {};

  function getT1Team(date) {
    if(dateT1Teams[date]) return dateT1Teams[date];

    // Находим все раунды для этой даты и сортируем по времени и номеру раунда
    const dateRounds = allRounds.filter(r => r.Date === date);
    if(dateRounds.length === 0) return null;

    // Сортируем по времени, затем по карте, затем по номеру раунда
    dateRounds.sort((a, b) => {
      if(a.Time !== b.Time) return a.Time.localeCompare(b.Time);
      if(a.Map !== b.Map) return a.Map.localeCompare(b.Map);
      return a.RoundNumber - b.RoundNumber;
    });

    // Первый раунд определяет Team 1
    const firstRound = dateRounds[0];
    const t1Players = firstRound.Players
      .filter(p => p.Team === 3) // CT в первом раунде = Team 1
      .map(p => p.AccountID);

    dateT1Teams[date] = new Set(t1Players);
    return dateT1Teams[date];
  }

  // Функция определения команды игрока (Team 1 или Team 2)
  function getPlayerTeam(accountID, date) {
    const t1Team = getT1Team(date);
    if(!t1Team) return 0;
    return t1Team.has(accountID) ? 1 : 2; // 1 = Team 1, 2 = Team 2
  }

  // Функция подсчета счета для Team 1 и Team 2
  function getTeamScores(round) {
    const t1Team = getT1Team(round.Date);
    if(!t1Team) return { t1: 0, t2: 0 };

    // Определяем какая команда (CT или T) является Team 1
    const ctPlayers = round.Players.filter(p => p.Team === 3);
    const ctIsT1 = ctPlayers.some(p => t1Team.has(p.AccountID));

    if(ctIsT1) {
      return { t1: round.ScoreCT, t2: round.ScoreT };
    } else {
      return { t1: round.ScoreT, t2: round.ScoreCT };
    }
  }

  const roundsList = document.getElementById('roundsList');
  const roundsCount = document.getElementById('roundsCount');

  // Группируем раунды: Дата → Карта → Раунды
  function groupRounds(rounds) {
    const grouped = {};

    rounds.forEach(round => {
      const date = round.Date || 'Unknown';
      const map = round.Map || 'Unknown';

      if(!grouped[date]) grouped[date] = {};
      if(!grouped[date][map]) grouped[date][map] = [];

      grouped[date][map].push(round);
    });

    return grouped;
  }

  // Отрисовка всех раундов
  function renderRounds() {
    // Используем отфильтрованные данные если доступны
    const roundsToRender = window.filteredRoundStats || allRounds;
    const grouped = groupRounds(roundsToRender);
    const dates = Object.keys(grouped).sort().reverse(); // Новые даты сверху

    let totalRounds = 0;
    let html = '';

    dates.forEach(date => {
      // Сортируем карты по времени первого раунда
      const mapFirstTimes = {};
      Object.keys(grouped[date]).forEach(mapName => {
        const rounds = grouped[date][mapName];
        if(rounds.length > 0) {
          // Находим самое раннее время среди раундов карты
          let minTime = rounds[0].Time;
          rounds.forEach(r => {
            if(r.Time < minTime) minTime = r.Time;
          });
          mapFirstTimes[mapName] = minTime;
        }
      });

      const maps = Object.keys(grouped[date]).sort((a, b) => {
        return (mapFirstTimes[a] || '').localeCompare(mapFirstTimes[b] || '');
      });
      let dateRoundsCount = 0;
      let mapsWonT1 = 0;
      let mapsWonT2 = 0;

      // Подсчитываем количество карт, выигранных каждой командой
      maps.forEach(map => {
        const rounds = grouped[date][map];
        dateRoundsCount += rounds.length;
        if(rounds.length > 0) {
          // Находим последний раунд карты (с максимальным счетом)
          let lastRound = rounds[0];
          rounds.forEach(round => {
            const currentTotal = (round.ScoreCT || 0) + (round.ScoreT || 0);
            const lastTotal = (lastRound.ScoreCT || 0) + (lastRound.ScoreT || 0);
            if(currentTotal > lastTotal) {
              lastRound = round;
            }
          });

          // Получаем счет Team 1 и Team 2 для последнего раунда
          const scores = getTeamScores(lastRound);

          // Определяем победителя карты
          if(scores.t1 > scores.t2) {
            mapsWonT1++;
          } else if(scores.t2 > scores.t1) {
            mapsWonT2++;
          }
        }
      });

      totalRounds += dateRoundsCount;

      // Заголовок даты
      html += '<div class="date-group" style="margin-bottom:30px;">';
      html += '<div class="date-header" data-date-id="' + date + '" style="padding:12px;background:var(--accent);color:white;cursor:pointer;border-radius:8px;display:flex;align-items:center;gap:12px;">';
      html += '<span class="expand-indicator" style="flex-shrink:0;">▶</span>';
      html += '<strong style="font-size:16px;flex-shrink:0;">' + date + '</strong>';
      html += '<span style="font-size:14px;font-weight:600;background:rgba(255,255,255,0.15);padding:4px 10px;border-radius:6px;">[Team 1] ' + mapsWonT1 + ' : ' + mapsWonT2 + ' [Team 2]</span>';
      html += '<span style="font-size:12px;color:rgba(255,255,255,0.7);margin-left:auto;">Карт: ' + maps.length + '</span>';
      html += '</div>';

      // Контент даты (карты) - изначально скрыт
      html += '<div id="date-body-' + date + '" style="display:none;margin-top:10px;padding-left:20px;">';

      maps.forEach(map => {
        const rounds = grouped[date][map];

        // Находим последний раунд карты (с максимальным счетом)
        let lastRound = rounds[0];
        rounds.forEach(round => {
          const currentTotal = (round.ScoreCT || 0) + (round.ScoreT || 0);
          const lastTotal = (lastRound.ScoreCT || 0) + (lastRound.ScoreT || 0);
          if(currentTotal > lastTotal) {
            lastRound = round;
          }
        });

        // Получаем итоговый счет карты в системе Team 1/Team 2
        const finalScores = getTeamScores(lastRound);

        // Определяем победителя для цветового выделения
        const t1Won = finalScores.t1 > finalScores.t2;
        const t2Won = finalScores.t2 > finalScores.t1;

        // Заголовок карты
        html += '<div class="map-group" style="margin-bottom:20px;">';
        html += '<div class="map-header" data-map-id="' + date + '-' + map + '" style="padding:10px;background:var(--panel);cursor:pointer;border-radius:8px;display:flex;align-items:center;gap:12px;border-left:4px solid var(--accent);">';
        html += '<span class="expand-indicator" style="flex-shrink:0;">▶</span>';
        html += '<strong style="flex-shrink:0;">' + map + '</strong>';
        html += '<span style="font-size:13px;font-weight:600;background:var(--panel-2);padding:3px 8px;border-radius:5px;">';
        html += '<span style="color:' + (t1Won ? '#22c55e' : '#c8c8c8') + ';">[Team 1] ' + finalScores.t1 + '</span>';
        html += ' : ';
        html += '<span style="color:' + (t2Won ? '#22c55e' : '#c8c8c8') + ';">' + finalScores.t2 + ' [Team 2]</span>';
        html += '</span>';
        html += '<span style="font-size:11px;color:var(--muted);margin-left:auto;">раундов: ' + rounds.length + '</span>';
        html += '</div>';

        // Контент карты (раунды) - изначально скрыт
        html += '<div id="map-body-' + date + '-' + map + '" style="display:none;margin-top:10px;padding-left:20px;">';

        rounds.forEach((round, idx) => {
          const roundId = date + '-' + map + '-' + idx;
          html += renderRound(round, roundId);
        });

        html += '</div>'; // map-body
        html += '</div>'; // map-group
      });

      html += '</div>'; // date-body
      html += '</div>'; // date-group
    });

    roundsCount.textContent = 'Раундов: ' + totalRounds;

    if(totalRounds === 0) {
      roundsList.innerHTML = '<div class="small" style="padding:20px;text-align:center;color:var(--muted)">Нет данных</div>';
      return;
    }

    roundsList.innerHTML = html;

    // Добавляем обработчики для раскрытия/скрытия дат
    document.querySelectorAll('.date-header').forEach(header => {
      header.addEventListener('click', function() {
        const dateId = this.dataset.dateId;
        const body = document.getElementById('date-body-' + dateId);
        const indicator = this.querySelector('.expand-indicator');

        if(body.style.display === 'none') {
          body.style.display = 'block';
          indicator.textContent = '▼';
        } else {
          body.style.display = 'none';
          indicator.textContent = '▶';
        }
      });
    });

    // Добавляем обработчики для раскрытия/скрытия карт
    document.querySelectorAll('.map-header').forEach(header => {
      header.addEventListener('click', function() {
        const mapId = this.dataset.mapId;
        const body = document.getElementById('map-body-' + mapId);
        const indicator = this.querySelector('.expand-indicator');

        if(body.style.display === 'none') {
          body.style.display = 'block';
          indicator.textContent = '▼';
        } else {
          body.style.display = 'none';
          indicator.textContent = '▶';
        }
      });
    });

    // Добавляем обработчики для раскрытия/скрытия раундов
    document.querySelectorAll('.round-header').forEach(header => {
      header.addEventListener('click', function() {
        const roundId = this.dataset.roundId;
        const body = document.getElementById('round-body-' + roundId);
        const indicator = this.querySelector('.expand-indicator');

        if(body.style.display === 'none') {
          body.style.display = 'block';
          indicator.textContent = '▼';
        } else {
          body.style.display = 'none';
          indicator.textContent = '▶';
        }
      });
    });

    // Добавляем обработчики для сортировки таблиц
    setupTableSorting();
  }

  // Настройка сортировки таблиц
  function setupTableSorting() {
    document.querySelectorAll('.sortable').forEach(header => {
      header.addEventListener('click', function(e) {
        e.stopPropagation(); // Предотвращаем закрытие раунда

        const table = this.closest('table');
        const tbody = table.querySelector('tbody');
        const rows = Array.from(tbody.querySelectorAll('tr'));
        const sortType = this.dataset.sort;
        const headers = table.querySelectorAll('.sortable');
        const columnIndex = Array.from(this.parentNode.children).indexOf(this);

        // Определяем текущее направление сортировки
        const currentIndicator = this.querySelector('.sort-indicator');
        const isAsc = currentIndicator.textContent === '▲';

        // Очищаем все индикаторы
        headers.forEach(h => h.querySelector('.sort-indicator').textContent = '');

        // Устанавливаем новый индикатор
        currentIndicator.textContent = isAsc ? '▼' : '▲';

        // Сортируем строки
        rows.sort((a, b) => {
          const aValue = a.children[columnIndex].dataset.value;
          const bValue = b.children[columnIndex].dataset.value;

          let comparison = 0;
          if(sortType === 'name') {
            // Текстовая сортировка
            comparison = aValue.localeCompare(bValue);
          } else {
            // Числовая сортировка
            comparison = parseFloat(aValue) - parseFloat(bValue);
          }

          return isAsc ? comparison : -comparison;
        });

        // Перестраиваем таблицу
        rows.forEach(row => tbody.appendChild(row));
      });
    });
  }

  // Отрисовка одного раунда
  function renderRound(round, roundId) {
    const ctPlayers = round.Players.filter(p => p.Team === 3);
    const tPlayers = round.Players.filter(p => p.Team === 2);

    // Определяем какая команда (CT или T) является Team 1 в этом раунде
    const t1Team = getT1Team(round.Date);
    const ctIsT1 = t1Team && ctPlayers.some(p => t1Team.has(p.AccountID));

    // Разделяем игроков на Team 1 и Team 2
    const t1Players = ctIsT1 ? ctPlayers : tPlayers;
    const t2Players = ctIsT1 ? tPlayers : ctPlayers;

    // Получаем счет в системе Team 1/Team 2
    const scores = getTeamScores(round);

    let html = '<div class="round-item" style="margin-bottom:15px;border:1px solid var(--grid);border-radius:8px;overflow:hidden;">';

    // Заголовок раунда
    html += '<div class="round-header" data-round-id="' + roundId + '" style="padding:10px;background:var(--panel);cursor:pointer;display:flex;align-items:center;justify-content:space-between;">';
    html += '<div>';
    html += '<span class="expand-indicator" style="margin-right:10px;color:var(--accent);">▶</span>';
    html += '<strong>Раунд ' + round.RoundNumber + '</strong>';
    html += ' — Счёт: <span style="color:#3b82f6;">[Team 1] ' + scores.t1 + '</span> : <span style="color:#ef4444;">[Team 2] ' + scores.t2 + '</span>';
    html += '</div>';
    html += '<div class="small">Игроков: ' + round.Players.length + '</div>';
    html += '</div>';

    // Тело раунда (изначально скрыто)
    html += '<div id="round-body-' + roundId + '" style="display:none;padding:16px;background:var(--panel-2);">';

    // Определяем текущую сторону для каждой команды
    const t1Side = ctIsT1 ? 'CT' : 'T';
    const t2Side = ctIsT1 ? 'T' : 'CT';

    // Team 1
    if(t1Players.length > 0) {
      html += '<div style="margin-bottom:20px;">';
      html += '<h4 style="color:#3b82f6;margin:0 0 10px;font-size:14px;">[Team 1] (играет за ' + t1Side + ') — ' + t1Players.length + ' игроков</h4>';
      html += renderTeamTable(t1Players, playerNames);
      html += '</div>';
    }

    // Team 2
    if(t2Players.length > 0) {
      html += '<div>';
      html += '<h4 style="color:#ef4444;margin:0 0 10px;font-size:14px;">[Team 2] (играет за ' + t2Side + ') — ' + t2Players.length + ' игроков</h4>';
      html += renderTeamTable(t2Players, playerNames);
      html += '</div>';
    }

    html += '</div>';
    html += '</div>';

    return html;
  }

  // Отрисовка таблицы команды
  function renderTeamTable(players, playerNames) {
    // Сортируем игроков по рейтингу (по убыванию)
    const sortedPlayers = [...players].sort((a, b) => (b.Rating || 0) - (a.Rating || 0));

    let html = '<div style="overflow-x:auto;">';
    html += '<table class="sortable-table" style="width:100%%;border-collapse:collapse;font-size:12px;">';
    html += '<thead><tr style="background:var(--sticky);text-align:left;">';
    html += '<th class="sortable" data-sort="name" style="padding:6px;cursor:pointer;">Игрок <span class="sort-indicator"></span></th>';
    html += '<th class="sortable" data-sort="rating" style="padding:6px;text-align:center;cursor:pointer;">Рейтинг <span class="sort-indicator">▼</span></th>';
    html += '<th class="sortable" data-sort="damage" style="padding:6px;text-align:center;cursor:pointer;">Урон <span class="sort-indicator"></span></th>';
    html += '<th class="sortable" data-sort="kills" style="padding:6px;text-align:center;cursor:pointer;">K <span class="sort-indicator"></span></th>';
    html += '<th class="sortable" data-sort="deaths" style="padding:6px;text-align:center;cursor:pointer;">D <span class="sort-indicator"></span></th>';
    html += '<th class="sortable" data-sort="assists" style="padding:6px;text-align:center;cursor:pointer;">A <span class="sort-indicator"></span></th>';
    html += '<th class="sortable" data-sort="hsp" style="padding:6px;text-align:center;cursor:pointer;">HS%% <span class="sort-indicator"></span></th>';
    html += '<th class="sortable" data-sort="mvp" style="padding:6px;text-align:center;cursor:pointer;">MVP <span class="sort-indicator"></span></th>';
    html += '<th class="sortable" data-sort="money" style="padding:6px;text-align:center;cursor:pointer;">$ <span class="sort-indicator"></span></th>';
    html += '</tr></thead>';
    html += '<tbody>';

    sortedPlayers.forEach(p => {
      const name = playerNames[p.AccountID] || ('Player_' + p.AccountID);
      const rating = (p.Rating || 0).toFixed(2);
      html += '<tr style="border-bottom:1px solid var(--grid);">';
      html += '<td style="padding:6px;" data-value="' + name + '">' + name + '</td>';
      html += '<td style="padding:6px;text-align:center;font-weight:bold;color:var(--accent);" data-value="' + rating + '">' + rating + '</td>';
      html += '<td style="padding:6px;text-align:center;" data-value="' + (p.Damage || 0) + '">' + (p.Damage || 0) + '</td>';
      html += '<td style="padding:6px;text-align:center;" data-value="' + p.Kills + '">' + p.Kills + '</td>';
      html += '<td style="padding:6px;text-align:center;" data-value="' + p.Deaths + '">' + p.Deaths + '</td>';
      html += '<td style="padding:6px;text-align:center;" data-value="' + p.Assists + '">' + p.Assists + '</td>';
      html += '<td style="padding:6px;text-align:center;" data-value="' + p.HSP.toFixed(2) + '">' + p.HSP.toFixed(0) + '%%</td>';
      html += '<td style="padding:6px;text-align:center;" data-value="' + p.MVP + '">' + p.MVP + '</td>';
      html += '<td style="padding:6px;text-align:center;" data-value="' + p.Money + '">$' + p.Money + '</td>';
      html += '</tr>';
    });

    html += '</tbody></table>';
    html += '</div>';
    return html;
  }

  // Начальная отрисовка
  renderRounds();

  // Обработчик изменения фильтра дат
  window.addEventListener('dateFilterChanged', function() {
    renderRounds();
  });
})();
`, string(jRounds), string(jPlayerMap))
}
