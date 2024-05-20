package initializers

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type GenderData struct {
	Name        string
	Gender      string
	Code        string
	Probability string
}

func InitializeCSV(filePath string) *[]GenderData {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read the header row: %v", err)
	}

	columnMap := createColumnMap(header)
	data := createGenderMap(reader, columnMap)

	log.Println("CSV initialized")

	return data
}

func createColumnMap(header []string) map[string]int {
	columnMap := make(map[string]int)
	for i, columnName := range header {
		columnMap[columnName] = i
	}
	return columnMap
}

func createGenderMap(reader *csv.Reader, columnMap map[string]int) *[]GenderData {
	data := make([]GenderData, 0)

	for {
		row, err := reader.Read()
		if err != nil {
			if err == csv.ErrFieldCount || err == io.EOF {
				break
			}
			log.Fatalf("Failed to read row: %v", err)
		}

		data = append(data, GenderData{
			Name:        row[columnMap["name"]],
			Gender:      row[columnMap["gender"]],
			Code:        row[columnMap["code"]],
			Probability: row[columnMap["wgt"]],
		})
	}

	return &data
}
