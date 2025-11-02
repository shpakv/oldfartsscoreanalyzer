package components

import (
	"encoding/json"
	"fmt"

	"oldfartscounter/internal/stats"
)

// TeamworkTabComponent отвечает за таб "Командная работа"
type TeamworkTabComponent struct{}

// NewTeamworkTab создает новый компонент таба командной работы
func NewTeamworkTab() *TeamworkTabComponent {
	return &TeamworkTabComponent{}
}

// GenerateHTML генерирует HTML для таба командной работы
func (t *TeamworkTabComponent) GenerateHTML(data *stats.StatsData) string {
	return `
<!-- TEAMWORK -->
<div id="tab-teamwork" class="view">
  <div class="toolbar">
    <input id="qTeamwork" type="search" placeholder="Поиск по именам…">
    <label class="small"><input id="heatTeamwork" type="checkbox" checked> Heatmap</label>
    <div class="legend"><div class="swatch"></div><span class="small">0 → макс</span></div>
  </div>
  <div class="table-wrap">
    <table id="gridTeamwork">
      <thead></thead>
      <tbody></tbody>
    </table>
  </div>

  <div style="margin-top: 20px;">
    <h3>Топ связки игроков</h3>
    <div class="table-wrap">
      <table id="teamCombos">
        <thead>
          <tr>
            <th>Игрок 1</th>
            <th>Игрок 2</th>
            <th>Ассисты</th>
          </tr>
        </thead>
        <tbody></tbody>
      </table>
    </div>
  </div>

  <div class="small" style="margin-top:6px">
    Матрица ассистов показывает, кто кому помогает убивать противников. Топ связки - игроки, которые чаще всего ассистируют друг другу. Клик по столбцам сортирует строки, клик по строкам сортирует столбцы.
  </div>
</div>`
}

// GenerateJS генерирует JavaScript для таба командной работы
func (t *TeamworkTabComponent) GenerateJS(data *stats.StatsData) string {
	players := make([]string, len(data.Players))
	for i, p := range data.Players {
		players[i] = p.Title
	}

	jPlayers, _ := json.Marshal(players)
	// TODO: добавить TeamworkData в StatsData
	jMatrix := []byte("[]") // Заглушка
	jCombos := []byte("[]") // Заглушка
	assistMax := 0          // Заглушка
	// jMatrix, _ := json.Marshal(data.TeamworkData.AssistMatrix)
	// jCombos, _ := json.Marshal(data.TeamworkData.TeamCombos)

	return fmt.Sprintf(`
// Init: Командная работа
renderMatrix({
  rootId:"#gridTeamwork",
  rowLabels: %s,
  colLabels: %s,
  data: %s,
  maxVal: %d,
  cornerTitle: "Ассистент",
  numFmt: v => String(v),
  heatOn: true
});

// Дополнительно: топ связки игроков
const teamCombos = %s;
const combosTable = document.getElementById('teamCombos');
let combosHtml = '';
teamCombos.forEach(combo => {
  combosHtml += '<tr><td>' + combo.Player1 + '</td><td>' + combo.Player2 + '</td><td>' + combo.Assists + '</td></tr>';
});
combosTable.querySelector('tbody').innerHTML = combosHtml;
`, jPlayers, jPlayers, jMatrix, assistMax, jCombos)
}
