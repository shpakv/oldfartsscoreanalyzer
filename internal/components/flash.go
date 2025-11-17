package components

import (
	"encoding/json"
	"fmt"

	"oldfartscounter/internal/stats"
)

// FlashTabComponent отвечает за таб "Индекс Пирога"
type FlashTabComponent struct{}

// NewFlashTab создает новый компонент таба флешек
func NewFlashTab() *FlashTabComponent {
	return &FlashTabComponent{}
}

// GenerateHTML генерирует HTML для таба флешек
func (f *FlashTabComponent) GenerateHTML(data *stats.StatsData) string {
	return `
<!-- FLASH (Pirog Index) -->
<div id="tab-flash" class="view">
  <div class="toolbar">
    <input id="qFlash" type="search" placeholder="Поиск по именам…">
    <button class="btn btn-sort-max">↗ Топ пересечение</button>
    <label class="small"><input id="heatFlash" type="checkbox" checked> Heatmap</label>
    <div class="legend"><div class="swatch"></div><span class="small">0 → макс</span></div>
  </div>
  <div class="table-wrap"><table id="gridFlash"><thead></thead><tbody></tbody></table></div>
  <div class="small" style="margin-top:6px">Индекс Пирога: сколько секунд кто кого слепил. Клик по столбцам сортирует строки, клик по строкам сортирует столбцы.</div>
</div>`
}

// GenerateJS генерирует JavaScript для таба флешек
func (f *FlashTabComponent) GenerateJS(data *stats.StatsData) string {
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
// Init: Индекс Пирога (только секунды)
window.flashTabState = (function() {
  const playerMappings = %s;
  const playerTitles = playerMappings.map(p => p.Title);
  const playerIndexMap = {};
  playerMappings.forEach((p, idx) => {
    playerIndexMap[p.Title] = idx;
    playerIndexMap[p.Key] = idx;
  });

  function recalcFlashMatrix(events) {
    const secondsMatrix = Array(playerMappings.length).fill(0).map(() => Array(playerMappings.length).fill(0));
    let secondsMax = 0;

    events.forEach(e => {
      let fIdx = playerIndexMap[e.FlasherName];
      if (fIdx === undefined) fIdx = playerIndexMap[e.FlasherSID];

      let vIdx = playerIndexMap[e.VictimName];
      if (vIdx === undefined) vIdx = playerIndexMap[e.VictimSID];

      if (fIdx !== undefined && vIdx !== undefined && fIdx !== vIdx) {
        secondsMatrix[fIdx][vIdx] += e.Duration || 0;
        if (secondsMatrix[fIdx][vIdx] > secondsMax) secondsMax = secondsMatrix[fIdx][vIdx];
      }
    });

    return {
      secondsMatrix,
      secondsMax: secondsMax || 1
    };
  }

  function renderFlashTab() {
    const { secondsMatrix, secondsMax } = recalcFlashMatrix(window.filteredFlashEvents || []);

    const legendEl = document.querySelector('#gridFlash .legend .small');
    if (legendEl) legendEl.textContent = '0 → ' + secondsMax.toFixed(2);

    renderMatrix({
      rootId:"#gridFlash",
      rowLabels: playerTitles,
      colLabels: playerTitles,
      data: secondsMatrix,
      maxVal: secondsMax,
      qInputId: "qFlash",
      csvBtnId: "csvFlash",
      heatToggleId: "heatFlash",
      cornerTitle: "Флешеры ↓ / Жертвы → (секунды)",
      numFmt: (v) => (typeof v === "number" ? v.toFixed(2) : String(v))
    });
  }

  window.addEventListener('dateFilterChanged', renderFlashTab);
  return { render: renderFlashTab };
})();

window.flashTabState.render();`,
		string(jPlayerMappings))
}
