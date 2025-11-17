package components

import (
	"encoding/json"
	"fmt"

	"oldfartscounter/internal/stats"
	"oldfartscounter/internal/tournament"
)

// TournamentTabComponent отвечает за таб "Турнир"
type TournamentTabComponent struct {
	config *tournament.Config
}

// NewTournamentTab создает новый компонент таба турнира
func NewTournamentTab(config *tournament.Config) *TournamentTabComponent {
	return &TournamentTabComponent{
		config: config,
	}
}

// GenerateHTML генерирует HTML для таба турнира
func (t *TournamentTabComponent) GenerateHTML() string {
	// Форматируем дату
	dateFormatted := t.formatDate(t.config.Date)

	return `
<!-- TOURNAMENT -->
<div id="tab-tournament" class="view active">
  <div style="max-width:1400px;margin:0 auto;padding:20px;">

    <!-- Christmas Decorations Styles -->
    <style>
      /* Snow Animation */
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
        z-index: 9999;
        animation: snowfall linear infinite;
      }

      /* Christmas Lights Animation */
      @keyframes twinkle {
        0%, 100% { opacity: 1; }
        50% { opacity: 0.3; }
      }

      @keyframes glow {
        0%, 100% { filter: brightness(1); }
        50% { filter: brightness(1.5); }
      }

      .christmas-lights {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 30px;
        display: flex;
        justify-content: space-around;
        align-items: center;
        z-index: 10;
      }

      .light {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        box-shadow: 0 0 10px currentColor;
        animation: twinkle 1.5s ease-in-out infinite;
      }

      .light:nth-child(1) { background: #ff0000; animation-delay: 0s; }
      .light:nth-child(2) { background: #00ff00; animation-delay: 0.3s; }
      .light:nth-child(3) { background: #0000ff; animation-delay: 0.6s; }
      .light:nth-child(4) { background: #ffff00; animation-delay: 0.9s; }
      .light:nth-child(5) { background: #ff00ff; animation-delay: 1.2s; }
      .light:nth-child(6) { background: #00ffff; animation-delay: 0.15s; }
      .light:nth-child(7) { background: #ff8800; animation-delay: 0.45s; }
      .light:nth-child(8) { background: #ff0088; animation-delay: 0.75s; }
      .light:nth-child(9) { background: #88ff00; animation-delay: 1.05s; }
      .light:nth-child(10) { background: #0088ff; animation-delay: 1.35s; }

      .hero-christmas {
        position: relative;
      }

      /* Tooltip Styles */
      .tooltip-wrapper {
        position: relative;
        display: inline-block;
        cursor: help;
      }

      .tooltip-wrapper .tooltip-text {
        visibility: hidden;
        opacity: 0;
        width: 280px;
        background: #1a1a1a;
        color: var(--text);
        text-align: center;
        border-radius: 6px;
        padding: 10px;
        position: absolute;
        z-index: 1000;
        bottom: 125%;
        left: 50%;
        margin-left: -140px;
        transition: opacity 0.3s;
        border: 1px solid var(--accent);
        box-shadow: 0 4px 12px rgba(0,0,0,0.5);
        font-size: 13px;
        line-height: 1.5;
      }

      .tooltip-wrapper .tooltip-text::after {
        content: "";
        position: absolute;
        top: 100%;
        left: 50%;
        margin-left: -5px;
        border-width: 5px;
        border-style: solid;
        border-color: var(--accent) transparent transparent transparent;
      }

      .tooltip-wrapper:hover .tooltip-text {
        visibility: visible;
        opacity: 1;
      }
    </style>

    <!-- Snowflakes Container -->
    <div id="snowflakes-container"></div>

    <!-- Hero Section -->
    <div class="hero-christmas" style="background:linear-gradient(135deg, #1a1a1a 0%, #2d1b4e 100%);border:2px solid #cfb53b;border-radius:16px;padding:40px;text-align:center;margin-bottom:30px;position:relative;overflow:hidden;">
      <!-- Christmas Lights -->
      <div class="christmas-lights">
        <div class="light"></div>
        <div class="light"></div>
        <div class="light"></div>
        <div class="light"></div>
        <div class="light"></div>
        <div class="light"></div>
        <div class="light"></div>
        <div class="light"></div>
        <div class="light"></div>
        <div class="light"></div>
      </div>

      <div style="position:absolute;top:-50px;right:-50px;width:200px;height:200px;background:radial-gradient(circle, rgba(207,181,59,0.1) 0%, transparent 70%);border-radius:50%;"></div>
      <div style="position:absolute;bottom:-30px;left:-30px;width:150px;height:150px;background:radial-gradient(circle, rgba(239,68,68,0.1) 0%, transparent 70%);border-radius:50%;"></div>

      <div style="font-size:48px;margin-bottom:10px;">🎄🏆🎄</div>
      <h1 style="color:#cfb53b;font-size:36px;margin:0 0 10px;text-transform:uppercase;letter-spacing:2px;text-shadow:0 2px 10px rgba(207,181,59,0.5);">
        3-й Ежегодный Турнир OldFarts
      </h1>
      <div style="color:var(--muted);font-size:18px;margin-bottom:20px;">
        Где геморрой не помеха победе
      </div>
      <div style="display:inline-block;background:rgba(239,68,68,0.2);border:1px solid #ef4444;border-radius:8px;padding:12px 24px;">
        <div style="color:var(--muted);font-size:12px;text-transform:uppercase;letter-spacing:1px;margin-bottom:4px;">Дата проведения</div>
        <div style="color:#ef4444;font-size:24px;font-weight:bold;">` + dateFormatted + `</div>
        <div style="color:var(--muted);font-size:14px;margin-top:4px;">` + t.config.StartTime + `</div>
      </div>
    </div>

    <!-- Info Cards -->
    <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:16px;margin-bottom:30px;">

      <!-- Duration Card -->
      <div style="background:var(--panel);border:1px solid var(--grid);border-radius:12px;padding:20px;text-align:center;">
        <div style="font-size:32px;margin-bottom:8px;">⏱️</div>
        <div style="color:var(--muted);font-size:12px;text-transform:uppercase;margin-bottom:6px;">Длительность</div>
        <div style="color:var(--text);font-size:20px;font-weight:bold;">~` + fmt.Sprintf("%d", t.config.DurationHours) + ` часа</div>
      </div>

      <!-- Format Card -->
      <div style="background:var(--panel);border:1px solid var(--grid);border-radius:12px;padding:20px;text-align:center;">
        <div style="font-size:32px;margin-bottom:8px;">🎮</div>
        <div style="color:var(--muted);font-size:12px;text-transform:uppercase;margin-bottom:6px;">Формат</div>
        <div style="color:var(--text);font-size:20px;font-weight:bold;">
          <span class="tooltip-wrapper">
            BO3
            <span style="font-size:10px;color:var(--muted);vertical-align:super;">ⓘ</span>
            <span class="tooltip-text">Best of 3 — для победы нужно выиграть 2 из 3 карт</span>
          </span>, 2 этапа
        </div>
      </div>

      <!-- Teams Card -->
      <div style="background:var(--panel);border:1px solid var(--grid);border-radius:12px;padding:20px;text-align:center;">
        <div style="font-size:32px;margin-bottom:8px;">👥</div>
        <div style="color:var(--muted);font-size:12px;text-transform:uppercase;margin-bottom:6px;">Участники</div>
        <div style="color:var(--text);font-size:20px;font-weight:bold;">` + fmt.Sprintf("%d", len(t.config.Participants)) + ` игроков</div>
        <div style="color:var(--muted);font-size:11px;margin-top:4px;">` + fmt.Sprintf("%d команды по %d", t.config.Teams.Count, t.config.Teams.Size) + `</div>
      </div>

      <!-- Prizes Card -->
      <div style="background:var(--panel);border:1px solid var(--grid);border-radius:12px;padding:20px;text-align:center;">
        <div style="font-size:32px;margin-bottom:8px;">💰</div>
        <div style="color:var(--muted);font-size:12px;text-transform:uppercase;margin-bottom:6px;">Призы</div>
        <div style="color:#f59e0b;font-size:14px;font-weight:bold;">Скоро...</div>
      </div>

    </div>

    <!-- Tournament Bracket -->
    <div style="background:var(--panel);border:2px solid #4b69ff;border-radius:12px;padding:30px;margin-bottom:30px;">
      <div style="text-align:center;margin-bottom:25px;">
        <div style="font-size:40px;margin-bottom:10px;">🎯</div>
        <h2 style="color:#4b69ff;font-size:28px;margin:0 0 8px;text-transform:uppercase;letter-spacing:2px;">Турнирная сетка</h2>
        <div style="color:var(--muted);font-size:16px;">Структура турнира Best of 3</div>
      </div>

      <!-- Bracket Visual -->
      <div style="display:flex;justify-content:space-around;align-items:center;flex-wrap:wrap;gap:30px;">

        <!-- Stage 1: Semifinals -->
        <div style="flex:1;min-width:300px;display:flex;flex-direction:column;">
          <div style="text-align:center;margin-bottom:20px;">
            <div style="background:#4b69ff;color:white;display:inline-block;padding:8px 20px;border-radius:20px;font-weight:bold;text-transform:uppercase;font-size:14px;">Этап 1: Полуфиналы</div>
          </div>

          <!-- Match 1 -->
          <div style="background:rgba(75,105,255,0.1);border:2px solid #4b69ff;border-radius:8px;padding:16px;">
            <div style="text-align:center;color:#4b69ff;font-size:13px;font-weight:bold;text-transform:uppercase;margin-bottom:8px;">Полуфинал 1</div>
            <div style="display:flex;justify-content:space-between;align-items:center;">
              <div style="color:var(--text);font-size:18px;font-weight:bold;">Team 1</div>
              <div style="color:var(--muted);font-size:14px;">vs</div>
              <div style="color:var(--text);font-size:18px;font-weight:bold;">Team 2</div>
            </div>
          </div>

          <div style="height:30px;"></div>

          <!-- Match 2 -->
          <div style="background:rgba(75,105,255,0.1);border:2px solid #4b69ff;border-radius:8px;padding:16px;">
            <div style="text-align:center;color:#4b69ff;font-size:13px;font-weight:bold;text-transform:uppercase;margin-bottom:8px;">Полуфинал 2</div>
            <div style="display:flex;justify-content:space-between;align-items:center;">
              <div style="color:var(--text);font-size:18px;font-weight:bold;">Team 3</div>
              <div style="color:var(--muted);font-size:14px;">vs</div>
              <div style="color:var(--text);font-size:18px;font-weight:bold;">Team 4</div>
            </div>
          </div>
        </div>

        <!-- Stage 2: Finals -->
        <div style="flex:1;min-width:300px;display:flex;flex-direction:column;">
          <div style="text-align:center;margin-bottom:20px;">
            <div style="background:#cfb53b;color:#1a1a1a;display:inline-block;padding:8px 20px;border-radius:20px;font-weight:bold;text-transform:uppercase;font-size:14px;">Этап 2: Финалы</div>
          </div>

          <!-- Final Match -->
          <div style="background:rgba(207,181,59,0.15);border:2px solid #cfb53b;border-radius:8px;padding:16px;">
            <div style="text-align:center;color:#cfb53b;font-size:13px;font-weight:bold;text-transform:uppercase;margin-bottom:8px;">🏆 Финал</div>
            <div style="display:flex;justify-content:space-between;align-items:center;">
              <div style="color:var(--text);font-size:18px;font-weight:bold;">Победитель 1</div>
              <div style="color:var(--muted);font-size:14px;">vs</div>
              <div style="color:var(--text);font-size:18px;font-weight:bold;">Победитель 2</div>
            </div>
          </div>

          <div style="height:30px;"></div>

          <!-- 3rd Place Match -->
          <div style="background:rgba(239,68,68,0.1);border:2px solid #ef4444;border-radius:8px;padding:16px;">
            <div style="text-align:center;color:#ef4444;font-size:13px;font-weight:bold;text-transform:uppercase;margin-bottom:8px;">🥉 За 3 место</div>
            <div style="display:flex;justify-content:space-between;align-items:center;">
              <div style="color:var(--text);font-size:18px;font-weight:bold;">Проигравший 1</div>
              <div style="color:var(--muted);font-size:14px;">vs</div>
              <div style="color:var(--text);font-size:18px;font-weight:bold;">Проигравший 2</div>
            </div>
          </div>
        </div>

      </div>

      <!-- Info note -->
      <div style="margin-top:25px;padding:15px;background:rgba(75,105,255,0.05);border-left:3px solid #4b69ff;border-radius:6px;">
        <div style="color:var(--muted);font-size:13px;line-height:1.6;">
          💡 <strong>Между этапами:</strong> перерыв 10 минут для выбора карт и отдыха. Все 20 участников играют во всех этапах.<br>
          🖥️ <strong>Сервера:</strong> Полуфинал 1 и Финал — Сервер 1. Полуфинал 2 и Матч за 3 место — Сервер 2.
        </div>
      </div>
    </div>

    <!-- Participants Section -->
    <div style="background:var(--panel);border:2px solid #8847ff;border-radius:12px;padding:30px;margin-bottom:30px;">
      <div style="text-align:center;margin-bottom:25px;">
        <div style="font-size:40px;margin-bottom:10px;">👤</div>
        <h2 style="color:#8847ff;font-size:28px;margin:0 0 8px;text-transform:uppercase;letter-spacing:2px;">Участники</h2>
        <div style="color:var(--muted);font-size:16px;">` + fmt.Sprintf("%d", len(t.config.Participants)) + ` подтвержденных игроков (все места заняты)</div>
      </div>

      <div id="participantsList"></div>
    </div>

    <!-- Teams Section -->
    <div style="background:var(--panel);border:2px solid #cfb53b;border-radius:12px;padding:30px;margin-bottom:30px;">
      <div style="text-align:center;margin-bottom:25px;">
        <div style="font-size:40px;margin-bottom:10px;">🛡️</div>
        <h2 style="color:#cfb53b;font-size:28px;margin:0 0 8px;text-transform:uppercase;letter-spacing:2px;">Составы команд</h2>
        <div style="color:var(--muted);font-size:16px;">Распределение по рейтингу</div>
      </div>

      <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(300px,1fr));gap:20px;">

        <!-- Team 1 -->
        <div style="background:linear-gradient(135deg, rgba(239,68,68,0.1) 0%, rgba(239,68,68,0.05) 100%);border:2px solid #ef4444;border-radius:12px;padding:20px;">
          <div style="text-align:center;margin-bottom:16px;">
            <div style="background:#ef4444;color:white;display:inline-block;padding:8px 20px;border-radius:20px;font-weight:bold;font-size:16px;text-transform:uppercase;">Team 1</div>
          </div>
          <div id="team1Roster" style="color:var(--muted);font-style:italic;text-align:center;font-size:14px;">Будет объявлено позже</div>
        </div>

        <!-- Team 2 -->
        <div style="background:linear-gradient(135deg, rgba(75,105,255,0.1) 0%, rgba(75,105,255,0.05) 100%);border:2px solid #4b69ff;border-radius:12px;padding:20px;">
          <div style="text-align:center;margin-bottom:16px;">
            <div style="background:#4b69ff;color:white;display:inline-block;padding:8px 20px;border-radius:20px;font-weight:bold;font-size:16px;text-transform:uppercase;">Team 2</div>
          </div>
          <div id="team2Roster" style="color:var(--muted);font-style:italic;text-align:center;font-size:14px;">Будет объявлено позже</div>
        </div>

        <!-- Team 3 -->
        <div style="background:linear-gradient(135deg, rgba(34,197,94,0.1) 0%, rgba(34,197,94,0.05) 100%);border:2px solid #22c55e;border-radius:12px;padding:20px;">
          <div style="text-align:center;margin-bottom:16px;">
            <div style="background:#22c55e;color:white;display:inline-block;padding:8px 20px;border-radius:20px;font-weight:bold;font-size:16px;text-transform:uppercase;">Team 3</div>
          </div>
          <div id="team3Roster" style="color:var(--muted);font-style:italic;text-align:center;font-size:14px;">Будет объявлено позже</div>
        </div>

        <!-- Team 4 -->
        <div style="background:linear-gradient(135deg, rgba(245,158,11,0.1) 0%, rgba(245,158,11,0.05) 100%);border:2px solid #f59e0b;border-radius:12px;padding:20px;">
          <div style="text-align:center;margin-bottom:16px;">
            <div style="background:#f59e0b;color:#1a1a1a;display:inline-block;padding:8px 20px;border-radius:20px;font-weight:bold;font-size:16px;text-transform:uppercase;">Team 4</div>
          </div>
          <div id="team4Roster" style="color:var(--muted);font-style:italic;text-align:center;font-size:14px;">Будет объявлено позже</div>
        </div>

      </div>

      <div style="margin-top:20px;padding:15px;background:rgba(207,181,59,0.05);border-left:3px solid #cfb53b;border-radius:6px;">
        <div style="color:var(--muted);font-size:13px;line-height:1.6;">
          👑 <strong>Капитаны:</strong> Каждая команда выбирает капитана любым способом. Капитан представляет команду при выборе карт.<br>
          ⚖️ <strong>Распределение:</strong> По рейтингу игроков. Ответственный: Слава Шпак (Mr. Titspervert)
        </div>
      </div>
    </div>

    <!-- Tournament Details -->
    <div style="background:var(--panel);border:2px solid var(--grid);border-radius:12px;padding:30px;margin-bottom:30px;">
      <div style="text-align:center;margin-bottom:25px;">
        <div style="font-size:40px;margin-bottom:10px;">📖</div>
        <h2 style="color:var(--accent);font-size:28px;margin:0 0 8px;text-transform:uppercase;letter-spacing:2px;">Детали турнира</h2>
      </div>

      <!-- Collapsible Sections -->
      <div id="tournamentDetails">

        <!-- Rules Section -->
        <div class="detail-section" style="border-bottom:1px solid var(--grid);padding:20px 0;">
          <div class="detail-header" style="cursor:pointer;display:flex;justify-content:space-between;align-items:center;" onclick="toggleDetail('rules')">
            <div style="display:flex;align-items:center;gap:12px;">
              <span style="font-size:24px;">📋</span>
              <span style="color:var(--text);font-size:18px;font-weight:bold;">Правила участия</span>
            </div>
            <span id="arrow-rules" style="color:var(--accent);font-size:20px;transition:transform 0.3s;">▼</span>
          </div>
          <div id="detail-rules" style="display:none;margin-top:15px;padding-left:36px;color:var(--muted);font-size:14px;line-height:1.8;">
            <p>• Турнир является дружеским соревнованием с собственными правилами</p>
            <p>• Посторонним запрещено заходить на сервера во время соревнования</p>
            <p>• Менять ники на время турнира запрещено</p>
            <p>• Любой софт, дающий игровое преимущество, запрещен. Нарушитель будет изгнан</p>
            <p>• <strong style="color:#22c55e;">Have fun и не забывайте пикать зайчика!</strong></p>
          </div>
        </div>

        <!-- Format Section -->
        <div class="detail-section" style="border-bottom:1px solid var(--grid);padding:20px 0;">
          <div class="detail-header" style="cursor:pointer;display:flex;justify-content:space-between;align-items:center;" onclick="toggleDetail('format')">
            <div style="display:flex;align-items:center;gap:12px;">
              <span style="font-size:24px;">🗺️</span>
              <span style="color:var(--text);font-size:18px;font-weight:bold;">Формат и выбор карт</span>
            </div>
            <span id="arrow-format" style="color:var(--accent);font-size:20px;transition:transform 0.3s;">▼</span>
          </div>
          <div id="detail-format" style="display:none;margin-top:15px;padding-left:36px;color:var(--muted);font-size:14px;line-height:1.8;">
            <p><strong>Best of 3 (BO3):</strong> для победы нужно выиграть 2 из 3 карт</p>
            <p><strong>Овертаймы включены</strong> для определения победителя</p>
            <p><strong>Выбор карт:</strong> через сервис mapban.gg капитанами команд</p>
            <p><strong>Пул карт:</strong> Ancient, Anubis, Dust2, Inferno, Mirage, Nuke, Overpass, Train, Vertigo</p>
            <p><strong>Порядок:</strong> бан → бан → пик → пик → бан → бан → бан → бан → выбор стороны</p>
            <p>Первой голосует команда, слабее по общему рейтингу</p>
          </div>
        </div>

        <!-- Teams Section -->
        <div class="detail-section" style="border-bottom:1px solid var(--grid);padding:20px 0;">
          <div class="detail-header" style="cursor:pointer;display:flex;justify-content:space-between;align-items:center;" onclick="toggleDetail('teams')">
            <div style="display:flex;align-items:center;gap:12px;">
              <span style="font-size:24px;">👥</span>
              <span style="color:var(--text);font-size:18px;font-weight:bold;">Команды и капитаны</span>
            </div>
            <span id="arrow-teams" style="color:var(--accent);font-size:20px;transition:transform 0.3s;">▼</span>
          </div>
          <div id="detail-teams" style="display:none;margin-top:15px;padding-left:36px;color:var(--muted);font-size:14px;line-height:1.8;">
            <p><strong>Команды:</strong> 4 команды по 5 человек (Team 1, Team 2, Team 3, Team 4)</p>
            <p><strong>Распределение:</strong> по рейтингу, отвечает Слава Шпак (Mr. Titspervert)</p>
            <p><strong>Капитаны:</strong> выбираются командой любым способом</p>
            <p><strong>Роль капитана:</strong> представитель команды при выборе карт</p>
            <p>Запрещено прыгать между командами даже на разминке</p>
          </div>
        </div>

        <!-- Communication Section -->
        <div class="detail-section" style="padding:20px 0;">
          <div class="detail-header" style="cursor:pointer;display:flex;justify-content:space-between;align-items:center;" onclick="toggleDetail('comm')">
            <div style="display:flex;align-items:center;gap:12px;">
              <span style="font-size:24px;">💬</span>
              <span style="color:var(--text);font-size:18px;font-weight:bold;">Общение и форс-мажоры</span>
            </div>
            <span id="arrow-comm" style="color:var(--accent);font-size:20px;transition:transform 0.3s;">▼</span>
          </div>
          <div id="detail-comm" style="display:none;margin-top:15px;padding-left:36px;color:var(--muted);font-size:14px;line-height:1.8;">
            <p><strong>Telegram:</strong> создайте отдельный чат команды для предварительной коммуникации</p>
            <p><strong>Discord:</strong> голосовые каналы по командам во время турнира</p>
            <p><strong>Перерывы:</strong> если игрок вылетел, команда может инициировать внутриигровое голосование</p>
            <p><strong>Рестарты:</strong> возможны по обоюдному согласию капитанов</p>
            <p>Официальной разминки нет, участники могут прийти заранее</p>
          </div>
        </div>

      </div>
    </div>

    <!-- Qualification Requirements -->
    <div style="background:var(--panel);border:2px solid #ef4444;border-radius:12px;padding:30px;margin-bottom:30px;">
      <div style="text-align:center;margin-bottom:20px;">
        <div style="font-size:40px;margin-bottom:10px;">⚔️</div>
        <h2 style="color:#ef4444;font-size:28px;margin:0 0 8px;text-transform:uppercase;letter-spacing:2px;">Квалификация</h2>
        <div style="color:var(--muted);font-size:16px;">Требования для участия в турнире</div>
      </div>

      <div style="background:rgba(239,68,68,0.1);border-radius:8px;padding:20px;margin-bottom:20px;">
        <div style="display:flex;align-items:center;justify-content:center;gap:12px;margin-bottom:12px;">
          <div style="font-size:32px;">📅</div>
          <div>
            <div style="color:#ef4444;font-size:18px;font-weight:bold;">Минимум 100 раундов</div>
            <div style="color:var(--muted);font-size:14px;">Сыграно с 1 сентября 2025 года</div>
          </div>
        </div>
        <div style="text-align:center;color:var(--muted);font-size:13px;font-style:italic;">
          Это требование гарантирует стабильность рейтинга и справедливость распределения по командам
        </div>
      </div>

      <div id="qualificationStatus"></div>
    </div>

  </div>
</div>

<script>
// Toggle function for collapsible sections
function toggleDetail(section) {
  const detail = document.getElementById('detail-' + section);
  const arrow = document.getElementById('arrow-' + section);

  if (detail.style.display === 'none') {
    detail.style.display = 'block';
    arrow.style.transform = 'rotate(180deg)';
  } else {
    detail.style.display = 'none';
    arrow.style.transform = 'rotate(0deg)';
  }
}

// Christmas Snow Effect
(function() {
  if (!shouldShowHolidayDecorations()) {
    // Hide all Christmas decorations
    const christmasLights = document.querySelectorAll('.christmas-lights');
    christmasLights.forEach(el => el.style.display = 'none');
    const snowContainer = document.getElementById('snowflakes-container');
    if (snowContainer) snowContainer.style.display = 'none';
    return;
  }

  const snowflakesContainer = document.getElementById('snowflakes-container');
  if (!snowflakesContainer) return;

  const snowflakeSymbols = ['❄', '❅', '❆', '⛄', '🎄'];
  const numberOfSnowflakes = 50;

  function createSnowflake() {
    const snowflake = document.createElement('div');
    snowflake.className = 'snowflake';
    snowflake.innerHTML = snowflakeSymbols[Math.floor(Math.random() * snowflakeSymbols.length)];

    // Random horizontal position
    snowflake.style.left = Math.random() * 100 + 'vw';

    // Random animation duration (slower = more realistic)
    const duration = Math.random() * 8 + 8; // 8-16 seconds
    snowflake.style.animationDuration = duration + 's';

    // Random delay to stagger the snowflakes
    snowflake.style.animationDelay = Math.random() * 5 + 's';

    // Random size
    const size = Math.random() * 0.7 + 0.5; // 0.5em to 1.2em
    snowflake.style.fontSize = size + 'em';

    // Random opacity
    snowflake.style.opacity = Math.random() * 0.6 + 0.4; // 0.4 to 1.0

    snowflakesContainer.appendChild(snowflake);

    // Remove and recreate snowflake after animation completes
    setTimeout(() => {
      snowflake.remove();
      createSnowflake();
    }, (duration + parseFloat(snowflake.style.animationDelay)) * 1000);
  }

  // Create initial snowflakes
  for (let i = 0; i < numberOfSnowflakes; i++) {
    setTimeout(() => createSnowflake(), i * 100);
  }
})();
</script>
`
}

// GenerateJS генерирует JavaScript для таба турнира
func (t *TournamentTabComponent) GenerateJS(data *stats.StatsData) string {
	jRatings, _ := json.Marshal(data.PlayerRatings)
	jParticipants, _ := json.Marshal(t.config.Participants)
	minRounds := data.MinRoundsForRating

	return fmt.Sprintf(`
// Init: Турнир
window.tournamentTabState = (function() {
  const originalRatings = %s;
  const MIN_ROUNDS = %v;
  const QUALIFICATION_DATE = new Date('2025-09-01');

  // Список подтвержденных участников турнира (из конфига)
  const confirmedParticipants = %s;

  // Дополнительные участники (лист ожидания) - для тестирования
  // Раскомментируй следующую строку для демонстрации системы приоритетов:
  const waitlistParticipants = [];
  /*
  const waitlistParticipants = [
    'ProGamer228', 'SneakyBeaky', 'Rush B No Stop', 'EZ4ENCE', 'Бабай',
    'Тапок', 'Silent Killer', 'Luntik', 'Макс Пейн', 'xXx_ProSniper_xXx'
  ];
  */

  // Составы команд (null = не объявлено)
  // Формат: { players: ['Имя1', 'Имя2', ...], captain: 'ИмяКапитана' }
  const teamRosters = {
    team1: { players: ['Пердун1', 'Пердун2', 'Пердун3', 'Пердун4', 'Пердун5'], captain: 'Пердун1' },
    team2: { players: ['Пердун6', 'Пердун7', 'Пердун8', 'Пердун9', 'Пердун10'], captain: 'Пердун6' },
    team3: { players: ['Пердун11', 'Пердун12', 'Пердун13', 'Пердун14', 'Пердун15'], captain: 'Пердун11' },
    team4: { players: ['Пердун16', 'Пердун17', 'Пердун18', 'Пердун19', 'Пердун20'], captain: 'Пердун16' }
  };

  function getPlayerData(name, playerDataMap) {
    for (var accountID in playerDataMap) {
      if (playerDataMap[accountID].Name === name) {
        return playerDataMap[accountID];
      }
    }
    return null;
  }

  function renderTeamRosters() {
    const teams = ['team1', 'team2', 'team3', 'team4'];

    teams.forEach(teamId => {
      const roster = teamRosters[teamId];
      const element = document.getElementById(teamId + 'Roster');

      if (!element) return;

      if (!roster || !roster.players || roster.players.length === 0) {
        element.innerHTML = '<div style="color:var(--muted);font-style:italic;text-align:center;font-size:14px;">Будет объявлено позже</div>';
        return;
      }

      let html = '<div style="display:flex;flex-direction:column;gap:8px;">';

      roster.players.forEach(playerName => {
        const isCaptain = playerName === roster.captain;
        const bgColor = isCaptain ? 'rgba(207,181,59,0.15)' : 'rgba(255,255,255,0.03)';
        const borderColor = isCaptain ? '#cfb53b' : 'var(--grid)';

        html += '<div style="background:' + bgColor + ';border-left:3px solid ' + borderColor + ';padding:10px 12px;border-radius:6px;display:flex;justify-content:space-between;align-items:center;">';
        html += '<div style="font-weight:' + (isCaptain ? 'bold' : 'normal') + ';color:var(--text);font-size:14px;">' + playerName + '</div>';

        if (isCaptain) {
          html += '<div style="background:#cfb53b;color:#1a1a1a;padding:2px 8px;border-radius:10px;font-size:10px;font-weight:bold;">👑 КАПИТАН</div>';
        }

        html += '</div>';
      });

      html += '</div>';

      element.innerHTML = html;
    });
  }

  function renderParticipantsList() {
    const roundStats = window.filteredRoundStats || [];
    const playerData = {};

    // Агрегируем раунды после 1 сентября 2025
    roundStats.forEach(function(round) {
      const roundDate = new Date(round.Date);
      if (roundDate < QUALIFICATION_DATE) return;

      round.Players.forEach(function(playerStats) {
        if (playerStats.AccountID === 0) return;

        if (!playerData[playerStats.AccountID]) {
          playerData[playerStats.AccountID] = {
            AccountID: playerStats.AccountID,
            Name: '',
            RoundsAfterSept: 0
          };
        }

        playerData[playerStats.AccountID].RoundsAfterSept++;
      });
    });

    // Находим имена игроков из оригинальных рейтингов
    originalRatings.forEach(function(orig) {
      if (playerData[orig.AccountID]) {
        playerData[orig.AccountID].Name = orig.Name;
      }
    });

    // Создаем полный список всех участников с приоритетом
    const allParticipants = [];

    // Добавляем подтвержденных участников
    confirmedParticipants.forEach((name, index) => {
      const found = getPlayerData(name, playerData);
      const rounds = found ? found.RoundsAfterSept : 0;
      const qualified = rounds >= MIN_ROUNDS;

      allParticipants.push({
        name: name,
        rounds: rounds,
        qualified: qualified,
        confirmed: true,
        registrationOrder: index + 1,
        priority: qualified ? 1 : 3 // Приоритет 1: квал+подтв, Приоритет 3: неквал+подтв
      });
    });

    // Добавляем участников из листа ожидания
    waitlistParticipants.forEach((name, index) => {
      const found = getPlayerData(name, playerData);
      const rounds = found ? found.RoundsAfterSept : 0;
      const qualified = rounds >= MIN_ROUNDS;

      allParticipants.push({
        name: name,
        rounds: rounds,
        qualified: qualified,
        confirmed: false,
        registrationOrder: confirmedParticipants.length + index + 1,
        priority: qualified ? 2 : 4 // Приоритет 2: квал+неподтв, Приоритет 4: неквал+неподтв
      });
    });

    // Сортируем по приоритету, затем по порядку регистрации
    allParticipants.sort((a, b) => {
      if (a.priority !== b.priority) return a.priority - b.priority;
      return a.registrationOrder - b.registrationOrder;
    });

    // Разделяем на допущенных (первые 20) и лист ожидания
    const admitted = allParticipants.slice(0, 20);
    const waitlist = allParticipants.slice(20);

    let html = '';

    // Секция допущенных участников
    html += '<div style="margin-bottom:30px;">';
    html += '<div style="background:rgba(34,197,94,0.1);border-left:4px solid #22c55e;padding:12px 16px;border-radius:8px;margin-bottom:16px;">';
    html += '<div style="display:flex;align-items:center;gap:10px;">';
    html += '<span style="font-size:24px;">✅</span>';
    html += '<div>';
    html += '<div style="color:#22c55e;font-weight:bold;font-size:16px;">ДОПУЩЕНЫ К ТУРНИРУ</div>';
    html += '<div style="color:var(--muted);font-size:13px;">' + admitted.length + ' из 20 мест занято</div>';
    html += '</div>';
    html += '</div>';
    html += '</div>';

    html += '<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:12px;">';

    admitted.forEach((p, idx) => {
      const borderColor = p.qualified ? '#22c55e' : '#ef4444';
      const bgColor = p.qualified ? 'rgba(34,197,94,0.05)' : 'rgba(239,68,68,0.05)';
      const statusIcon = p.qualified ? '✅' : '❌';
      const statusText = p.qualified ? 'Квалифицирован' : p.rounds + '/' + MIN_ROUNDS;

      const priorityBadge = !p.confirmed && p.qualified
        ? '<span style="background:#4b69ff;color:white;padding:2px 8px;border-radius:10px;font-size:10px;margin-left:6px;">ЛИСТ → ОСНОВА</span>'
        : '';

      html += '<div style="background:' + bgColor + ';border-left:3px solid ' + borderColor + ';border-radius:8px;padding:12px;position:relative;">';
      html += '<div style="position:absolute;top:8px;right:8px;color:var(--muted);font-size:11px;font-weight:bold;z-index:1;">#' + (idx + 1) + '</div>';
      html += '<div style="display:flex;justify-content:space-between;align-items:flex-start;margin-bottom:4px;padding-right:30px;">';
      html += '<div style="font-weight:bold;color:var(--text);font-size:14px;word-break:break-word;">' + p.name + '</div>';
      html += '<div style="font-size:18px;flex-shrink:0;margin-left:8px;">' + statusIcon + '</div>';
      html += '</div>';
      html += '<div style="color:var(--muted);font-size:12px;">' + statusText + '</div>';
      if (priorityBadge) {
        html += '<div style="margin-top:6px;">' + priorityBadge + '</div>';
      }
      html += '</div>';
    });

    // Пустые слоты
    const emptySlots = 20 - admitted.length;
    for (let i = 0; i < emptySlots; i++) {
      html += '<div style="background:rgba(156,163,175,0.05);border:2px dashed var(--grid);border-radius:8px;padding:12px;display:flex;align-items:center;justify-content:center;position:relative;">';
      html += '<div style="position:absolute;top:8px;right:8px;color:var(--muted);font-size:11px;font-weight:bold;">#' + (admitted.length + i + 1) + '</div>';
      html += '<div style="color:var(--muted);font-size:14px;text-align:center;">🎭<br>Свободно</div>';
      html += '</div>';
    }

    html += '</div></div>';

    // Секция листа ожидания
    if (waitlist.length > 0) {
      // Группируем по приоритету
      const priority2 = waitlist.filter(p => p.priority === 2);
      const priority3 = waitlist.filter(p => p.priority === 3);
      const priority4 = waitlist.filter(p => p.priority === 4);

      html += '<div style="margin-top:30px;">';
      html += '<div style="background:rgba(245,158,11,0.1);border-left:4px solid #f59e0b;padding:12px 16px;border-radius:8px;margin-bottom:16px;">';
      html += '<div style="display:flex;align-items:center;gap:10px;">';
      html += '<span style="font-size:24px;">⏳</span>';
      html += '<div>';
      html += '<div style="color:#f59e0b;font-weight:bold;font-size:16px;">ЛИСТ ОЖИДАНИЯ</div>';
      html += '<div style="color:var(--muted);font-size:13px;">' + waitlist.length + ' участников в очереди</div>';
      html += '</div>';
      html += '</div>';
      html += '</div>';

      // Приоритет 2: Квалифицированы, но не в основном списке
      if (priority2.length > 0) {
        html += '<div style="margin-bottom:20px;">';
        html += '<div style="color:#4b69ff;font-weight:bold;font-size:14px;margin-bottom:10px;padding-left:4px;">🔵 Приоритет 2 — Квалифицированы (первые в очереди)</div>';
        html += '<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:12px;">';

        priority2.forEach((p, idx) => {
          const isFirst = idx === 0;
          const borderStyle = isFirst ? 'border:2px solid #4b69ff;' : 'border-left:3px solid #4b69ff;';
          const glowStyle = isFirst ? 'box-shadow:0 0 12px rgba(75,105,255,0.4);' : '';

          html += '<div style="background:rgba(75,105,255,0.05);' + borderStyle + 'border-radius:8px;padding:12px;position:relative;' + glowStyle + '">';
          html += '<div style="position:absolute;top:8px;right:8px;color:var(--muted);font-size:11px;font-weight:bold;">#' + (20 + waitlist.indexOf(p) + 1) + '</div>';
          if (isFirst) {
            html += '<div style="background:#4b69ff;color:white;display:inline-block;padding:2px 8px;border-radius:10px;font-size:10px;margin-bottom:6px;font-weight:bold;">СЛЕДУЮЩИЙ</div>';
          }
          html += '<div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:4px;">';
          html += '<div style="font-weight:bold;color:var(--text);font-size:14px;max-width:180px;overflow:hidden;text-overflow:ellipsis;">' + p.name + '</div>';
          html += '<div style="font-size:18px;">✅</div>';
          html += '</div>';
          html += '<div style="color:var(--muted);font-size:12px;">Квалифицирован (' + p.rounds + ' раундов)</div>';
          html += '</div>';
        });

        html += '</div></div>';
      }

      // Приоритет 3: Не квалифицированы, но в основном списке (не должно быть в листе ожидания)
      if (priority3.length > 0) {
        html += '<div style="margin-bottom:20px;">';
        html += '<div style="color:#8847ff;font-weight:bold;font-size:14px;margin-bottom:10px;padding-left:4px;">🟣 Приоритет 3 — Подтверждены, но не квалифицированы</div>';
        html += '<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:12px;">';

        priority3.forEach(p => {
          html += '<div style="background:rgba(136,71,255,0.05);border-left:3px solid #8847ff;border-radius:8px;padding:12px;position:relative;">';
          html += '<div style="position:absolute;top:8px;right:8px;color:var(--muted);font-size:11px;font-weight:bold;">#' + (20 + waitlist.indexOf(p) + 1) + '</div>';
          html += '<div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:4px;">';
          html += '<div style="font-weight:bold;color:var(--text);font-size:14px;max-width:180px;overflow:hidden;text-overflow:ellipsis;">' + p.name + '</div>';
          html += '<div style="font-size:18px;">❌</div>';
          html += '</div>';
          html += '<div style="color:var(--muted);font-size:12px;">' + p.rounds + '/' + MIN_ROUNDS + '</div>';
          html += '</div>';
        });

        html += '</div></div>';
      }

      // Приоритет 4: Не квалифицированы и не в основном списке
      if (priority4.length > 0) {
        html += '<div style="margin-bottom:20px;">';
        html += '<div style="color:#9ca3af;font-weight:bold;font-size:14px;margin-bottom:10px;padding-left:4px;">⚪ Приоритет 4 — Не квалифицированы</div>';
        html += '<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:12px;">';

        priority4.forEach(p => {
          html += '<div style="background:rgba(156,163,175,0.05);border-left:3px solid #9ca3af;border-radius:8px;padding:12px;position:relative;">';
          html += '<div style="position:absolute;top:8px;right:8px;color:var(--muted);font-size:11px;font-weight:bold;">#' + (20 + waitlist.indexOf(p) + 1) + '</div>';
          html += '<div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:4px;">';
          html += '<div style="font-weight:bold;color:var(--text);font-size:14px;max-width:180px;overflow:hidden;text-overflow:ellipsis;">' + p.name + '</div>';
          html += '<div style="font-size:18px;">❌</div>';
          html += '</div>';
          html += '<div style="color:var(--muted);font-size:12px;">' + p.rounds + '/' + MIN_ROUNDS + '</div>';
          html += '</div>';
        });

        html += '</div></div>';
      }

      html += '</div>';
    }

    const participantsListEl = document.getElementById('participantsList');
    if (participantsListEl) {
      participantsListEl.innerHTML = html;
    }
  }

  function calculateQualificationStatus() {
    // Пересчитываем раунды с учетом фильтра дат
    const roundStats = window.filteredRoundStats || [];
    const playerData = {};

    // Агрегируем раунды после 1 сентября 2025
    roundStats.forEach(function(round) {
      const roundDate = new Date(round.Date);
      if (roundDate < QUALIFICATION_DATE) return;

      round.Players.forEach(function(playerStats) {
        if (playerStats.AccountID === 0) return;

        if (!playerData[playerStats.AccountID]) {
          playerData[playerStats.AccountID] = {
            AccountID: playerStats.AccountID,
            Name: '',
            RoundsAfterSept: 0
          };
        }

        playerData[playerStats.AccountID].RoundsAfterSept++;
      });
    });

    // Находим имена игроков из оригинальных рейтингов
    originalRatings.forEach(function(orig) {
      if (playerData[orig.AccountID]) {
        playerData[orig.AccountID].Name = orig.Name;
      }
    });

    // Преобразуем в массив и сортируем по количеству раундов (убывание)
    const players = [];
    for (var accountID in playerData) {
      const player = playerData[accountID];
      if (!player.Name) {
        player.Name = 'Player_' + player.AccountID;
      }
      player.RoundsNeeded = Math.max(0, MIN_ROUNDS - player.RoundsAfterSept);
      player.IsQualified = player.RoundsAfterSept >= MIN_ROUNDS;
      players.push(player);
    }

    players.sort(function(a, b) {
      return b.RoundsAfterSept - a.RoundsAfterSept;
    });

    return players;
  }

  function renderQualificationStatus() {
    const players = calculateQualificationStatus();
    const qualifiedCount = players.filter(p => p.IsQualified).length;
    const totalCount = players.length;

    let html = '<div style="margin-bottom:20px;">';
    html += '<div style="display:flex;gap:20px;justify-content:center;margin-bottom:20px;">';

    // Qualified Counter
    html += '<div style="background:rgba(34,197,94,0.1);border:2px solid #22c55e;border-radius:12px;padding:20px;flex:1;max-width:250px;text-align:center;">';
    html += '<div style="color:#22c55e;font-size:36px;font-weight:bold;">' + qualifiedCount + '</div>';
    html += '<div style="color:var(--muted);font-size:14px;margin-top:4px;">✅ Квалифицировано</div>';
    html += '</div>';

    // Not Qualified Counter
    html += '<div style="background:rgba(239,68,68,0.1);border:2px solid #ef4444;border-radius:12px;padding:20px;flex:1;max-width:250px;text-align:center;">';
    html += '<div style="color:#ef4444;font-size:36px;font-weight:bold;">' + (totalCount - qualifiedCount) + '</div>';
    html += '<div style="color:var(--muted);font-size:14px;margin-top:4px;">❌ Не квалифицировано</div>';
    html += '</div>';

    html += '</div></div>';

    // Players Table
    html += '<div style="overflow-x:auto;">';
    html += '<table style="width:100%%;border-collapse:collapse;">';
    html += '<thead><tr style="background:var(--sticky);text-align:left;border-bottom:2px solid var(--grid);">';
    html += '<th style="padding:12px;text-align:center;width:60px;">#</th>';
    html += '<th style="padding:12px;">Игрок</th>';
    html += '<th style="padding:12px;text-align:center;">Раундов с 01.09.2025</th>';
    html += '<th style="padding:12px;text-align:center;">Осталось до 100</th>';
    html += '<th style="padding:12px;text-align:center;">Статус</th>';
    html += '</tr></thead>';
    html += '<tbody>';

    players.forEach((player, idx) => {
      const rowStyle = player.IsQualified
        ? 'background:rgba(34,197,94,0.05);border-left:3px solid #22c55e;'
        : 'background:rgba(239,68,68,0.05);border-left:3px solid #ef4444;';

      html += '<tr style="border-bottom:1px solid var(--grid);' + rowStyle + '">';
      html += '<td style="padding:12px;text-align:center;color:var(--muted);font-weight:bold;">' + (idx + 1) + '</td>';
      html += '<td style="padding:12px;font-weight:bold;">' + player.Name + '</td>';
      html += '<td style="padding:12px;text-align:center;font-size:18px;font-weight:bold;color:' + (player.IsQualified ? '#22c55e' : '#ef4444') + ';">' + player.RoundsAfterSept + '</td>';

      if (player.IsQualified) {
        html += '<td style="padding:12px;text-align:center;color:#22c55e;font-weight:bold;">—</td>';
        html += '<td style="padding:12px;text-align:center;"><span style="background:#22c55e;color:white;padding:6px 16px;border-radius:20px;font-size:12px;font-weight:bold;text-transform:uppercase;">✅ Допущен</span></td>';
      } else {
        html += '<td style="padding:12px;text-align:center;font-size:18px;font-weight:bold;color:#ef4444;">' + player.RoundsNeeded + '</td>';
        html += '<td style="padding:12px;text-align:center;"><span style="background:#ef4444;color:white;padding:6px 16px;border-radius:20px;font-size:12px;font-weight:bold;text-transform:uppercase;">❌ Не допущен</span></td>';
      }

      html += '</tr>';
    });

    html += '</tbody></table>';
    html += '</div>';

    document.getElementById('qualificationStatus').innerHTML = html;
  }

  // Подписываемся на изменение фильтра дат
  window.addEventListener('dateFilterChanged', function() {
    renderParticipantsList();
    renderQualificationStatus();
  });

  return {
    render: function() {
      renderTeamRosters();
      renderParticipantsList();
      renderQualificationStatus();
    }
  };
})();

// Начальная отрисовка
window.tournamentTabState.render();
`, string(jRatings), minRounds, string(jParticipants))
}

// formatDate форматирует дату в читаемый вид
func (t *TournamentTabComponent) formatDate(dateStr string) string {
	// Форматируем дату из "2025-12-20" в "20 декабря 2025"
	months := map[string]string{
		"01": "января", "02": "февраля", "03": "марта", "04": "апреля",
		"05": "мая", "06": "июня", "07": "июля", "08": "августа",
		"09": "сентября", "10": "октября", "11": "ноября", "12": "декабря",
	}

	if len(dateStr) != 10 {
		return dateStr
	}

	year := dateStr[0:4]
	month := dateStr[5:7]
	day := dateStr[8:10]

	// Убираем ведущий ноль у дня
	if day[0] == '0' {
		day = day[1:]
	}

	monthName, ok := months[month]
	if !ok {
		return dateStr
	}

	return fmt.Sprintf("%s %s %s", day, monthName, year)
}
