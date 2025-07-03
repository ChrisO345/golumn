package dfio

import (
	"encoding/csv"
	"fmt"
	"github.com/chriso345/golumn"
	"github.com/chriso345/golumn/series"
	"io"
	"os"
)

// CSVSettings defines a struct that contains settings for reading a CSV file, allows for optional settings
type CSVSettings struct {
	Header      bool
	Separator   rune
	IndexColumn string
	SkipRows    []int
}

var defaultCSVSettings = CSVSettings{
	Header:      true,
	Separator:   ',',
	IndexColumn: "",
}

// FromCSV reads a CSV file and returns a DataFrame
func FromCSV(path string, settings ...CSVSettings) *golumn.DataFrame {
	if len(settings) == 0 {
		settings = append(settings, defaultCSVSettings)
	} else if len(settings) > 1 {
		fmt.Println(fmt.Errorf("only one settings struct allowed"))
		os.Exit(1)
	}

	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("error opening file: %v", err))
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(fmt.Errorf("error closing file: %v", err))
		}
	}(file)

	reader := csv.NewReader(file)
	reader.Comma = settings[0].Separator

	record, err := reader.Read()
	if err == io.EOF {
		panic(fmt.Errorf("empty CSV file"))
	}

	if err != nil {
		panic(fmt.Errorf("error reading CSV: %v", err))
	}

	names := record
	if !settings[0].Header {
		names = make([]string, len(record))
		for idx := range record {
			names[idx] = fmt.Sprintf("Column %d", idx)
		}
	}

	se := make([]series.Series, len(record))
	for idx, val := range names {
		se[idx] = series.NewEmptySeries(series.String, 0, val)
	}

	if !settings[0].Header {
		for idx, val := range record {
			se[idx].Append(val)
		}
	}

	idx := 0
	for {
		idx++
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(fmt.Errorf("error reading CSV: %v", err))
		}

		for jdx, val := range record {
			se[jdx].Append(val)
		}
	}

	df := golumn.New(se...)
	return &df
}
