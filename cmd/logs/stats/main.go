package main

import (
	"flag"
	"fmt"
	"log"

	"oldfartscounter/internal/logparser"
	"oldfartscounter/internal/output"
	"oldfartscounter/internal/stats"
)

var (
	dirFlag         = flag.String("dir", "logs", "Папка с логами (рекурсивно)")
	extFlag         = flag.String("ext", "", "Фильтр по расширению (например, .log). Пусто = все файлы")
	outCSV          = flag.String("out", "", "Сохранить CSV для матрицы убийств (опционально)")
	outHTML         = flag.String("html", "cs2_stats.html", "Путь к HTML (всегда пишется)")
	highlightPlayer = flag.String("highlight", "maslina420", "Игрок для золотой подсветки в табе 'Сорян, Братан'")
)

func main() {
	flag.Parse()

	// Создание компонентов
	parser := logparser.New()
	processor := stats.New()
	csvExporter := output.NewCSVExporter()
	htmlGenerator := output.NewHTMLGenerator()

	// Парсинг логов
	parseResult, err := parser.ParseDirectory(*dirFlag, *extFlag)
	if err != nil {
		log.Fatalf("ошибка парсинга логов: %v", err)
	}

	// Обработка статистики (всегда группируем по SteamID)
	statsData := processor.Process(parseResult)
	statsData.HighlightedPlayer = *highlightPlayer

	// Экспорт CSV (опционально)
	if *outCSV != "" {
		if err := csvExporter.WriteKillMatrix(*outCSV, statsData); err != nil {
			log.Fatalf("не удалось записать CSV: %v", err)
		}
		fmt.Printf("CSV сохранён: %s\n", *outCSV)
	}

	// Генерация HTML
	if err = htmlGenerator.Generate(*outHTML, statsData); err != nil {
		log.Fatalf("ошибка записи HTML: %v", err)
	}
	fmt.Printf("HTML сохранён: %s\n", *outHTML)
}
