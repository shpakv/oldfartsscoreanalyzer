package components

import (
	"encoding/json"
	"fmt"

	"oldfartscounter/internal/stats"
)

// WeaponsTabComponent отвечает за табы с оружием
type WeaponsTabComponent struct{}

// NewWeaponsTab создает новый компонент таба оружия
func NewWeaponsTab() *WeaponsTabComponent {
	return &WeaponsTabComponent{}
}

// GenerateKillerWeaponHTML генерирует HTML для таба "Кто с чего убивает"
func (w *WeaponsTabComponent) GenerateKillerWeaponHTML(data *stats.StatsData) string {
	return fmt.Sprintf(`
<!-- KILLER x WEAPON (rows=weapons, cols=players) -->
<div id="tab-kw" class="view">
  <div class="toolbar">
    <input id="qKW" type="search" placeholder="Поиск по именам/оружию…">
    <button class="btn btn-sort-max">↗ Топ пересечение</button>
    <label class="small"><input id="heatKW" type="checkbox" checked> Heatmap</label>
    <div class="legend"><div class="swatch"></div><span class="small">0 → %d</span></div>
  </div>
  <div class="table-wrap"><table id="gridKW"><thead></thead><tbody></tbody></table></div>
</div>`,
		data.WeaponData.KillerMax)
}

// GenerateVictimWeaponHTML генерирует HTML для таба "Кого чем убивают"
func (w *WeaponsTabComponent) GenerateVictimWeaponHTML(data *stats.StatsData) string {
	return fmt.Sprintf(`
<!-- VICTIM x WEAPON (rows=players-victims, cols=weapons) -->
<div id="tab-vw" class="view">
  <div class="toolbar">
    <input id="qVW" type="search" placeholder="Поиск по именам/оружию…">
    <button class="btn btn-sort-max">↗ Топ пересечение</button>
    <label class="small"><input id="heatVW" type="checkbox" checked> Heatmap</label>
    <div class="legend"><div class="swatch"></div><span class="small">0 → %d</span></div>
  </div>
  <div class="table-wrap"><table id="gridVW"><thead></thead><tbody></tbody></table></div>
</div>`,
		data.WeaponData.VictimMax)
}

// GenerateKillerWeaponJS генерирует JavaScript для таба "Кто с чего убивает"
func (w *WeaponsTabComponent) GenerateKillerWeaponJS(data *stats.StatsData) string {
	type PlayerMapping struct {
		Title string
		Key   string
	}
	playerMappings := make([]PlayerMapping, len(data.Players))
	for i, p := range data.Players {
		playerMappings[i] = PlayerMapping{Title: p.Title, Key: p.Key}
	}

	jPlayerMappings, _ := json.Marshal(playerMappings)

	return fmt.Sprintf(`
// Init: Кто с чего убивает (строки=оружие, столбцы=игроки)
window.kwTabState = (function() {
  const playerMappings = %s;
  const playerTitles = playerMappings.map(p => p.Title);
  const playerIndexMap = {};
  playerMappings.forEach((p, idx) => {
    playerIndexMap[p.Title] = idx;
    playerIndexMap[p.Key] = idx;
  });

  function recalcKWMatrix(events) {
    const weaponSet = new Set();
    const matrix = {};

    events.forEach(e => {
      if (!e.Weapon) return;
      weaponSet.add(e.Weapon);

      let kIdx = playerIndexMap[e.KillerName];
      if (kIdx === undefined) kIdx = playerIndexMap[e.KillerSID];

      if (kIdx !== undefined) {
        if (!matrix[e.Weapon]) matrix[e.Weapon] = Array(playerMappings.length).fill(0);
        matrix[e.Weapon][kIdx]++;
      }
    });

    const weapons = weaponSet.size > 0 ? Array.from(weaponSet).sort() : ["(нет данных)"];
    const matrixArray = weapons.map(w => matrix[w] || Array(playerMappings.length).fill(0));

    let maxVal = 0;
    matrixArray.forEach(row => {
      row.forEach(v => { if (v > maxVal) maxVal = v; });
    });

    return { weapons, matrix: matrixArray, maxVal: maxVal || 1 };
  }

  function renderKWTab() {
    const { weapons, matrix, maxVal } = recalcKWMatrix(window.filteredKillEvents || []);
    const legendEl = document.querySelector('#gridKW .legend .small');
    if (legendEl) legendEl.textContent = '0 → ' + maxVal;

    renderMatrix({
      rootId:"#gridKW",
      rowLabels: weapons,
      colLabels: playerTitles,
      data: matrix,
      maxVal: maxVal,
      qInputId: "qKW",
      csvBtnId: "csvKW",
      heatToggleId: "heatKW",
      cornerTitle: "Оружие ↓ / Убийцы →"
    });
  }

  window.addEventListener('dateFilterChanged', renderKWTab);
  return { render: renderKWTab };
})();

window.kwTabState.render();`,
		string(jPlayerMappings))
}

// GenerateVictimWeaponJS генерирует JavaScript для таба "Кого чем убивают"
func (w *WeaponsTabComponent) GenerateVictimWeaponJS(data *stats.StatsData) string {
	type PlayerMapping struct {
		Title string
		Key   string
	}
	playerMappings := make([]PlayerMapping, len(data.Players))
	for i, p := range data.Players {
		playerMappings[i] = PlayerMapping{Title: p.Title, Key: p.Key}
	}

	jPlayerMappings, _ := json.Marshal(playerMappings)

	return fmt.Sprintf(`
// Init: Кого чем убивают (строки=жертвы, столбцы=оружие)
window.vwTabState = (function() {
  const playerMappings = %s;
  const playerTitles = playerMappings.map(p => p.Title);
  const playerIndexMap = {};
  playerMappings.forEach((p, idx) => {
    playerIndexMap[p.Title] = idx;
    playerIndexMap[p.Key] = idx;
  });

  function recalcVWMatrix(events) {
    const weaponSet = new Set();
    const matrix = Array(playerMappings.length).fill(0).map(() => ({}));

    events.forEach(e => {
      if (!e.Weapon) return;
      weaponSet.add(e.Weapon);

      let vIdx = playerIndexMap[e.VictimName];
      if (vIdx === undefined) vIdx = playerIndexMap[e.VictimSID];

      if (vIdx !== undefined) {
        if (!matrix[vIdx][e.Weapon]) matrix[vIdx][e.Weapon] = 0;
        matrix[vIdx][e.Weapon]++;
      }
    });

    const weapons = weaponSet.size > 0 ? Array.from(weaponSet).sort() : ["(нет данных)"];
    const matrixArray = matrix.map(row => weapons.map(w => row[w] || 0));

    let maxVal = 0;
    matrixArray.forEach(row => {
      row.forEach(v => { if (v > maxVal) maxVal = v; });
    });

    return { weapons, matrix: matrixArray, maxVal: maxVal || 1 };
  }

  function renderVWTab() {
    const { weapons, matrix, maxVal } = recalcVWMatrix(window.filteredKillEvents || []);
    const legendEl = document.querySelector('#gridVW .legend .small');
    if (legendEl) legendEl.textContent = '0 → ' + maxVal;

    renderMatrix({
      rootId:"#gridVW",
      rowLabels: playerTitles,
      colLabels: weapons,
      data: matrix,
      maxVal: maxVal,
      qInputId: "qVW",
      csvBtnId: "csvVW",
      heatToggleId: "heatVW",
      cornerTitle: "Жертвы ↓ / Оружие →"
    });
  }

  window.addEventListener('dateFilterChanged', renderVWTab);
  return { render: renderVWTab };
})();

window.vwTabState.render();`,
		string(jPlayerMappings))
}
