package initializers

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type GenderCountryData struct {
	Name        string
	Gender      string
	Code        string
	Probability string
}

type GenderData struct {
	Name   string
	Gender string
}

func InitializeGenderCountryCSV(filePath string) *[]GenderCountryData {
	return parseCSV(filePath, createGenderCountryMap)
}

func InitializeGenderCSV(filePath string) *[]GenderData {
	return parseCSV(filePath, createGenderMap)
}

func parseCSV[T any](filePath string, parseFn func(*csv.Reader, map[string]int) *[]T) *[]T {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Unable to read input file %s: %v", filePath, err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	header, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read the header row: %v", err)
	}

	columnMap := createColumnMap(header)
	return parseFn(reader, columnMap)
}

func createColumnMap(header []string) map[string]int {
	columnMap := make(map[string]int)
	for i, columnName := range header {
		columnMap[columnName] = i
	}
	return columnMap
}

func createGenderCountryMap(reader *csv.Reader, columnMap map[string]int) *[]GenderCountryData {
	data := make([]GenderCountryData, 0)

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Failed to read row: %v", err)
		}

		data = append(data, GenderCountryData{
			Name:        row[columnMap["name"]],
			Gender:      row[columnMap["gender"]],
			Code:        row[columnMap["code"]],
			Probability: row[columnMap["wgt"]],
		})
	}

	return &data
}

func createGenderMap(reader *csv.Reader, columnMap map[string]int) *[]GenderData {
	data := make([]GenderData, 0)

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Failed to read row: %v", err)
		}

		data = append(data, GenderData{
			Name:   row[columnMap["name"]],
			Gender: row[columnMap["gender"]],
		})
	}

	return &data
}
