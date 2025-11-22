package components

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// TreeTabComponent отвечает за таб "Древо Пердунов"
type TreeTabComponent struct{}

// NewTreeTab создает новый компонент таба древа
func NewTreeTab() *TreeTabComponent {
	return &TreeTabComponent{}
}

// TreeNode представляет узел в древе приглашений
type TreeNode struct {
	Telegram     string
	FirstName    string
	LastName     string
	Nickname     string
	Inviter      string
	DateInvited  string
	DateExcluded string
	Status       string
	Children     []*TreeNode
}

// GenerateHTML генерирует HTML для таба древа
func (t *TreeTabComponent) GenerateHTML() string {
	return `
<!-- TREE (Old Farts Tree) -->
<div id="tab-tree" class="view">
  <div class="toolbar">
    <label style="display:flex;align-items:center;gap:8px">
      <input id="hideInactive" type="checkbox"> Скрыть пассивов
    </label>
    <span class="small" style="margin-left:auto">🌳 Древо приглашений в Old Farts</span>
  </div>
  <div id="tree-container" style="padding:40px 20px;overflow-x:auto;overflow-y:auto;max-height:calc(100vh - 200px)"></div>
</div>`
}

// GenerateJS генерирует JavaScript для таба древа
func (t *TreeTabComponent) GenerateJS() string {
	// Загружаем и парсим CSV
	tree, err := t.loadTreeData()
	if err != nil {
		return fmt.Sprintf(`
// Tree data loading error: %s
document.getElementById('tree-container').innerHTML = '<div style="color:red;text-align:center;padding:20px">Ошибка загрузки данных: %s</div>';
`, err.Error(), err.Error())
	}

	// Строим древо
	treeHTML := t.buildTreeHTML(tree)

	return fmt.Sprintf(`
// Init: Древо Пердунов
(function() {
  const treeContainer = document.getElementById('tree-container');
  const hideInactiveCheckbox = document.getElementById('hideInactive');
  const treeData = %s;

  function renderTree(showInactive) {
    if (showInactive) {
      treeContainer.innerHTML = treeData;
    } else {
      // Удаляем неактивные узлы
      const temp = document.createElement('div');
      temp.innerHTML = treeData;
      const inactiveNodes = temp.querySelectorAll('.tree-node.inactive');
      inactiveNodes.forEach(node => {
        const branch = node.closest('.tree-branch');
        if (!branch) return;

        // Проверяем, есть ли у этого узла дети
        const hasChildren = branch.querySelector('.tree-children');

        if (hasChildren) {
          // Если есть дети, скрываем только сам узел, но оставляем детей
          node.style.display = 'none';
        } else {
          // Если детей нет, удаляем весь branch
          branch.remove();
        }
      });
      treeContainer.innerHTML = temp.innerHTML;
    }
  }

  hideInactiveCheckbox.addEventListener('change', function() {
    renderTree(!this.checked);
  });

  // Начальный рендер
  renderTree(true);
})();
`, "`"+treeHTML+"`")
}

// loadTreeData загружает данные из CSV файла
func (t *TreeTabComponent) loadTreeData() (map[string]*TreeNode, error) {
	file, err := os.Open("members.csv") // #nosec G304 -- path is controlled by application code
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	nodes := make(map[string]*TreeNode)

	// Пропускаем заголовок
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) < 11 {
			continue
		}

		// Новый формат CSV:
		// 0: Статус членства, 1: Статус активности, 2: Имя, 3: Фамилия, 4: Никнейм
		// 5: SteamID3, 6: TG tag, 7: E-mail, 8: Кем приглашен, 9: Дата вступления, 10: Дата исключения
		firstName := strings.TrimSpace(record[2])
		lastName := strings.TrimSpace(record[3])
		nickname := strings.TrimSpace(record[4])
		invited := strings.TrimSpace(record[6])
		inviter := strings.TrimSpace(record[8])
		dateInvited := strings.TrimSpace(record[9])
		dateExcluded := strings.TrimSpace(record[10])
		status := strings.TrimSpace(record[1])

		// Пропускаем пустые записи
		if invited == "" {
			continue
		}

		// Нормализуем inviter: "-" означает "нет приглашающего"
		if inviter == "-" || inviter == "" {
			inviter = ""
		}

		// Создаем узел для приглашенного если его нет
		if _, exists := nodes[invited]; !exists {
			nodes[invited] = &TreeNode{
				Telegram:     invited,
				FirstName:    firstName,
				LastName:     lastName,
				Nickname:     nickname,
				Inviter:      inviter,
				DateInvited:  dateInvited,
				DateExcluded: dateExcluded,
				Status:       status,
				Children:     []*TreeNode{},
			}
		} else {
			// Обновляем информацию если узел уже существует
			nodes[invited].FirstName = firstName
			nodes[invited].LastName = lastName
			nodes[invited].Nickname = nickname
			nodes[invited].Inviter = inviter
			nodes[invited].DateInvited = dateInvited
			nodes[invited].DateExcluded = dateExcluded
			nodes[invited].Status = status
		}

		// Создаем узел для приглашающего если он указан
		if inviter != "" {
			if _, exists := nodes[inviter]; !exists {
				nodes[inviter] = &TreeNode{
					Telegram: inviter,
					Children: []*TreeNode{},
				}
			}

			// Добавляем связь
			nodes[inviter].Children = append(nodes[inviter].Children, nodes[invited])
		}
	}

	return nodes, nil
}

// buildTreeHTML строит HTML представление древа
func (t *TreeTabComponent) buildTreeHTML(nodes map[string]*TreeNode) string {
	// Находим корневые узлы (тех, кого никто не пригласил)
	roots := []*TreeNode{}
	invitedSet := make(map[string]bool)

	for _, node := range nodes {
		if node.Inviter != "" {
			invitedSet[node.Telegram] = true
		}
	}

	for _, node := range nodes {
		if !invitedSet[node.Telegram] {
			roots = append(roots, node)
		}
	}

	// Сортируем корневые узлы по дате
	sort.Slice(roots, func(i, j int) bool {
		dateI, _ := time.Parse("1/2/2006", roots[i].DateInvited)
		dateJ, _ := time.Parse("1/2/2006", roots[j].DateInvited)
		return dateI.Before(dateJ)
	})

	var sb strings.Builder
	sb.WriteString(`<div class="tree-root">`)

	for i, root := range roots {
		if i > 0 {
			sb.WriteString(`<div class="tree-root-separator"></div>`)
		}
		t.renderNode(&sb, root, true)
	}

	sb.WriteString(`</div>`)
	return sb.String()
}

// renderNode рекурсивно рендерит узел и его детей
func (t *TreeTabComponent) renderNode(sb *strings.Builder, node *TreeNode, isRoot bool) {
	// Определяем класс статуса
	statusClass := "active"
	if node.Status == "Inactive" {
		statusClass = "inactive"
	}

	// Проверяем, исключен ли участник
	isExcluded := node.DateExcluded != ""
	if isExcluded {
		statusClass += " excluded"
	}

	// Форматируем дату из mm/dd/yyyy в dd.mm.yyyy
	dateStr := ""
	if node.DateInvited != "" {
		parsedDate, err := time.Parse("1/2/2006", node.DateInvited)
		if err == nil {
			dateStr = parsedDate.Format("02.01.2006")
		} else {
			dateStr = node.DateInvited
		}
	}

	// Убираем @ из telegram handle
	telegram := strings.TrimPrefix(node.Telegram, "@")

	// Формируем отображаемое имя (только никнейм, без реальных имен)
	displayName := telegram
	if node.Nickname != "" {
		displayName = node.Nickname
	}

	// Начинаем branch
	sb.WriteString(`<div class="tree-branch">`)

	// Рендерим узел
	rootClass := ""
	if isRoot {
		rootClass = " root"
	}
	sb.WriteString(fmt.Sprintf(`<div class="tree-node %s%s" data-telegram="%s">`, statusClass, rootClass, telegram))
	sb.WriteString(`<div class="tree-node-content">`)

	// Аватар (первая буква ника или telegram)
	avatar := "?"
	// Специальный случай для Mr. Titspervert
	if node.Telegram == "@shpak_vv" || telegram == "shpak_vv" {
		avatar = "🖤"
	} else if node.Nickname != "" {
		avatar = strings.ToUpper(string([]rune(node.Nickname)[0]))
	} else if len(telegram) > 0 {
		avatar = strings.ToUpper(string(telegram[0]))
	}
	sb.WriteString(fmt.Sprintf(`<div class="tree-node-avatar">%s</div>`, avatar))

	// Информация
	sb.WriteString(`<div class="tree-node-info">`)
	sb.WriteString(fmt.Sprintf(`<div class="tree-node-name">%s</div>`, displayName))
	if dateStr != "" {
		sb.WriteString(fmt.Sprintf(`<div class="tree-node-date">%s</div>`, dateStr))
	}

	// Если участник исключен, показываем дату исключения
	if isExcluded {
		excludedDateStr := node.DateExcluded
		if parsedDate, err := time.Parse("1/2/2006", node.DateExcluded); err == nil {
			excludedDateStr = parsedDate.Format("02.01.2006")
		}
		sb.WriteString(fmt.Sprintf(`<div class="tree-node-excluded">❌ Исключен: %s</div>`, excludedDateStr))
	}
	sb.WriteString(`</div>`)

	// Статус индикатор
	statusIcon := "●"
	sb.WriteString(fmt.Sprintf(`<div class="tree-node-status">%s</div>`, statusIcon))

	sb.WriteString(`</div>`) // tree-node-content
	sb.WriteString(`</div>`) // tree-node

	// Если есть дети, рендерим их
	if len(node.Children) > 0 {
		// Сортируем детей по дате (формат mm/dd/yyyy)
		sort.Slice(node.Children, func(i, j int) bool {
			dateI, _ := time.Parse("1/2/2006", node.Children[i].DateInvited)
			dateJ, _ := time.Parse("1/2/2006", node.Children[j].DateInvited)
			return dateI.Before(dateJ)
		})

		sb.WriteString(`<div class="tree-children">`)
		for _, child := range node.Children {
			t.renderNode(sb, child, false)
		}
		sb.WriteString(`</div>`) // tree-children
	}

	sb.WriteString(`</div>`) // tree-branch
}
