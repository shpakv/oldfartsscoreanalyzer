package components

import (
	"encoding/json"
	"fmt"

	"oldfartscounter/internal/stats"
)

// KillsTabComponent отвечает за таб "Сорян, братан"
type KillsTabComponent struct{}

// NewKillsTab создает новый компонент таба убийств
func NewKillsTab() *KillsTabComponent {
	return &KillsTabComponent{}
}

// GenerateHTML генерирует HTML для таба убийств
func (k *KillsTabComponent) GenerateHTML(data *stats.StatsData) string {

	return fmt.Sprintf(`
<!-- KILLS -->
<div id="tab-kills" class="view">
  <!-- История "Сорян, Братан" -->
  <div style="margin-bottom:30px;">
    <h3 style="color:var(--accent);font-size:20px;margin-bottom:16px;">🎯 "Сорян, Братан" — История целей</h3>

    <!-- Текущие цели ноября -->
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-bottom:20px;">
      <!-- Основная цель (золото) -->
      <div style="background:linear-gradient(135deg, rgba(207,181,59,0.15) 0%%, rgba(207,181,59,0.05) 100%%);border:2px solid #cfb53b;border-radius:12px;padding:24px;">
        <div style="font-size:24px;font-weight:bold;color:#cfb53b;">🥇 maslina420</div>
        <div style="font-size:12px;color:var(--muted);margin-top:4px;">Ноябрь 2025 — Основная цель</div>
      </div>

      <!-- Специальная цель (серебро) -->
      <div style="background:linear-gradient(135deg, rgba(192,192,192,0.15) 0%%, rgba(192,192,192,0.05) 100%%);border:2px solid #c0c0c0;border-radius:12px;padding:24px;">
        <div style="font-size:24px;font-weight:bold;color:#c0c0c0;">🥈 Ai, Bobo!</div>
        <div style="font-size:12px;color:var(--muted);margin-top:4px;">Ноябрь 2025 — Спец. цель</div>
      </div>
    </div>

    <!-- История прошлых целей -->
    <details style="margin-top:16px;">
      <summary style="cursor:pointer;padding:12px;background:var(--panel);border-radius:8px;color:var(--accent);font-weight:bold;user-select:none;">📜 История прошлых целей</summary>
      <div style="margin-top:12px;padding-left:12px;">

        <!-- Октябрь 2025 -->
        <div style="border-left:3px solid #4b69ff;padding:12px 16px;margin-bottom:12px;background:rgba(75,105,255,0.05);border-radius:4px;">
          <div style="font-weight:bold;color:#4b69ff;font-size:16px;">Mr. Titspervert</div>
          <div style="font-size:13px;color:var(--muted);">Октябрь 2025 | 01.10 — 31.10</div>
        </div>

        <!-- Сентябрь 2025 -->
        <div style="border-left:3px solid #8847ff;padding:12px 16px;background:rgba(136,71,255,0.05);border-radius:4px;">
          <div style="font-weight:bold;color:#8847ff;font-size:16px;">Баба Валя</div>
          <div style="font-size:13px;color:var(--muted);">Сентябрь 2025 | 05.09 — 30.09</div>
        </div>

      </div>
    </details>
  </div>

  <div class="toolbar">
    <input id="qKills" type="search" placeholder="Поиск по именам…">
    <button class="btn btn-sort-max">↗ Топ пересечение</button>
    <label class="small"><input id="heatKills" type="checkbox" checked> Heatmap</label>
    <div class="legend"><div class="swatch"></div><span class="small">0 → %d</span></div>
  </div>
  <div class="table-wrap"><table id="gridKills"><thead></thead><tbody></tbody></table></div>
  <div class="small" style="margin-top:6px">Подсказка: <strong>Клик на имя в заголовке (сверху)</strong> сортирует строки по этому столбцу (кто больше убил этого игрока). <strong>Клик на имя слева</strong> сортирует столбцы по этой строке (кого этот игрок больше убивал). <strong>Повторный клик</strong> меняет направление (↓/↑). Клик по левому углу — сброс сортировки. <strong>Золотая подсветка</strong> 🥇 — основная жертва "Сорян, Братан", <strong>серебряная подсветка</strong> 🥈 — специальная цель.</div>
</div>`,
		data.KillMatrix.Max)
}

// GenerateJS генерирует JavaScript для таба убийств
func (k *KillsTabComponent) GenerateJS(data *stats.StatsData) string {
	// Создаем map для быстрого поиска индекса игрока по имени или SteamID
	type PlayerMapping struct {
		Title string
		Key   string
	}
	playerMappings := make([]PlayerMapping, len(data.Players))
	for i, p := range data.Players {
		playerMappings[i] = PlayerMapping{Title: p.Title, Key: p.Key}
	}

	jPlayerMappings, _ := json.Marshal(playerMappings)
	jHighlightedPlayer, _ := json.Marshal(data.HighlightedPlayer)

	return fmt.Sprintf(`
// Init: Сорян, братан
window.killsTabState = (function() {
  const playerMappings = %s;
  const playerTitles = playerMappings.map(p => p.Title);
  const highlightedPlayer = %s;
  const secondaryTarget = "Ai, Bobo!"; // Специальная цель (серебро)

  // Создаем индекс: и по Title, и по Key могут искать один и тот же индекс
  const playerIndexMap = {};
  playerMappings.forEach((p, idx) => {
    playerIndexMap[p.Title] = idx;
    playerIndexMap[p.Key] = idx;
  });

  function recalcKillMatrix(events) {
    const matrix = Array(playerMappings.length).fill(0).map(() => Array(playerMappings.length).fill(0));
    let maxKills = 0;

    events.forEach(e => {
      // Ищем индекс убийцы: сначала по имени, затем по SteamID
      let kIdx = playerIndexMap[e.KillerName];
      if (kIdx === undefined) kIdx = playerIndexMap[e.KillerSID];

      // Ищем индекс жертвы: сначала по имени, затем по SteamID
      let vIdx = playerIndexMap[e.VictimName];
      if (vIdx === undefined) vIdx = playerIndexMap[e.VictimSID];

      if (kIdx !== undefined && vIdx !== undefined) {
        matrix[kIdx][vIdx]++;
        if (matrix[kIdx][vIdx] > maxKills) maxKills = matrix[kIdx][vIdx];
      }
    });

    return { matrix, maxKills: maxKills || 1 };
  }

  function renderKillsTab() {
    const { matrix, maxKills } = recalcKillMatrix(window.filteredKillEvents || []);
    const legendEl = document.querySelector('#gridKills .legend .small');
    if (legendEl) legendEl.textContent = '0 → ' + maxKills;

    renderMatrix({
      rootId:"#gridKills",
      rowLabels: playerTitles,
      colLabels: playerTitles,
      data: matrix,
      maxVal: maxKills,
      qInputId: "qKills",
      csvBtnId: "csvKills",
      heatToggleId: "heatKills",
      cornerTitle: "Сорян, братан — Убийцы ↓ / Жертвы →",
      highlightedPlayer: highlightedPlayer,
      secondaryTarget: secondaryTarget
    });
  }

  // Переотрисовка при изменении фильтра дат
  window.addEventListener('dateFilterChanged', renderKillsTab);

  return { render: renderKillsTab };
})();

// Начальная отрисовка
window.killsTabState.render();`,
		string(jPlayerMappings),
		string(jHighlightedPlayer))
}
