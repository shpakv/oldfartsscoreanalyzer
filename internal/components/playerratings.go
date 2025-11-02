package components

import (
	"encoding/json"
	"fmt"

	"oldfartscounter/internal/stats"
)

// PlayerRatingsTabComponent отвечает за таб "Рейтинг игроков"
type PlayerRatingsTabComponent struct{}

// NewPlayerRatingsTab создает новый компонент таба рейтингов
func NewPlayerRatingsTab() *PlayerRatingsTabComponent {
	return &PlayerRatingsTabComponent{}
}

// GenerateHTML генерирует HTML для таба рейтингов
func (r *PlayerRatingsTabComponent) GenerateHTML() string {
	return `
<!-- PLAYER RATINGS -->
<div id="tab-player-ratings" class="view">
  <!-- Warning Popup -->
  <div id="ratingsWarning" style="position:fixed;top:0;left:0;width:100%;height:100%;background:rgba(0,0,0,0.95);z-index:9999;display:flex;align-items:center;justify-content:center;">
    <div style="background:#1a1a1a;border:2px solid #f59e0b;border-radius:16px;padding:40px;max-width:500px;text-align:center;box-shadow:0 20px 60px rgba(0,0,0,0.8);">
      <div style="font-size:64px;margin-bottom:20px;">⚠️</div>
      <h2 style="color:#f59e0b;margin:0 0 20px;font-size:24px;font-weight:bold;">ВНИМАНИЕ!</h2>
      <p style="color:var(--text);font-size:16px;line-height:1.8;margin-bottom:30px;">
        Следующий раздел содержит <strong style="color:#ef4444;">объективную математическую оценку</strong> вашей игры.<br><br>
        Эти данные могут <strong>серьёзно навредить вашему эго</strong> и вызвать:<br>
        • Осознание своей бездарности<br>
        • Желание обвинить teammates<br>
        • Непреодолимую тягу к оправданиям<br><br>
        <span style="color:var(--muted);font-size:14px;">Только для лиц старше 18 лет и с крепкими нервами.</span>
      </p>
      <div style="display:flex;gap:12px;justify-content:center;">
        <button onclick="acceptRatings()" onmouseover="this.style.transform='scale(1.05)'" onmouseout="this.style.transform='scale(1)'" style="background:#22c55e;color:white;border:none;padding:14px 32px;border-radius:8px;font-size:16px;font-weight:bold;cursor:pointer;transition:all 0.2s;">
          Да, я большой мальчик 💪
        </button>
        <button onclick="rejectRatings()" onmouseover="this.style.transform='scale(1.05)'" onmouseout="this.style.transform='scale(1)'" style="background:#ef4444;color:white;border:none;padding:14px 32px;border-radius:8px;font-size:16px;font-weight:bold;cursor:pointer;transition:all 0.2s;">
          Нет, я плакса 😭
        </button>
      </div>
    </div>
  </div>

  <div class="toolbar">
    <span class="small" style="margin-left: auto;" id="ratingsCount">Игроков: 0</span>
  </div>

  <div style="margin-top: 20px;">
    <div style="background:var(--panel);padding:20px;border-radius:8px;margin-bottom:20px;">
      <div style="display:flex;align-items:center;justify-content:space-between;cursor:pointer;user-select:none;" id="ratingsInfoToggle">
        <h3 style="margin:0;font-size:18px;color:var(--accent);">📚 О системе рейтингов</h3>
        <span id="ratingsInfoArrow" style="font-size:20px;transition:transform 0.3s;color:var(--accent);">▼</span>
      </div>
      <p style="margin:10px 0 0 0;color:var(--muted);font-size:13px;line-height:1.7;">
        Справедливая оценка мастерства с учетом урона, контекста боя и стабильности игры.
        <span style="color:var(--accent);cursor:pointer;" onclick="document.getElementById('ratingsInfoToggle').click()">Подробнее...</span>
      </p>

      <div id="ratingsInfoDetails" style="display:none;margin-top:15px;border-top:1px solid var(--grid);padding-top:15px;">
        <div style="margin-bottom:15px;">
          <h4 style="margin:0 0 8px;font-size:15px;color:#3b82f6;">Что это такое?</h4>
          <p style="margin:0;color:var(--muted);font-size:13px;line-height:1.7;">
            Это <strong>справедливая оценка мастерства</strong> каждого игрока, которая учитывает не только убийства,
            но и урон, контекст боя (3v5 или 5v3), ассисты и смерти. Система защищена от случайных удачных игр —
            чтобы попасть в топ, нужно <strong>стабильно</strong> хорошо играть.
          </p>
        </div>

      <div style="margin-bottom:15px;">
        <h4 style="margin:0 0 8px;font-size:15px;color:#3b82f6;">Почему не просто K/D?</h4>
        <p style="margin:0;color:var(--muted);font-size:13px;line-height:1.7;">
          K/D не учитывает множество важных факторов:
        </p>
        <ul style="margin:4px 0 0 20px;padding:0;color:var(--muted);font-size:13px;line-height:1.7;">
          <li>Вы нанесли 99 урона, а киллом добил союзник → K/D не видит ваш вклад</li>
          <li>Вы выиграли клатч 1v3 → K/D оценит это как обычные 3 килла</li>
          <li>У игрока 2 раунда с ace → K/D покажет его как "лучшего"</li>
          <li>Саппорт с кучей ассистов → K/D недооценивает его роль</li>
        </ul>
        <p style="margin:8px 0 0 0;color:var(--muted);font-size:13px;line-height:1.7;">
          <strong style="color:#10b981;">EPI решает эти проблемы</strong> — учитывает урон, ассисты, численное преимущество и штрафует за смерти.
        </p>
      </div>

      <div style="margin-bottom:15px;">
        <h4 style="margin:0 0 8px;font-size:15px;color:#3b82f6;">Зачем байесовский рейтинг?</h4>
        <p style="margin:0;color:var(--muted);font-size:13px;line-height:1.7;">
          <strong>Проблема:</strong> Игрок сыграл 3 раунда с ace → его средний EPI = 2.5 (невероятно круто!).
          Но это <strong>случайность</strong>, не стабильный уровень игры.
        </p>
        <p style="margin:8px 0 0 0;color:var(--muted);font-size:13px;line-height:1.7;">
          <strong style="color:#10b981;">Решение:</strong> Система "дофантазирует" за вас <strong><span id="minRoundsValue">100</span> виртуальных раундов</strong>
          со средним результатом (μ = <span id="averageMuValue3" style="color:#10b981;font-weight:bold;">0.6</span>). Пока вы не наиграли достаточно, система не поверит в ваш "гениальный" рейтинг.
        </p>
        <p style="margin:8px 0 0 0;color:var(--muted);font-size:13px;line-height:1.7;">
          <strong>Пример:</strong> У вас 10 раундов с EPI = 1.2 → байесовский рейтинг = 0.72 (система скептична).
          У вас 100 раундов с EPI = 1.2 → байесовский рейтинг = 1.09 (система вам верит).
        </p>
        <p style="margin:8px 0 0 0;color:#10b981;font-size:13px;line-height:1.7;">
          💡 <strong>Шкала ниже автоматически адаптируется</strong> под средний уровень вашей группы (μ = <span id="averageMuValue2" style="font-weight:bold;">0.6</span>)
        </p>
      </div>

      <div style="margin-bottom:15px;">
        <h4 style="margin:0 0 8px;font-size:15px;color:#3b82f6;">Формула (для гиков)</h4>
        <div style="background:var(--panel-2);padding:10px;border-radius:6px;border-left:3px solid var(--accent);">
          <code style="color:#10b981;font-size:14px;">Rating = (ΣE + K×μ) / (N + K)</code>
        </div>
        <p style="margin:8px 0 0 0;color:var(--muted);font-size:12px;line-height:1.6;">
          <strong>ΣE</strong> = сумма EPI по всем вашим раундам<br>
          <strong>N</strong> = количество сыгранных раундов<br>
          <strong>K</strong> = <span id="minRoundsValue">100</span> (константа: минимум раундов для достоверности)<br>
          <strong>μ</strong> = <span id="averageMuValue" style="color:#10b981;font-weight:bold;">0.6</span> (средний EPI всех игроков вашей группы)
        </p>
        <p style="margin:8px 0 0 0;color:var(--muted);font-size:12px;line-height:1.6;">
          💡 <strong>Наведите курсор на рейтинг</strong> в таблице, чтобы увидеть детальный расчет!
        </p>
      </div>

      <div>
        <h4 style="margin:0 0 8px;font-size:15px;color:#cfb53b;">Как читать рейтинг?</h4>
        <div style="display:flex;gap:12px;flex-wrap:wrap;margin-top:8px;" id="ratingScaleCards">
          <div style="background:rgba(156,163,175,0.15);border-left:3px solid #9ca3af;padding:8px 12px;border-radius:4px;flex:1;min-width:150px;">
            <div style="color:#9ca3af;font-weight:bold;font-size:12px;" id="scaleWeak">< 0.5 — Подпивас</div>
            <div style="color:var(--muted);font-size:11px;margin-top:2px;">Mil-Spec</div>
          </div>
          <div style="background:rgba(75,105,255,0.15);border-left:3px solid #4b69ff;padding:8px 12px;border-radius:4px;flex:1;min-width:150px;">
            <div style="color:#4b69ff;font-weight:bold;font-size:12px;" id="scaleAverage">0.5 - 0.6 — Пердун</div>
            <div style="color:var(--muted);font-size:11px;margin-top:2px;">Restricted</div>
          </div>
          <div style="background:rgba(239,68,68,0.15);border-left:3px solid #ef4444;padding:8px 12px;border-radius:4px;flex:1;min-width:150px;">
            <div style="color:#ef4444;font-weight:bold;font-size:12px;" id="scaleGood">0.6 - 0.7 — Ебака</div>
            <div style="color:var(--muted);font-size:11px;margin-top:2px;">Classified</div>
          </div>
          <div style="background:rgba(207,181,59,0.15);border-left:3px solid #cfb53b;padding:8px 12px;border-radius:4px;flex:1;min-width:150px;">
            <div style="color:#cfb53b;font-weight:bold;font-size:12px;" id="scaleMonster">≥ 0.7 — Гиперебака</div>
            <div style="color:var(--muted);font-size:11px;margin-top:2px;">Covert</div>
          </div>
        </div>
      </div>

        <div style="margin-top:15px;padding:10px;background:rgba(245,158,11,0.1);border-radius:6px;border:1px solid rgba(245,158,11,0.3);">
          <div style="color:#f59e0b;font-size:12px;font-weight:bold;margin-bottom:4px;">⚠️ Внимание</div>
          <div style="color:var(--muted);font-size:12px;line-height:1.6;">
            Если рядом с именем значок "⚠ недостаточно данных" — это значит игрок сыграл меньше 100 раундов.
            Его рейтинг еще <strong>не стабилизировался</strong> и может сильно измениться.
          </div>
        </div>
      </div>
    </div>

    <script>
      // Toggle для раскрытия информации о рейтинге
      document.getElementById('ratingsInfoToggle').addEventListener('click', function() {
        const details = document.getElementById('ratingsInfoDetails');
        const arrow = document.getElementById('ratingsInfoArrow');

        if (details.style.display === 'none') {
          details.style.display = 'block';
          arrow.style.transform = 'rotate(180deg)';
        } else {
          details.style.display = 'none';
          arrow.style.transform = 'rotate(0deg)';
        }
      });
    </script>

    <div id="ratingsTable"></div>
  </div>

  <!-- Custom Tooltip -->
  <div id="ratingTooltip" style="display:none;position:fixed;background:#1a1a1a;border:2px solid var(--accent);border-radius:8px;padding:12px;z-index:10000;pointer-events:none;font-family:monospace;font-size:12px;line-height:1.6;box-shadow:0 4px 12px rgba(0,0,0,0.5);max-width:400px;"></div>
</div>

<style>
.rating-cell {
  cursor: help;
  position: relative;
}
</style>`
}

// GenerateJS генерирует JavaScript для таба рейтингов
func (r *PlayerRatingsTabComponent) GenerateJS(data *stats.StatsData) string {
	jRatings, _ := json.Marshal(data.PlayerRatings)
	minRounds := data.MinRoundsForRating
	averageMu := data.AverageMu

	return fmt.Sprintf(`
// Init: Рейтинг игроков
window.playerRatingsTabState = (function() {
  const originalRatings = %s;
  const K = %v; // Минимальное количество раундов для достоверности
  const AVERAGE_MU = %v; // Средний EPI всех игроков (вычислено из реальных данных)

  const ratingsCount = document.getElementById('ratingsCount');
  const ratingsTable = document.getElementById('ratingsTable');
  const warningPopup = document.getElementById('ratingsWarning');
  const minRoundsValueEl = document.getElementById('minRoundsValue');
  const averageMuValueEl = document.getElementById('averageMuValue');

  // Устанавливаем значения K и μ в описании
  if (minRoundsValueEl) {
    minRoundsValueEl.textContent = K.toFixed(0);
  }
  if (averageMuValueEl) {
    averageMuValueEl.textContent = AVERAGE_MU.toFixed(3);
  }
  const averageMuValueEl2 = document.getElementById('averageMuValue2');
  if (averageMuValueEl2) {
    averageMuValueEl2.textContent = AVERAGE_MU.toFixed(3);
  }
  const averageMuValueEl3 = document.getElementById('averageMuValue3');
  if (averageMuValueEl3) {
    averageMuValueEl3.textContent = AVERAGE_MU.toFixed(3);
  }

  // Вычисляем адаптивные границы для шкалы на основе μ
  const THRESHOLDS = {
    weak: (AVERAGE_MU * 0.85).toFixed(2),      // < μ * 0.85 → Подпивас
    average: (AVERAGE_MU * 1.05).toFixed(2),   // μ * 0.85 - μ * 1.05 → Пердун, ≥ μ * 1.05 → Ебака
    monster: (AVERAGE_MU * 1.25).toFixed(2)    // ≥ μ * 1.25 → Гиперебака
  };

  // Обновляем шкалу "Как читать рейтинг?" с реальными значениями
  const scaleWeak = document.getElementById('scaleWeak');
  const scaleAverage = document.getElementById('scaleAverage');
  const scaleGood = document.getElementById('scaleGood');
  const scaleMonster = document.getElementById('scaleMonster');
  if (scaleWeak) scaleWeak.textContent = '< ' + THRESHOLDS.weak + ' — Подпивас';
  if (scaleAverage) scaleAverage.textContent = THRESHOLDS.weak + ' - ' + THRESHOLDS.average + ' — Пердун';
  if (scaleGood) scaleGood.textContent = THRESHOLDS.average + ' - ' + THRESHOLDS.monster + ' — Ебака';
  if (scaleMonster) scaleMonster.textContent = '≥ ' + THRESHOLDS.monster + ' — Гиперебака';

  // Проверяем, принял ли пользователь предупреждение ранее
  const hasAccepted = localStorage.getItem('ratingsWarningAccepted') === 'true';
  if (hasAccepted) {
    warningPopup.style.display = 'none';
  }

  // Функции для кнопок (доступны глобально)
  window.acceptRatings = function() {
    localStorage.setItem('ratingsWarningAccepted', 'true');
    warningPopup.style.display = 'none';
  };

  window.rejectRatings = function() {
    // Переключаемся на первый таб (убийства)
    const killsBtn = document.querySelector('.tab-btn[data-tab="kills"]');
    if (killsBtn) {
      killsBtn.click();
    }
    alert('Правильное решение! Берегите своё эго 😌');
  };

  // Функция для расчета рейтингов из раундов
  function calculateRatings(roundStats) {
    const playerData = {};

    // Агрегируем данные по игрокам
    roundStats.forEach(function(round) {
      round.Players.forEach(function(playerStats) {
        if (playerStats.AccountID === 0) return;

        if (!playerData[playerStats.AccountID]) {
          playerData[playerStats.AccountID] = {
            AccountID: playerStats.AccountID,
            Name: '', // Имя установим позже
            RoundsPlayed: 0,
            TotalEPI: 0,
            TotalDamage: 0,
            TotalKills: 0,
            TotalDeaths: 0,
            TotalAssists: 0,
            WinRounds: 0,
            LastPlayed: ''
          };
        }

        var rating = playerData[playerStats.AccountID];
        rating.RoundsPlayed++;
        rating.TotalEPI += playerStats.Rating;
        rating.TotalDamage += playerStats.Damage;
        rating.TotalKills += playerStats.Kills;
        rating.TotalDeaths += playerStats.Deaths;
        rating.TotalAssists += playerStats.Assists;

        // Проверяем победу
        if ((playerStats.Team === 3 && round.Winner === 3) || (playerStats.Team === 2 && round.Winner === 2)) {
          rating.WinRounds++;
        }

        // Обновляем последнюю дату игры
        if (round.Date && (!rating.LastPlayed || round.Date > rating.LastPlayed)) {
          rating.LastPlayed = round.Date;
        }
      });
    });

    // Находим имена игроков из оригинальных рейтингов
    originalRatings.forEach(function(orig) {
      if (playerData[orig.AccountID]) {
        playerData[orig.AccountID].Name = orig.Name;
      }
    });

    // Рассчитываем среднее EPI по всем игрокам для μ
    var totalEPI = 0;
    var totalRounds = 0;
    for (var accountID in playerData) {
      totalEPI += playerData[accountID].TotalEPI;
      totalRounds += playerData[accountID].RoundsPlayed;
    }

    var mu = totalRounds > 0 ? totalEPI / totalRounds : 0.6;

    // Рассчитываем финальные рейтинги
    var ratings = [];
    for (var accountID in playerData) {
      var rating = playerData[accountID];

      // Простое среднее
      rating.AverageEPI = rating.RoundsPlayed > 0 ? rating.TotalEPI / rating.RoundsPlayed : 0;

      // Байесовский рейтинг
      rating.BayesianEPI = (rating.TotalEPI + K * mu) / (rating.RoundsPlayed + K);

      // Если нет имени, используем AccountID
      if (!rating.Name) {
        rating.Name = 'Player_' + rating.AccountID;
      }

      ratings.push(rating);
    }

    // Сортируем по байесовскому рейтингу (по убыванию)
    ratings.sort(function(a, b) {
      return b.BayesianEPI - a.BayesianEPI;
    });

    return ratings;
  }

  function renderRatingsTable() {
    var ratings = calculateRatings(window.filteredRoundStats || []);

    if (ratings.length === 0) {
      ratingsTable.innerHTML = '<div class="small" style="padding:20px;text-align:center;color:var(--muted)">Нет данных</div>';
      ratingsCount.textContent = 'Игроков: 0';
      return;
    }

    ratingsCount.textContent = 'Игроков: ' + ratings.length;

    // Вычисляем μ (средний EPI всех игроков) для tooltip
    var totalEPI = 0;
    var totalRounds = 0;
    for (var i = 0; i < ratings.length; i++) {
      totalEPI += ratings[i].TotalEPI;
      totalRounds += ratings[i].RoundsPlayed;
    }
    var mu = totalRounds > 0 ? totalEPI / totalRounds : 0.6;

    let html = '<div style="overflow-x:auto;">';
    html += '<table style="width:100%%;border-collapse:collapse;">';
    html += '<thead><tr style="background:var(--sticky);text-align:left;">';
    html += '<th style="padding:12px;text-align:center;width:60px;">#</th>';
    html += '<th style="padding:12px;">Игрок</th>';
    html += '<th style="padding:12px;text-align:center;">Рейтинг</th>';
    html += '<th style="padding:12px;text-align:center;">Раундов</th>';
    html += '<th style="padding:12px;text-align:center;">K/D/A</th>';
    html += '<th style="padding:12px;text-align:center;">Урон</th>';
    html += '<th style="padding:12px;text-align:center;">Побед</th>';
    html += '<th style="padding:12px;text-align:center;">Win%%</th>';
    html += '<th style="padding:12px;text-align:center;">Последняя игра</th>';
    html += '</tr></thead>';
    html += '<tbody>';

    ratings.forEach((player, idx) => {
      const kd = player.TotalDeaths > 0 ? (player.TotalKills / player.TotalDeaths).toFixed(2) : player.TotalKills.toFixed(2);
      const winRate = player.RoundsPlayed > 0 ? ((player.WinRounds / player.RoundsPlayed) * 100).toFixed(1) : '0.0';

      // Определяем цвет рейтинга на основе адаптивных границ (цвета предметов CS2)
      let ratingColor = 'var(--text)';
      let playerStatus = '';
      if (player.BayesianEPI >= parseFloat(THRESHOLDS.monster)) {
        ratingColor = '#cfb53b'; // gold - Гиперебака (Covert)
        playerStatus = 'Гиперебака';
      } else if (player.BayesianEPI >= parseFloat(THRESHOLDS.average)) {
        ratingColor = '#ef4444'; // red - Ебака (Classified)
        playerStatus = 'Ебака';
      } else if (player.BayesianEPI >= parseFloat(THRESHOLDS.weak)) {
        ratingColor = '#4b69ff'; // blue - Пердун (Restricted)
        playerStatus = 'Пердун';
      } else {
        ratingColor = '#9ca3af'; // gray - Подпивас (Mil-Spec)
        playerStatus = 'Подпивас';
      }

      // Определяем достоверность (если раундов < K, выделяем)
      const isUnreliable = player.RoundsPlayed < K;
      const nameStyle = isUnreliable ? 'color:var(--muted);font-style:italic;' : '';
      const unreliableWarning = isUnreliable ? ' <span style="font-size:11px;color:#f59e0b;">⚠ недостаточно данных</span>' : '';

      // Формируем tooltip с детальной калькуляцией
      const numerator = (player.TotalEPI + K * mu).toFixed(3);
      const denominator = (player.RoundsPlayed + K).toFixed(0);
      const tooltipData = {
        totalEPI: player.TotalEPI.toFixed(3),
        rounds: player.RoundsPlayed,
        K: K.toFixed(0),
        mu: mu.toFixed(3),
        numerator: numerator,
        denominator: denominator,
        result: player.BayesianEPI.toFixed(3),
        status: playerStatus,
        color: ratingColor
      };

      html += '<tr style="border-bottom:1px solid var(--grid);">';
      html += '<td style="padding:12px;text-align:center;color:var(--muted);font-weight:bold;">' + (idx + 1) + '</td>';
      html += '<td style="padding:12px;' + nameStyle + '">' + player.Name + unreliableWarning + '</td>';
      html += '<td class="rating-cell" style="padding:12px;text-align:center;font-weight:bold;font-size:16px;color:' + ratingColor + ';" data-tooltip="' + encodeURIComponent(JSON.stringify(tooltipData)) + '">' + player.BayesianEPI.toFixed(3) + '</td>';
      html += '<td style="padding:12px;text-align:center;">' + player.RoundsPlayed + '</td>';
      html += '<td style="padding:12px;text-align:center;">' + player.TotalKills + ' / ' + player.TotalDeaths + ' / ' + player.TotalAssists + '</td>';
      html += '<td style="padding:12px;text-align:center;">' + player.TotalDamage + '</td>';
      html += '<td style="padding:12px;text-align:center;">' + player.WinRounds + '</td>';
      html += '<td style="padding:12px;text-align:center;">' + winRate + '%%</td>';
      html += '<td style="padding:12px;text-align:center;color:var(--muted);font-size:13px;">' + (player.LastPlayed || '-') + '</td>';
      html += '</tr>';
    });

    html += '</tbody></table>';
    html += '</div>';

    ratingsTable.innerHTML = html;

    // Подключаем обработчики tooltip для всех ячеек с рейтингом
    const ratingCells = ratingsTable.querySelectorAll('.rating-cell');
    const tooltip = document.getElementById('ratingTooltip');

    ratingCells.forEach(function(cell) {
      cell.addEventListener('mouseenter', function(e) {
        const data = JSON.parse(decodeURIComponent(cell.dataset.tooltip));

        const tooltipHTML = '<div style="color:' + data.color + ';font-weight:bold;margin-bottom:8px;font-size:16px;">📊 ' + data.status + '</div>' +
                           '<div style="border-bottom:1px solid #444;margin-bottom:8px;"></div>' +
                           '<div style="color:#3b82f6;">ΣE (сумма EPI):</div>' +
                           '<div style="margin-left:12px;margin-bottom:4px;">' + data.totalEPI + '</div>' +
                           '<div style="color:#3b82f6;">N (раундов):</div>' +
                           '<div style="margin-left:12px;margin-bottom:4px;">' + data.rounds + '</div>' +
                           '<div style="color:#3b82f6;">K (константа):</div>' +
                           '<div style="margin-left:12px;margin-bottom:4px;">' + data.K + '</div>' +
                           '<div style="color:#3b82f6;">μ (средний EPI):</div>' +
                           '<div style="margin-left:12px;margin-bottom:8px;">' + data.mu + '</div>' +
                           '<div style="border-bottom:1px solid #444;margin-bottom:8px;"></div>' +
                           '<div style="color:#10b981;font-weight:bold;">Формула:</div>' +
                           '<div style="margin-left:12px;margin-top:4px;color:var(--muted);">Rating = (ΣE + K×μ) / (N + K)</div>' +
                           '<div style="margin-left:12px;margin-top:4px;">Rating = (' + data.totalEPI + ' + ' + data.K + '×' + data.mu + ') / (' + data.rounds + ' + ' + data.K + ')</div>' +
                           '<div style="margin-left:12px;margin-top:4px;">Rating = ' + data.numerator + ' / ' + data.denominator + '</div>' +
                           '<div style="margin-left:12px;margin-top:8px;color:#f59e0b;font-weight:bold;font-size:14px;">Rating = ' + data.result + '</div>';

        tooltip.innerHTML = tooltipHTML;
        tooltip.style.display = 'block';

        // Позиционируем tooltip рядом с курсором
        const updatePosition = function(event) {
          tooltip.style.left = (event.clientX + 15) + 'px';
          tooltip.style.top = (event.clientY + 15) + 'px';
        };
        updatePosition(e);
        cell.addEventListener('mousemove', updatePosition);
        cell._updatePosition = updatePosition;
      });

      cell.addEventListener('mouseleave', function() {
        tooltip.style.display = 'none';
        if (cell._updatePosition) {
          cell.removeEventListener('mousemove', cell._updatePosition);
          delete cell._updatePosition;
        }
      });
    });
  }

  // Подписываемся на изменение фильтра дат
  window.addEventListener('dateFilterChanged', renderRatingsTable);

  return { render: renderRatingsTable };
})();

// Начальная отрисовка
window.playerRatingsTabState.render();
`, string(jRatings), minRounds, averageMu)
}
