package components

import (
	"encoding/json"
	"fmt"

	"oldfartscounter/internal/stats"
)

// DefuseTabComponent отвечает за таб "Герои Дефьюза"
type DefuseTabComponent struct{}

// NewDefuseTab создает новый компонент таба дефьюза
func NewDefuseTab() *DefuseTabComponent {
	return &DefuseTabComponent{}
}

// GenerateHTML генерирует HTML для таба дефьюза
func (d *DefuseTabComponent) GenerateHTML(data *stats.StatsData) string {
	return `
<!-- DEFUSE -->
<div id="tab-defuse" class="view">
  <div class="toolbar">
    <input id="qDefuse" type="search" placeholder="Поиск по именам…">
    <label class="small"><input id="heatDefuse" type="checkbox" checked> Heatmap</label>
    <div class="legend"><div class="swatch"></div><span class="small">0 → макс</span></div>
  </div>
  <div class="table-wrap">
    <table id="gridDefuse">
      <thead></thead>
      <tbody></tbody>
    </table>
  </div>
  <div class="small" style="margin-top:6px">
    Статистика дефьюза бомбы. "С дефьюз-китом" = 5 секунд, "Без дефьюз-кита" = 10 секунд.
  </div>
</div>`
}

// GenerateJS генерирует JavaScript для таба дефьюза
func (d *DefuseTabComponent) GenerateJS(data *stats.StatsData) string {
	players := make([]string, len(data.Players))
	for i, p := range data.Players {
		players[i] = p.Title
	}

	jPlayers, _ := json.Marshal(players)
	jAttempts, _ := json.Marshal(data.DefuseData.Attempts)
	jWithKit, _ := json.Marshal(data.DefuseData.WithKit)
	jWithoutKit, _ := json.Marshal(data.DefuseData.WithoutKit)
	jSuccessWithKit, _ := json.Marshal(data.DefuseData.SuccessWithKit)
	jSuccessWithoutKit, _ := json.Marshal(data.DefuseData.SuccessWithoutKit)
	jAbandoned, _ := json.Marshal(data.DefuseData.Abandoned)
	jFailed, _ := json.Marshal(data.DefuseData.Failed)

	return fmt.Sprintf(`
// Init: Герои Дефьюза (объединенная таблица)
function drawDefuse(){
  const players = %s;
  const attempts = %s;
  const withKit = %s;
  const withoutKit = %s;
  const successWithKit = %s;
  const successWithoutKit = %s;
  const abandoned = %s;
  const failed = %s;

  const columns = [
    {title: "Всего", data: attempts},
    {title: "Успешно без кита", data: successWithoutKit},
    {title: "Успешно с китом", data: successWithKit},
    {title: "Не успел", data: failed}
  ];

  renderDefuseMatrix({
    rootId: "#gridDefuse",
    players: players,
    columns: columns,
    qInputId: "qDefuse",
    csvBtnId: "csvDefuse",
    heatToggleId: "heatDefuse"
  });
}

function renderDefuseMatrix(opts){
  const {rootId, players, columns, qInputId, csvBtnId, heatToggleId} = opts;
  const wrap = document.querySelector(rootId);
  const thead = wrap.querySelector("thead");
  const tbody = wrap.querySelector("tbody");
  const q = document.getElementById(qInputId);
  const csvBtn = document.getElementById(csvBtnId);
  const heatToggle = document.getElementById(heatToggleId);

  let filter = "";
  let heatOn = true;

  // Вычисляем максимальное значение для каждого столбца
  const maxValues = columns.map(col => Math.max(...col.data));

  q?.addEventListener("input", ()=>{ filter = q.value.trim().toLowerCase(); draw(); });
  heatToggle?.addEventListener("change", ()=>{ heatOn = heatToggle.checked; draw(); });
  csvBtn?.addEventListener("click", ()=>{
    const rows = [];
    const keepIdx = players.map((_,i)=>i).filter(i=>players[i].toLowerCase().includes(filter));
    rows.push(["Игрок", ...columns.map(col => col.title)]);
    keepIdx.forEach(i=>{
      rows.push([players[i], ...columns.map(col => String(col.data[i] || 0))]);
    });
    const csv = rows.map(r=>r.map(escCSV).join(",")).join("\\n");
    const blob = new Blob([csv], {type:"text/csv;charset=utf-8"});
    const a=document.createElement("a"); a.href=URL.createObjectURL(blob); a.download="defuse_stats.csv"; a.click();
  });

  function draw(){
    const keepIdx = players.map((_,i)=>i).filter(i=>players[i].toLowerCase().includes(filter));

    // Сортируем по общему количеству попыток (первый столбец)
    keepIdx.sort((a,b)=> (columns[0].data[b] || 0) - (columns[0].data[a] || 0));

    // Создаем заголовок
    thead.innerHTML = "";
    const tr = document.createElement("tr");

    const th1 = document.createElement("th");
    th1.textContent = "Игрок";
    th1.className = "sticky-left sortable";
    th1.onclick = ()=> {
      const sorted = keepIdx.slice().sort((a,b)=> players[a].localeCompare(players[b]));
      renderBody(sorted);
    };
    tr.appendChild(th1);

    columns.forEach((col, colIdx) => {
      const th = document.createElement("th");
      th.textContent = col.title;
      th.className = "sortable";
      th.onclick = ()=> {
        const sorted = keepIdx.slice().sort((a,b)=> (col.data[b] || 0) - (col.data[a] || 0));
        renderBody(sorted);
      };
      tr.appendChild(th);
    });

    thead.appendChild(tr);
    renderBody(keepIdx);
  }

  function renderBody(indices){
    tbody.innerHTML = "";
    indices.forEach(i=>{
      const tr = document.createElement("tr");

      // Колонка с именем игрока
      const td1 = document.createElement("td");
      td1.textContent = trimLabel(players[i]);
      td1.title = players[i];
      td1.className = "sticky-left";
      tr.appendChild(td1);

      // Колонки с данными
      columns.forEach((col, colIdx) => {
        const td = document.createElement("td");
        td.className = "cell";
        const v = col.data[i] || 0;
        td.textContent = String(v);
        td.title = players[i] + " - " + col.title + ": " + v;

        if(heatOn && maxValues[colIdx] > 0){
          td.style.background = heatColor(v, maxValues[colIdx]);
          td.style.color = textColor(v, maxValues[colIdx]);
        } else {
          td.style.background = "";
          td.style.color = "";
        }

        tr.appendChild(td);
      });

      tbody.appendChild(tr);
    });
  }

  draw();
}

drawDefuse();`,
		string(jPlayers),           // %s для players
		string(jAttempts),          // %s для attempts
		string(jWithKit),           // %s для withKit
		string(jWithoutKit),        // %s для withoutKit
		string(jSuccessWithKit),    // %s для successWithKit
		string(jSuccessWithoutKit), // %s для successWithoutKit
		string(jAbandoned),         // %s для abandoned
		string(jFailed))            // %s для failed
}
