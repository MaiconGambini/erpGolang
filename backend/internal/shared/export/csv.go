package export

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
)

// UTF-8 BOM so Excel opens accented pt-BR text correctly.
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

func WriteCSV(w http.ResponseWriter, filename string, header []string, rows [][]string) error {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, sanitizeFilename(filename)))
	if _, err := w.Write(utf8BOM); err != nil {
		return err
	}
	writer := csv.NewWriter(w)
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, `"`, "")
	name = strings.ReplaceAll(name, "\n", "")
	if name == "" {
		return "export.csv"
	}
	return name
}
