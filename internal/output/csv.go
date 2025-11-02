package output

import (
	"encoding/csv"
	"os"
	"strconv"

	"oldfartscounter/internal/stats"
)

// CSVExporter отвечает за экспорт в CSV
type CSVExporter struct{}

// NewCSVExporter создает новый экспортер CSV
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

// WriteKillMatrix записывает матрицу убийств в CSV файл
func (c *CSVExporter) WriteKillMatrix(path string, data *stats.StatsData) error {
	f, err := os.Create(path) // #nosec G304 - path is controlled by user input for CSV export
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Заголовок
	header := make([]string, 0, len(data.Players)+1)
	header = append(header, "Сорян, братан — Убийцы ↓ / Жертвы →")
	for _, player := range data.Players {
		header = append(header, player.Title)
	}
	_ = w.Write(header)

	// Строки данных
	for i, killer := range data.Players {
		row := make([]string, 0, len(data.Players)+1)
		row = append(row, killer.Title)
		for j := range data.Players {
			row = append(row, strconv.Itoa(data.KillMatrix.Matrix[i][j]))
		}
		_ = w.Write(row)
	}

	return w.Error()
}
