package output

import (
	"encoding/json"
	"html"
	"os"
	"os/exec"
	"strings"
	"time"

	"oldfartscounter/internal/components"
	"oldfartscounter/internal/stats"
)

// HTMLGenerator отвечает за генерацию HTML
type HTMLGenerator struct {
	tournamentTab    *components.TournamentTabComponent
	killsTab         *components.KillsTabComponent
	weaponsTab       *components.WeaponsTabComponent
	flashTab         *components.FlashTabComponent
	defuseTab        *components.DefuseTabComponent
	roundsTab        *components.RoundsTabComponent
	playerRatingsTab *components.PlayerRatingsTabComponent
}

// NewHTMLGenerator создает новый генератор HTML
func NewHTMLGenerator() *HTMLGenerator {
	return &HTMLGenerator{
		tournamentTab:    components.NewTournamentTab(),
		killsTab:         components.NewKillsTab(),
		weaponsTab:       components.NewWeaponsTab(),
		flashTab:         components.NewFlashTab(),
		defuseTab:        components.NewDefuseTab(),
		roundsTab:        components.NewRoundsTab(),
		playerRatingsTab: components.NewPlayerRatingsTab(),
	}
}

// getBuildVersion возвращает информацию о сборке (дата + git hash)
func getBuildVersion() string {
	now := time.Now().Format("2006-01-02 15:04:05")

	// Пытаемся получить короткий хеш коммита
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		// Если git недоступен, возвращаем только дату
		return now
	}

	gitHash := strings.TrimSpace(string(output))
	return now + " • " + gitHash
}

// Generate генерирует полный HTML файл
func (h *HTMLGenerator) Generate(path string, data *stats.StatsData, bySteamID bool) error {
	title := "Сколько пёрнул старый?"

	// Создаем данные для JavaScript
	players := make([]string, len(data.Players))
	for i, p := range data.Players {
		players[i] = p.Title
	}

	jPlayers, _ := json.Marshal(players)
	jWeapons, _ := json.Marshal(data.Weapons)

	// Передаём агрегированные данные по дням
	jDailyKills, _ := json.Marshal(data.DailyKills)
	jDailyFlash, _ := json.Marshal(data.DailyFlash)
	jDailyDefuse, _ := json.Marshal(data.DailyDefuse)
	jDailyRounds, _ := json.Marshal(data.DailyRounds)

	// Извлекаем диапазон дат для плейсхолдеров
	minDate, maxDate := h.extractDateRange(data)

	// Получаем версию сборки
	buildVersion := getBuildVersion()

	// Генерируем HTML с правильной структурой
	htmlContent := `<!doctype html><html lang="ru"><meta charset="utf-8">
<title>` + html.EscapeString(title) + `</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="icon" type="image/png" href="favicon.png">
` + h.generateCSS() + `

<!-- Christmas Decorations -->
<div id="snowflakes-container"></div>
<div class="christmas-lights-global">
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
  <div class="light-global"></div>
</div>

<div id="loading-indicator" style="position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);background:#2b2b2b;padding:20px;border-radius:12px;border:1px solid #3a3a3a;z-index:10000;color:#c8c8c8">
  <div id="load-step-1" style="text-align:center">Шаг 1: Загрузка HTML...</div>
  <div id="load-step-2" style="text-align:center;margin-top:10px;color:#666">Шаг 2: Загрузка JS...</div>
  <div id="load-step-3" style="text-align:center;margin-top:10px;color:#666">Шаг 3: Парсинг данных...</div>
  <div id="load-step-4" style="text-align:center;margin-top:10px;color:#666">Шаг 4: Инициализация таблиц...</div>
</div>

<script>
// Самая базовая проверка - работает ли JS вообще
document.getElementById('load-step-1').style.color = '#22c55e';
</script>

<div class="header">
  <h1>Сколько пёрнул старый?</h1>
  <div class="date-filter">
    <label>Фильтр по датам:</label>
    <input type="date" id="date-from" class="date-input" min="` + minDate + `" max="` + maxDate + `">
    <span>—</span>
    <input type="date" id="date-to" class="date-input" min="` + minDate + `" max="` + maxDate + `">
    <button id="reset-dates" class="btn">Сбросить</button>
    <span class="small" style="margin-left:8px">Доступный период: ` + minDate + ` — ` + maxDate + `</span>
  </div>
</div>

<div class="tabs">
  <button class="tab-btn active" data-tab="tournament">🏆 Турнир</button>
  <button class="tab-btn" data-tab="kills">Сорян, братан</button>
  <button class="tab-btn" data-tab="player-ratings">Рейтинг игроков</button>
  <button class="tab-btn" data-tab="kw">Кто с чего убивает</button>
  <button class="tab-btn" data-tab="vw">Кого чем убивают</button>
  <button class="tab-btn" data-tab="flash">Индекс Пирога</button>
  <button class="tab-btn" data-tab="rounds">Игры</button>
  <button class="tab-btn" data-tab="defuse" style="display:none">Герои Дефьюза</button>
</div>

` + h.tournamentTab.GenerateHTML() + `
` + h.killsTab.GenerateHTML(data) + `
` + h.weaponsTab.GenerateKillerWeaponHTML(data) + `
` + h.weaponsTab.GenerateVictimWeaponHTML(data) + `
` + h.flashTab.GenerateHTML(data) + `
` + h.playerRatingsTab.GenerateHTML() + `
` + h.roundsTab.GenerateHTML() + `
` + h.defuseTab.GenerateHTML(data) + `

<div class="footer">Сборка: ` + html.EscapeString(buildVersion) + `</div>

<script>
// Шаг 2: JS начал выполняться
document.getElementById('load-step-2').style.color = '#22c55e';

// Отладка для Telegram браузера
window.onerror = function(msg, url, line, col, error) {
  var indicator = document.getElementById('loading-indicator');
  if (indicator) {
    indicator.innerHTML = '<div style="background:red;color:white;padding:20px;text-align:center"><strong>ОШИБКА JS:</strong><br>' +
      msg + '<br>Строка: ' + line + '</div>';
  }
  return false;
};

var PLAYERS, WEAPONS, DAILY_KILLS, DAILY_FLASH, DAILY_DEFUSE, DAILY_ROUNDS;
try {
  // Шаг 3: Парсим данные
  document.getElementById('load-step-3').style.color = '#fde047';

  PLAYERS = ` + string(jPlayers) + `;
  WEAPONS = ` + string(jWeapons) + `;
  DAILY_KILLS = ` + string(jDailyKills) + `;
  DAILY_FLASH = ` + string(jDailyFlash) + `;
  DAILY_DEFUSE = ` + string(jDailyDefuse) + `;
  DAILY_ROUNDS = ` + string(jDailyRounds) + `;

  document.getElementById('load-step-3').style.color = '#22c55e';
  document.getElementById('load-step-4').style.color = '#fde047';
} catch(e) {
  var indicator = document.getElementById('loading-indicator');
  if (indicator) {
    indicator.innerHTML = '<div style="background:orange;color:white;padding:20px;text-align:center"><strong>Ошибка парсинга данных:</strong><br>' +
      e.message + '</div>';
  }
  throw e;
}

` + h.generateCommonJS() + `
` + h.generateDateFilterJS() + `
` + h.generateTabsJS() + `
` + h.tournamentTab.GenerateJS(data) + `
` + h.killsTab.GenerateJS(data) + `
` + h.weaponsTab.GenerateKillerWeaponJS(data) + `
` + h.weaponsTab.GenerateVictimWeaponJS(data) + `
` + h.flashTab.GenerateJS(data) + `
` + h.playerRatingsTab.GenerateJS(data) + `
` + h.roundsTab.GenerateJS(data) + `
` + h.defuseTab.GenerateJS(data) + `

// Шаг 4 завершен
try {
  document.getElementById('load-step-4').style.color = '#22c55e';

  // Скрываем индикатор загрузки после инициализации всех табов
  setTimeout(function() {
    var loader = document.getElementById('loading-indicator');
    if (loader) {
      loader.style.display = 'none';
    }
  }, 500);
} catch(e) {
  var indicator = document.getElementById('loading-indicator');
  if (indicator) {
    indicator.innerHTML = '<div style="background:red;color:white;padding:20px;text-align:center"><strong>Ошибка финализации:</strong><br>' +
      e.message + '</div>';
  }
}

// Christmas Snow Effect
(function() {
  // Check if decorations should be shown
  var urlParams = new URLSearchParams(window.location.search);
  var hasHnyParam = urlParams.has('hny');
  var now = new Date();
  var month = now.getMonth(); // 0-11 (0=January, 11=December)
  var day = now.getDate();
  var isHolidaySeason = (month === 11 && day >= 1) || (month === 0 && day <= 10);
  var shouldShowDecorations = hasHnyParam || isHolidaySeason;

  if (!shouldShowDecorations) {
    // Hide all Christmas decorations
    var christmasLightsGlobal = document.querySelectorAll('.christmas-lights-global');
    christmasLightsGlobal.forEach(function(el) { el.style.display = 'none'; });
    var snowContainer = document.getElementById('snowflakes-container');
    if (snowContainer) snowContainer.style.display = 'none';
    return;
  }

  var snowflakesContainer = document.getElementById('snowflakes-container');
  if (!snowflakesContainer) return;

  var snowflakeSymbols = ['❄', '❅', '❆', '⛄', '🎄'];
  var numberOfSnowflakes = 50;

  function createSnowflake() {
    var snowflake = document.createElement('div');
    snowflake.className = 'snowflake';
    snowflake.innerHTML = snowflakeSymbols[Math.floor(Math.random() * snowflakeSymbols.length)];

    // Random horizontal position
    snowflake.style.left = Math.random() * 100 + 'vw';

    // Random animation duration (slower = more realistic)
    var duration = Math.random() * 8 + 8; // 8-16 seconds
    snowflake.style.animationDuration = duration + 's';

    // Random delay to stagger the snowflakes
    snowflake.style.animationDelay = Math.random() * 5 + 's';

    // Random size
    var size = Math.random() * 0.7 + 0.5; // 0.5em to 1.2em
    snowflake.style.fontSize = size + 'em';

    // Random opacity
    snowflake.style.opacity = Math.random() * 0.6 + 0.4; // 0.4 to 1.0

    snowflakesContainer.appendChild(snowflake);

    // Remove and recreate snowflake after animation completes
    setTimeout(function() {
      snowflake.remove();
      createSnowflake();
    }, (duration + parseFloat(snowflake.style.animationDelay)) * 1000);
  }

  // Create initial snowflakes
  for (var i = 0; i < numberOfSnowflakes; i++) {
    setTimeout(function() { createSnowflake(); }, i * 100);
  }
})();
</script>
</html>`

	return os.WriteFile(path, []byte(htmlContent), 0600) // #nosec G306 - HTML output file permissions
}

// generateCSS возвращает CSS стили
func (h *HTMLGenerator) generateCSS() string {
	return `<style>
:root{
  --bg:#1e1e1e; --panel:#2b2b2b; --panel-2:#232323; --text:#c8c8c8; --muted:#9aa0a6;
  --accent:#7c5cff; --grid:#3a3a3a; --sticky:#242424;
}
*{box-sizing:border-box}
body{margin:0;background:var(--bg);color:var(--text);font-family:Inter,system-ui,-apple-system,Segoe UI,Roboto,Ubuntu,Cantarell,'Noto Sans',sans-serif;position:relative;overflow-x:hidden}

/* Christmas Snow Animation */
@keyframes snowfall {
  0% {
    transform: translateY(-10px) translateX(0);
    opacity: 1;
  }
  100% {
    transform: translateY(100vh) translateX(50px);
    opacity: 0.3;
  }
}

.snowflake {
  position: fixed;
  top: -10px;
  color: white;
  font-size: 1em;
  opacity: 0.8;
  pointer-events: none;
  z-index: 1;
  animation: snowfall linear infinite;
  text-shadow: 0 0 5px rgba(255,255,255,0.5);
}

/* Christmas Lights Animation */
@keyframes twinkle {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.christmas-lights-global {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 20px;
  display: flex;
  justify-content: space-around;
  align-items: center;
  z-index: 1;
  pointer-events: none;
}

.light-global {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  box-shadow: 0 0 8px currentColor;
  animation: twinkle 1.5s ease-in-out infinite;
}

.light-global:nth-child(1) { background: #ff0000; animation-delay: 0s; }
.light-global:nth-child(2) { background: #00ff00; animation-delay: 0.3s; }
.light-global:nth-child(3) { background: #0000ff; animation-delay: 0.6s; }
.light-global:nth-child(4) { background: #ffff00; animation-delay: 0.9s; }
.light-global:nth-child(5) { background: #ff00ff; animation-delay: 1.2s; }
.light-global:nth-child(6) { background: #00ffff; animation-delay: 0.15s; }
.light-global:nth-child(7) { background: #ff8800; animation-delay: 0.45s; }
.light-global:nth-child(8) { background: #ff0088; animation-delay: 0.75s; }
.light-global:nth-child(9) { background: #88ff00; animation-delay: 1.05s; }
.light-global:nth-child(10) { background: #0088ff; animation-delay: 1.35s; }
.light-global:nth-child(11) { background: #ff4444; animation-delay: 0.2s; }
.light-global:nth-child(12) { background: #44ff44; animation-delay: 0.5s; }
.light-global:nth-child(13) { background: #4444ff; animation-delay: 0.8s; }
.light-global:nth-child(14) { background: #ffaa00; animation-delay: 1.1s; }
.light-global:nth-child(15) { background: #ff00aa; animation-delay: 1.4s; }
.header{padding:16px}
h1{margin:0 0 6px;font-size:18px}
.sub{color:var(--muted);margin:0 0 12px}
.date-filter{display:flex;align-items:center;gap:8px;margin-top:12px;flex-wrap:wrap}
.date-filter label{color:var(--text);font-size:14px}
.date-input{background:var(--panel);border:1px solid var(--grid);border-radius:8px;padding:6px 10px;color:var(--text);font-size:14px;outline:none;min-width:150px}
.date-input:focus{border-color:var(--accent)}
.date-input::-webkit-calendar-picker-indicator{filter:invert(0.7);cursor:pointer}
.tabs{display:flex;gap:6px;padding:0 16px;overflow-x:auto}
.tab-btn{background:transparent;border:1px solid var(--grid);color:var(--text);padding:8px 12px;border-radius:10px;cursor:pointer;white-space:nowrap}
.tab-btn.active{border-color:var(--accent);color:#fff;background:#2a2440}
.view{display:none;padding:16px}
.view.active{display:block}
.toolbar{display:flex;flex-wrap:wrap;gap:10px;align-items:center;margin:10px 0 14px}
.toolbar input[type="search"]{background:var(--panel);border:1px solid var(--grid);border-radius:8px;padding:8px 10px;color:var(--text);min-width:240px;outline:none}
.btn{background:var(--panel);border:1px solid var(--grid);border-radius:8px;padding:8px 12px;color:var(--text);cursor:pointer}
.btn:hover{border-color:var(--accent);color:#fff}
.select{background:var(--panel);border:1px solid var(--grid);border-radius:8px;padding:8px 10px;color:var(--text)}
.table-wrap{position:relative;overflow:auto;border:1px solid var(--grid);border-radius:12px;background:var(--panel-2);box-shadow:0 6px 24px rgba(0,0,0,.35)}
table{border-collapse:separate;border-spacing:0;min-width:max(900px, 100%)}
th,td{padding:8px 10px;border-bottom:1px solid #1f1f1f;border-right:1px solid #1f1f1f;white-space:nowrap;text-align:center}
th:last-child, td:last-child{border-right:none}
tr:last-child td{border-bottom:none}
thead th{position:sticky;top:0;background:var(--sticky);z-index:2;color:#e5e5e5;font-weight:600;font-size:12px}
th.sticky-left, td.sticky-left{position:sticky;left:0;z-index:3;background:var(--sticky);text-align:right}
.corner{z-index:4}
.sortable{cursor:pointer;position:relative;user-select:none}
.sortable:hover{background:rgba(124,92,255,0.1)}
.sortable::after{content:"⇅";margin-left:4px;opacity:0.5;font-size:10px}
.sortable:hover::after{opacity:1;color:var(--accent)}
.cell{transition:background-color .15s}
.dragging{opacity:.6}
th.sticky-left.sortable{cursor:move}
th.sticky-left.sortable::after{content:"⇄";margin-left:4px;opacity:0.5;font-size:10px}
td.sticky-left{cursor:move}
td.sticky-left:hover{background:rgba(124,92,255,0.1)}
.legend{display:flex;align-items:center;gap:8px;margin-left:auto}
.swatch{width:120px;height:10px;background:linear-gradient(90deg,#0b1020,#1e3a8a,#0ea5e9,#22c55e,#fde047,#f59e0b,#ef4444);border-radius:4px}
.small{font-size:12px;color:var(--muted)}
.footer{opacity:.7;font-size:12px;margin:12px 16px 24px}

/* Мобильная версия */
@media (max-width: 768px) {
  body{overflow-x:hidden}
  .header{padding:12px}
  h1{font-size:16px}
  .sub{font-size:11px;margin:0 0 8px}
  .date-filter{gap:4px;margin-top:8px}
  .date-filter label{font-size:12px}
  .date-input{padding:4px 8px;font-size:12px;min-width:120px}
  .date-filter .small{display:none}
  .tabs{padding:0 12px;gap:4px;-webkit-overflow-scrolling:touch;scrollbar-width:none}
  .tabs::-webkit-scrollbar{display:none}
  .tab-btn{padding:6px 10px;font-size:12px;flex-shrink:0}
  .view{padding:8px}
  .toolbar{gap:6px;margin:8px 0 10px}
  .toolbar input[type="search"]{min-width:140px;padding:6px 8px;font-size:12px}
  .btn{padding:6px 10px;font-size:12px}
  .select{padding:6px 8px;font-size:12px}
  .legend{display:none}
  .table-wrap{border-radius:8px;-webkit-overflow-scrolling:touch;margin:0 -8px;border-left:none;border-right:none;border-radius:0}
  table{min-width:600px}
  th,td{padding:4px 6px;font-size:11px}
  thead th{font-size:10px}
  th.sticky-left, td.sticky-left{font-size:10px;max-width:70px;overflow:hidden;text-overflow:ellipsis}
  .sortable::after{font-size:8px;margin-left:2px}
  th.sticky-left.sortable::after{font-size:8px;margin-left:2px}
  .footer{font-size:10px;margin:8px 12px 16px}
}

/* Очень маленькие экраны */
@media (max-width: 480px) {
  .header{padding:8px}
  h1{font-size:14px}
  .sub{font-size:10px}
  .date-filter{flex-direction:column;align-items:flex-start;width:100%}
  .date-input{width:100%;max-width:none}
  .tabs{padding:0 8px}
  .tab-btn{padding:5px 8px;font-size:11px}
  .view{padding:4px}
  .toolbar{gap:4px;flex-direction:column;align-items:stretch}
  .toolbar input[type="search"]{width:100%;min-width:auto;font-size:11px}
  .btn{padding:5px 8px;font-size:11px}
  .select{font-size:11px}
  .table-wrap{margin:0 -4px}
  table{min-width:500px}
  th,td{padding:3px 4px;font-size:10px}
  thead th{font-size:9px}
  th.sticky-left, td.sticky-left{font-size:9px;max-width:60px;overflow:hidden;text-overflow:ellipsis}
}
</style>`
}

// generateCommonJS возвращает общие JavaScript функции
func (h *HTMLGenerator) generateCommonJS() string {
	return `function heatColor(v, max){ if(!max) return "#0b1020"; const t = v/max; const h=(220*(1-t))/360, s=0.85, l=0.16+0.39*t; return hslToHex(h,s,l); }
function textColor(v, max){ const t = max? v/max : 0; return t<0.55 ? "#e5e5e5" : "#0b0b0b"; }
function hslToHex(h,s,l){
  const hue2rgb=(p,q,t)=>{ if(t<0) t+=1; if(t>1) t-=1;
    if(t<1/6) return p+(q-p)*6*t;
    if(t<1/2) return q;
    if(t<2/3) return p+(q-p)*(2/3-t)*6;
    return p;
  };
  let r,g,b;
  if(s===0){ r=g=b=l; }
  else{
    const q = l<0.5 ? l*(1+s) : l+s-l*s;
    const p = 2*l-q;
    r = hue2rgb(p,q,h+1/3);
    g = hue2rgb(p,q,h);
    b = hue2rgb(p,q,h-1/3);
  }
  const toHex=x=>Math.round(x*255).toString(16).padStart(2,'0');
  return '#'+toHex(r)+toHex(g)+toHex(b);
}
function escCSV(s){ s = String(s); if(/[",\n]/.test(s)){ return '"' + s.replace(/"/g,'""') + '"'; } return s; }
function trimLabel(s, n=24){ const arr=[...s]; return arr.length<=n ? s : arr.slice(0,n-1).join("")+"…"; }

function renderMatrix(opts){
  const {rootId, rowLabels, colLabels, data, maxVal, qInputId, csvBtnId, heatToggleId, cornerTitle, numFmt, highlightedPlayer, secondaryTarget} = opts;
  const wrap = document.querySelector(rootId);
  const thead = wrap.querySelector("thead");
  const tbody = wrap.querySelector("tbody");
  const q = document.getElementById(qInputId);
  const csvBtn = document.getElementById(csvBtnId);
  const heatToggle = document.getElementById(heatToggleId);

  let orderCols = [...Array(colLabels.length).keys()];
  let orderRows = [...Array(rowLabels.length).keys()];
  let filter = "";
  let heatOn = true;
  let dragColFrom = null, dragRowFrom = null;
  let sortedByCol = -1; // Индекс столбца, по которому идёт сортировка
  let sortedByRow = -1; // Индекс строки, по которой идёт сортировка
  let sortColAsc = false; // Направление сортировки столбцов (false = по убыванию)
  let sortRowAsc = false; // Направление сортировки строк (false = по убыванию)

  // Сохраняем индекс выделенного игрока для последующей сортировки
  let highlightedPlayerIndex = -1;
  let secondaryTargetIndex = -1;

  if(highlightedPlayer) {
    // Ищем ВСЕ вхождения игрока с таким именем и выбираем того, кого больше убивали
    const allIndices = [];
    for(let i = 0; i < colLabels.length; i++) {
      if(colLabels[i] === highlightedPlayer) {
        allIndices.push(i);
      }
    }

    console.log('DEBUG: highlightedPlayer =', highlightedPlayer);
    console.log('DEBUG: Found indices:', allIndices);

    if(allIndices.length > 0) {
      // Если несколько игроков с одинаковым ником, выбираем того, кого больше всего убивали
      if(allIndices.length > 1) {
        let maxKills = -1;
        for(const idx of allIndices) {
          let totalKills = 0;
          for(let i = 0; i < rowLabels.length; i++) {
            totalKills += (data[i] && data[i][idx] != null) ? data[i][idx] : 0;
          }
          if(totalKills > maxKills) {
            maxKills = totalKills;
            highlightedPlayerIndex = idx;
          }
        }
        console.log('DEBUG: Multiple players found, selected index', highlightedPlayerIndex, 'with', maxKills, 'total kills');
      } else {
        highlightedPlayerIndex = allIndices[0];
      }

      // Перемещаем на первое место, если еще не там
      if(highlightedPlayerIndex > 0) {
        orderCols.splice(orderCols.indexOf(highlightedPlayerIndex), 1);
        orderCols.unshift(highlightedPlayerIndex);
        console.log('DEBUG: Moved column to first position');
      }
    }
  }

  // Ищем вторичную цель (серебро)
  // Ищем по частичному совпадению (например, "Boberto" найдет "Bobo")
  if(secondaryTarget) {
    const targetLower = secondaryTarget.toLowerCase();
    for(let i = 0; i < colLabels.length; i++) {
      const labelLower = colLabels[i].toLowerCase();
      // Проверяем точное совпадение или частичное (bobo/boberto)
      if(labelLower === targetLower ||
         labelLower.includes('bobo') && targetLower.includes('bobo') ||
         labelLower.includes('boberto') && targetLower.includes('boberto')) {
        secondaryTargetIndex = i;
        console.log('DEBUG: Found secondary target:', colLabels[i], 'at index', i);
        break;
      }
    }

    // Перемещаем на второе место (после основной цели)
    if(secondaryTargetIndex > 0) {
      orderCols.splice(orderCols.indexOf(secondaryTargetIndex), 1);
      // Вставляем после highlightedPlayer (на позицию 1)
      orderCols.splice(1, 0, secondaryTargetIndex);
      console.log('DEBUG: Moved secondary target to second position');
    }
  }

  q?.addEventListener("input", ()=>{ filter = q.value.trim().toLowerCase(); draw(); });
  heatToggle?.addEventListener("change", ()=>{ heatOn = heatToggle.checked; draw(); });
  csvBtn?.addEventListener("click", ()=>{
    const rows=[];
    const keepRows = orderRows.filter(i=>rowLabels[i].toLowerCase().includes(filter));
    const keepCols = orderCols.filter(j=>colLabels[j].toLowerCase().includes(filter));
    rows.push([cornerTitle, ...keepCols.map(j=>colLabels[j])]);
    keepRows.forEach(i=>{
      rows.push([rowLabels[i], ...keepCols.map(j=>{
        const v = (data[i] && data[i][j] != null) ? data[i][j] : 0;
        return numFmt? numFmt(v) : String(v);
      })]);
    });
    const csv = rows.map(r=>r.map(escCSV).join(",")).join("\\n");
    const blob = new Blob([csv], {type:"text/csv;charset=utf-8"});
    const a=document.createElement("a"); a.href=URL.createObjectURL(blob); a.download="matrix.csv"; a.click();
  });

  // Функция сортировки по максимальному значению
  function sortByMaxValue(){
    // Находим максимальное значение и его позицию
    let maxValue = -Infinity;
    let maxRow = -1;
    let maxCol = -1;

    for(let i = 0; i < rowLabels.length; i++){
      for(let j = 0; j < colLabels.length; j++){
        const v = (data[i] && data[i][j] != null) ? data[i][j] : 0;
        const numV = (typeof v === "number") ? v : parseFloat(v) || 0;
        if(numV > maxValue){
          maxValue = numV;
          maxRow = i;
          maxCol = j;
        }
      }
    }

    if(maxRow !== -1 && maxCol !== -1){
      // Перемещаем строку с максимальным значением на первое место
      const rowPos = orderRows.indexOf(maxRow);
      if(rowPos > 0){
        orderRows.splice(rowPos, 1);
        orderRows.unshift(maxRow);
      }

      // Перемещаем столбец с максимальным значением на первое место
      const colPos = orderCols.indexOf(maxCol);
      if(colPos > 0){
        orderCols.splice(colPos, 1);
        orderCols.unshift(maxCol);
      }

      sortedByCol = -1;
      sortedByRow = -1;
      draw();
    }
  }

  // Добавляем обработчик на кнопку сортировки по максимуму (если она есть)
  const maxSortBtn = wrap.closest('.view')?.querySelector('.btn-sort-max');
  if(maxSortBtn){
    maxSortBtn.addEventListener('click', sortByMaxValue);
  }

  function sortRowsByCol(colIndex){
    // Если кликнули по тому же столбцу - меняем направление
    if(sortedByCol === colIndex) {
      sortColAsc = !sortColAsc;
    } else {
      sortColAsc = false; // По умолчанию сортируем по убыванию
    }

    const idx = orderRows.filter(i=>rowLabels[i].toLowerCase().includes(filter));
    if(sortColAsc) {
      idx.sort((a,b)=> ((data[a] && data[a][colIndex])||0) - ((data[b] && data[b][colIndex])||0));
    } else {
      idx.sort((a,b)=> ((data[b] && data[b][colIndex])||0) - ((data[a] && data[a][colIndex])||0));
    }
    orderRows = idx.concat(orderRows.filter(i=>!idx.includes(i)));
    sortedByCol = colIndex;
    sortedByRow = -1;
    draw();
  }

  function sortColsByRow(rowIndex){
    // Если кликнули по той же строке - меняем направление
    if(sortedByRow === rowIndex) {
      sortRowAsc = !sortRowAsc;
    } else {
      sortRowAsc = false; // По умолчанию сортируем по убыванию
    }

    const idx = orderCols.filter(j=>colLabels[j].toLowerCase().includes(filter));
    if(sortRowAsc) {
      idx.sort((a,b)=> ((data[rowIndex] && data[rowIndex][a])||0) - ((data[rowIndex] && data[rowIndex][b])||0));
    } else {
      idx.sort((a,b)=> ((data[rowIndex] && data[rowIndex][b])||0) - ((data[rowIndex] && data[rowIndex][a])||0));
    }
    orderCols = idx.concat(orderCols.filter(j=>!idx.includes(j)));
    sortedByRow = rowIndex;
    sortedByCol = -1;
    draw();
  }

  function dragStartCol(e){ dragColFrom = +e.target.dataset.col; e.dataTransfer.setData("text/plain","col"); e.target.classList.add("dragging"); }
  function dragOverCol(e){ if(dragColFrom!==null){ e.preventDefault(); } }
  function dropCol(e){
    e.preventDefault();
    const to = +e.target.closest("th").dataset.col;
    const fromPos = orderCols.indexOf(dragColFrom);
    const toPos = orderCols.indexOf(to);
    if(fromPos>-1 && toPos>-1 && fromPos!==toPos){
      orderCols.splice(toPos,0,...orderCols.splice(fromPos,1));
      draw();
    }
    document.querySelectorAll("th.dragging").forEach(x=>x.classList.remove("dragging"));
    dragColFrom=null;
  }
  function dragStartRow(e){ dragRowFrom = +e.target.dataset.row; e.dataTransfer.setData("text/plain","row"); e.target.classList.add("dragging"); }
  function dragOverRow(e){ if(dragRowFrom!==null){ e.preventDefault(); } }
  function dropRow(e){
    e.preventDefault();
    const to = +e.target.closest("th").dataset.row;
    const fromPos = orderRows.indexOf(dragRowFrom);
    const toPos = orderRows.indexOf(to);
    if(fromPos>-1 && toPos>-1 && fromPos!==toPos){
      orderRows.splice(toPos,0,...orderRows.splice(fromPos,1));
      draw();
    }
    document.querySelectorAll("th.dragging").forEach(x=>x.classList.remove("dragging"));
    dragRowFrom=null;
  }

  function draw(){
    // Фильтруем строки по поиску
    let keepRows = orderRows.filter(i=>rowLabels[i].toLowerCase().includes(filter));
    let keepCols = orderCols.filter(j=>colLabels[j].toLowerCase().includes(filter));

    // Убираем строки где все значения == 0
    keepRows = keepRows.filter(i=>{
      for(let j=0; j<colLabels.length; j++){
        const v = (data[i] && data[i][j] != null) ? data[i][j] : 0;
        if(v !== 0) return true;
      }
      return false;
    });

    // Убираем столбцы где все значения == 0
    keepCols = keepCols.filter(j=>{
      for(let i=0; i<rowLabels.length; i++){
        const v = (data[i] && data[i][j] != null) ? data[i][j] : 0;
        if(v !== 0) return true;
      }
      return false;
    });

    thead.innerHTML = "";
    const tr = document.createElement("tr");
    const corner = document.createElement("th"); corner.textContent = cornerTitle; corner.className="corner sticky-left sortable";
    corner.onclick = ()=> {
      orderRows = [...Array(rowLabels.length).keys()].sort((a,b)=> rowLabels[a].localeCompare(rowLabels[b]));
      orderCols = [...Array(colLabels.length).keys()].sort((a,b)=> colLabels[a].localeCompare(colLabels[b]));
      sortedByCol = -1;
      sortedByRow = -1;
      draw();
    };
    tr.appendChild(corner);
    keepCols.forEach(j=>{
      const th = document.createElement("th");
      let label = trimLabel(colLabels[j]);
      if(sortedByCol === j) {
        label += sortColAsc ? " ↑" : " ↓";
      }
      th.textContent = label;
      th.title = colLabels[j] + " — кликните для сортировки строк, перетаскивайте для смены порядка";
      th.className = "sortable";

      // Подсветка золотым для выделенного игрока
      if(highlightedPlayer && colLabels[j] === highlightedPlayer) {
        th.style.background = "linear-gradient(135deg, #FFD700 0%, #FFA500 100%)";
        th.style.color = "#000";
        th.style.fontWeight = "bold";
        th.style.boxShadow = "0 0 10px rgba(255, 215, 0, 0.5)";
      } else if(j === secondaryTargetIndex && secondaryTargetIndex >= 0) {
        // Подсветка серебром для вторичной цели (проверяем по индексу)
        th.style.background = "linear-gradient(135deg, #C0C0C0 0%, #A8A8A8 100%)";
        th.style.color = "#000";
        th.style.fontWeight = "bold";
        th.style.boxShadow = "0 0 10px rgba(192, 192, 192, 0.5)";
      } else if(sortedByCol === j) {
        th.style.background = "rgba(124,92,255,0.2)";
        th.style.color = "#fff";
      }

      th.dataset.col = j;
      th.draggable = true;
      th.addEventListener("dragstart", dragStartCol);
      th.addEventListener("dragover", dragOverCol);
      th.addEventListener("drop", dropCol);
      th.onclick = ()=> sortRowsByCol(j);
      tr.appendChild(th);
    });
    thead.appendChild(tr);
    renderBody(keepRows, keepCols);
  }
  function renderBody(rows, cols){
    tbody.innerHTML = "";
    rows.forEach(i=>{
      const tr = document.createElement("tr");
      const th = document.createElement("th");
      let label = trimLabel(rowLabels[i]);
      if(sortedByRow === i) {
        label += sortRowAsc ? " ↑" : " ↓";
      }
      th.textContent = label;
      th.title = rowLabels[i] + " — кликните для сортировки столбцов, перетаскивайте для смены порядка";
      th.className="sticky-left sortable";

      // Подсветка золотым для выделенного игрока
      if(highlightedPlayer && rowLabels[i] === highlightedPlayer) {
        th.style.background = "linear-gradient(135deg, #FFD700 0%, #FFA500 100%)";
        th.style.color = "#000";
        th.style.fontWeight = "bold";
        th.style.boxShadow = "0 0 10px rgba(255, 215, 0, 0.5)";
      } else if(i === secondaryTargetIndex && secondaryTargetIndex >= 0) {
        // Подсветка серебром для вторичной цели (проверяем по индексу)
        th.style.background = "linear-gradient(135deg, #C0C0C0 0%, #A8A8A8 100%)";
        th.style.color = "#000";
        th.style.fontWeight = "bold";
        th.style.boxShadow = "0 0 10px rgba(192, 192, 192, 0.5)";
      } else if(sortedByRow === i) {
        th.style.background = "rgba(124,92,255,0.2)";
        th.style.color = "#fff";
      }

      th.dataset.row = i;
      th.draggable = true;
      th.addEventListener("dragstart", dragStartRow);
      th.addEventListener("dragover", dragOverRow);
      th.addEventListener("drop", dropRow);
      th.onclick = ()=> sortColsByRow(i);
      tr.appendChild(th);
      cols.forEach(j=>{
        const td = document.createElement("td"); td.className="cell";
        const v = (data[i] && data[i][j] != null) ? data[i][j] : 0;
        td.textContent = numFmt ? numFmt(v) : String(v);
        td.title = rowLabels[i] + " × " + colLabels[j] + " = " + (numFmt? numFmt(v) : String(v));
        if(heatOn){
          const mv = maxVal || 1;
          const vv = (typeof v==="number")? v : parseFloat(v)||0;
          td.style.background = heatColor(vv, mv);
          td.style.color = textColor(vv, mv);
        } else { td.style.background=""; td.style.color=""; }
        tr.appendChild(td);
      });
      tbody.appendChild(tr);
    });
  }
  draw();

  // Если есть выделенный игрок, сортируем строки по его столбцу (по убыванию)
  // Иначе используем стандартную сортировку по максимальному значению
  if(highlightedPlayerIndex >= 0) {
    console.log('DEBUG: Sorting rows by column', highlightedPlayerIndex);
    sortRowsByCol(highlightedPlayerIndex);
  } else {
    console.log('DEBUG: Using sortByMaxValue');
    sortByMaxValue();
  }
}`
}

// generateDateFilterJS возвращает JavaScript для фильтрации по датам
func (h *HTMLGenerator) generateDateFilterJS() string {
	return `// Date Filter State
let DATE_FROM = null;
let DATE_TO = null;

// Функция для получения событий из агрегированных данных по датам
function getFilteredEvents(dailyData) {
  var result = [];
  for (var date in dailyData) {
    if (DATE_FROM && date < DATE_FROM) continue;
    if (DATE_TO && date > DATE_TO) continue;
    result = result.concat(dailyData[date]);
  }
  return result;
}

// Функция для пересчета статистики с учетом фильтра по датам
function recalculateStats() {
  window.filteredKillEvents = getFilteredEvents(DAILY_KILLS);
  window.filteredFlashEvents = getFilteredEvents(DAILY_FLASH);
  window.filteredDefuseEvents = getFilteredEvents(DAILY_DEFUSE);
  window.filteredRoundStats = getFilteredEvents(DAILY_ROUNDS);

  // Триггерим событие для обновления всех таблиц (совместимо со старыми браузерами)
  var event;
  if (typeof CustomEvent === 'function') {
    event = new CustomEvent('dateFilterChanged');
  } else {
    event = document.createEvent('Event');
    event.initEvent('dateFilterChanged', true, true);
  }
  window.dispatchEvent(event);
}

// Обработчики для фильтра дат
var dateFromEl = document.getElementById('date-from');
var dateToEl = document.getElementById('date-to');
var resetDatesEl = document.getElementById('reset-dates');

if (dateFromEl) {
  dateFromEl.addEventListener('change', function(e) {
    DATE_FROM = e.target.value || null;
    recalculateStats();
  });
}

if (dateToEl) {
  dateToEl.addEventListener('change', function(e) {
    DATE_TO = e.target.value || null;
    recalculateStats();
  });
}

if (resetDatesEl) {
  resetDatesEl.addEventListener('click', function() {
    DATE_FROM = null;
    DATE_TO = null;
    if (dateFromEl) dateFromEl.value = '';
    if (dateToEl) dateToEl.value = '';
    recalculateStats();
  });
}

// Инициализация фильтрованных событий (без фильтра по умолчанию)
window.filteredKillEvents = getFilteredEvents(DAILY_KILLS);
window.filteredFlashEvents = getFilteredEvents(DAILY_FLASH);
window.filteredDefuseEvents = getFilteredEvents(DAILY_DEFUSE);
window.filteredRoundStats = getFilteredEvents(DAILY_ROUNDS);
`
}

// extractDateRange извлекает минимальную и максимальную даты из событий
func (h *HTMLGenerator) extractDateRange(data *stats.StatsData) (minDate, maxDate string) {
	minDate = "1970-01-01"
	maxDate = "2099-12-31"

	// Ищем минимальную и максимальную даты среди всех событий
	dates := make(map[string]bool)

	for _, e := range data.KillEvents {
		if e.Date != "" {
			dates[e.Date] = true
		}
	}
	for _, e := range data.FlashEvents {
		if e.Date != "" {
			dates[e.Date] = true
		}
	}
	for _, e := range data.DefuseEvents {
		if e.Date != "" {
			dates[e.Date] = true
		}
	}

	if len(dates) > 0 {
		first := true
		for date := range dates {
			if first {
				minDate = date
				maxDate = date
				first = false
			} else {
				if date < minDate {
					minDate = date
				}
				if date > maxDate {
					maxDate = date
				}
			}
		}
	}

	return minDate, maxDate
}

// generateTabsJS возвращает JavaScript для управления табами
func (h *HTMLGenerator) generateTabsJS() string {
	return `// Tabs with hash link support
const btns=[...document.querySelectorAll(".tab-btn")];

// Маппинг между data-tab атрибутом и hash
const tabHashMap = {
  'tournament': 'tournament',
  'kills': 'sorrybro',
  'kw': 'killer-weapon',
  'vw': 'victim-weapon',
  'flash': 'whereispie',
  'player-ratings': 'ratings',
  'rounds': 'games',
  'defuse': 'defuse'
};

// Обратный маппинг
const hashTabMap = {};
Object.keys(tabHashMap).forEach(k => hashTabMap[tabHashMap[k]] = k);

// Сохраняем оригинальные значения фильтра для восстановления
var savedDateFrom = null;
var savedDateTo = null;
var currentTab = null; // Изначально null, чтобы первая загрузка сработала корректно

// Функция для получения дат текущего месяца
function getCurrentMonthDates() {
  var now = new Date();
  var year = now.getFullYear();
  var month = now.getMonth(); // 0-11

  // Первый день текущего месяца
  var firstDay = new Date(year, month, 1);
  var firstDayStr = firstDay.getFullYear() + '-' +
    String(firstDay.getMonth() + 1).padStart(2, '0') + '-' +
    String(firstDay.getDate()).padStart(2, '0');

  // Последний день текущего месяца
  var lastDay = new Date(year, month + 1, 0);
  var lastDayStr = lastDay.getFullYear() + '-' +
    String(lastDay.getMonth() + 1).padStart(2, '0') + '-' +
    String(lastDay.getDate()).padStart(2, '0');

  return { from: firstDayStr, to: lastDayStr };
}

// Функция переключения таба
function switchTab(tabId, updateHash) {
  var dateFromEl = document.getElementById('date-from');
  var dateToEl = document.getElementById('date-to');

  // Если переключаемся С таба kills на другой таб - восстанавливаем сохраненный фильтр
  if (currentTab === 'kills' && tabId !== 'kills') {
    DATE_FROM = savedDateFrom;
    DATE_TO = savedDateTo;
    if (dateFromEl) dateFromEl.value = savedDateFrom || '';
    if (dateToEl) dateToEl.value = savedDateTo || '';
    recalculateStats();
  }

  // Если переключаемся НА таб kills - сохраняем текущий фильтр и устанавливаем месячный
  if (tabId === 'kills' && currentTab !== 'kills') {
    savedDateFrom = DATE_FROM;
    savedDateTo = DATE_TO;

    var monthDates = getCurrentMonthDates();
    DATE_FROM = monthDates.from;
    DATE_TO = monthDates.to;
    if (dateFromEl) dateFromEl.value = monthDates.from;
    if (dateToEl) dateToEl.value = monthDates.to;
    recalculateStats();
  }

  currentTab = tabId;

  btns.forEach(x=>x.classList.remove("active"));
  document.querySelectorAll(".view").forEach(v=>v.classList.remove("active"));

  const btn = document.querySelector('.tab-btn[data-tab="'+tabId+'"]');
  const view = document.querySelector("#tab-"+tabId);

  if(btn) btn.classList.add("active");
  if(view) view.classList.add("active");

  // Обновляем hash в URL если нужно
  if(updateHash !== false) {
    const hash = tabHashMap[tabId];
    if(hash) {
      history.replaceState(null, null, '#'+hash);
    }
  }
}

// Клик по кнопке таба
btns.forEach(b=>b.addEventListener("click", ()=>{
  const tab=b.dataset.tab;
  switchTab(tab, true);
}));

// Открытие таба из hash при загрузке
function loadTabFromHash() {
  const hash = window.location.hash.substring(1); // Убираем #
  if(hash && hashTabMap[hash]) {
    switchTab(hashTabMap[hash], false);
  } else {
    // Если нет hash, по умолчанию открывается tournament
    switchTab('tournament', false);
  }
}

// Загружаем таб при загрузке страницы
loadTabFromHash();

// Обработка изменения hash (кнопки назад/вперед)
window.addEventListener('hashchange', ()=>{
  loadTabFromHash();
});`
}
