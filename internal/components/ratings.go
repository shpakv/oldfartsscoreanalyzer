package components

import (
	"fmt"

	"oldfartscounter/internal/stats"
)

// RatingsTabComponent отвечает за таб "Рейтинг игроков"
type RatingsTabComponent struct{}

// NewRatingsTab создает новый компонент таба рейтингов
func NewRatingsTab() *RatingsTabComponent {
	return &RatingsTabComponent{}
}

// GenerateHTML генерирует HTML для таба рейтингов
func (r *RatingsTabComponent) GenerateHTML(data *stats.StatsData) string {
	return `
<!-- RATINGS -->
<div id="tab-ratings" class="view">
  <div class="toolbar">
    <input id="qRatings" type="search" placeholder="Поиск по именам…">
  </div>
  <div class="table-wrap">
    <table id="ratingsTable">
      <thead>
        <tr>
          <th class="sortable" data-sort="PlayerName">Игрок <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="Rating">Рейтинг <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="GamesPlayed">Игр <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="TotalKills">Убийства <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="TotalDeaths">Смерти <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="TotalAssists">Ассисты <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="AverageKD">K/D <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="AverageADR">ADR <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="TotalDamage">Урон <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="TotalScore">Счет <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="TotalMVP">MVP <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="AverageHeadshot">Попадания в голову % <span class="sort-indicator"></span></th>
          <th class="sortable" data-sort="AverageEfficiency">Эффективность <span class="sort-indicator"></span></th>
        </tr>
      </thead>
      <tbody></tbody>
    </table>
  </div>
  <div class="small" style="margin-top:6px">
    Рейтинг игроков основан на JSON статистике из логов CS2. Учитывается K/D, ADR, эффективность и процент попаданий в голову. Чем выше рейтинг, тем лучше игрок.
  </div>
</div>`
}

// GenerateJS генерирует JavaScript для таба рейтингов
func (r *RatingsTabComponent) GenerateJS(data *stats.StatsData) string {
	// TODO: добавить PlayerStats в StatsData
	jRatings := []byte("[]") // Заглушка
	// jRatings, _ := json.Marshal(data.PlayerStats.PlayerRatings)

	return fmt.Sprintf(`
// Init: Рейтинг игроков
const ratingsData = %s;
let currentRatingsData = [...ratingsData]; // Копия для сортировки
const ratingsTable = document.getElementById('ratingsTable');
const ratingsSearch = document.getElementById('qRatings');

// Состояние сортировки
let currentSort = { field: 'Rating', order: 'desc' };

function renderRatingsTable(filteredData = currentRatingsData) {
  let html = '';
  filteredData.forEach((rating, index) => {
    html += '<tr draggable="true" data-index="' + index + '">' +
      '<td>' + rating.PlayerName + '</td>' +
      '<td>' + rating.Rating.toFixed(2) + '</td>' +
      '<td>' + rating.GamesPlayed + '</td>' +
      '<td>' + rating.TotalKills + '</td>' +
      '<td>' + rating.TotalDeaths + '</td>' +
      '<td>' + rating.TotalAssists + '</td>' +
      '<td>' + rating.AverageKD.toFixed(2) + '</td>' +
      '<td>' + rating.AverageADR.toFixed(1) + '</td>' +
      '<td>' + rating.TotalDamage + '</td>' +
      '<td>' + rating.TotalScore + '</td>' +
      '<td>' + rating.TotalMVP + '</td>' +
      '<td>' + rating.AverageHeadshot.toFixed(1) + '%%</td>' +
      '<td>' + rating.AverageEfficiency.toFixed(2) + '</td>' +
      '</tr>';
  });
  ratingsTable.querySelector('tbody').innerHTML = html;
  updateSortIndicators();
  initDragAndDrop();
}

// Функция сортировки
function sortRatingsData(field, order = 'asc') {
  currentRatingsData.sort((a, b) => {
    let aVal = a[field];
    let bVal = b[field];

    // Для строковых полей
    if (typeof aVal === 'string') {
      aVal = aVal.toLowerCase();
      bVal = bVal.toLowerCase();
    }

    if (order === 'asc') {
      return aVal < bVal ? -1 : aVal > bVal ? 1 : 0;
    } else {
      return aVal > bVal ? -1 : aVal < bVal ? 1 : 0;
    }
  });
}

// Обновление индикаторов сортировки
function updateSortIndicators() {
  document.querySelectorAll('.sort-indicator').forEach(indicator => {
    indicator.textContent = '';
  });

  const activeHeader = document.querySelector('[data-sort="' + currentSort.field + '"] .sort-indicator');
  if (activeHeader) {
    activeHeader.textContent = currentSort.order === 'asc' ? ' ↑' : ' ↓';
  }
}

// Обработчики кликов на заголовки
document.querySelectorAll('.sortable').forEach(header => {
  header.addEventListener('click', () => {
    const field = header.dataset.sort;

    // Если кликнули по тому же полю, меняем направление
    if (currentSort.field === field) {
      currentSort.order = currentSort.order === 'asc' ? 'desc' : 'asc';
    } else {
      currentSort.field = field;
      currentSort.order = field === 'PlayerName' ? 'asc' : 'desc'; // Имена по алфавиту, числа по убыванию
    }

    sortRatingsData(currentSort.field, currentSort.order);
    renderRatingsTable();
  });
});

// Drag and Drop функциональность
function initDragAndDrop() {
  const rows = ratingsTable.querySelectorAll('tbody tr');

  rows.forEach(row => {
    row.addEventListener('dragstart', (e) => {
      e.dataTransfer.setData('text/plain', row.dataset.index);
      row.classList.add('dragging');
    });

    row.addEventListener('dragend', () => {
      row.classList.remove('dragging');
    });

    row.addEventListener('dragover', (e) => {
      e.preventDefault();
    });

    row.addEventListener('drop', (e) => {
      e.preventDefault();
      const fromIndex = parseInt(e.dataTransfer.getData('text/plain'));
      const toIndex = parseInt(row.dataset.index);

      if (fromIndex !== toIndex) {
        // Перемещаем элемент в массиве
        const item = currentRatingsData.splice(fromIndex, 1)[0];
        currentRatingsData.splice(toIndex, 0, item);
        renderRatingsTable();
      }
    });
  });
}

// Поиск
ratingsSearch.addEventListener('input', (e) => {
  const query = e.target.value.toLowerCase();
  if (query === '') {
    // Если поиск пуст, показываем все данные с текущей сортировкой
    currentRatingsData = [...ratingsData];
    sortRatingsData(currentSort.field, currentSort.order);
    renderRatingsTable();
  } else {
    // Фильтруем данные и отображаем
    const filtered = currentRatingsData.filter(rating =>
      rating.PlayerName.toLowerCase().includes(query)
    );
    renderRatingsTable(filtered);
  }
});

// Начальная сортировка по рейтингу (по убыванию) и отрисовка
sortRatingsData(currentSort.field, currentSort.order);
renderRatingsTable();
`, string(jRatings))
}
