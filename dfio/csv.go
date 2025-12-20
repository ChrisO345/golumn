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
	Header           bool
	Separator        rune
	IndexColumn      string
	SkipRows         []int
	TreatEmptyAsNull bool
	NullToken        string // if non-empty, token representing null in CSV
}

var defaultCSVSettings = CSVSettings{
	Header:           true,
	Separator:        ',',
	IndexColumn:      "",
	TreatEmptyAsNull: false,
	NullToken:        "",
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

	// helper to interpret null token
	nullToken := settings[0].NullToken
	for idx, val := range record {
		if !settings[0].Header {
			if (settings[0].TreatEmptyAsNull && val == "") || (nullToken != "" && val == nullToken) {
				se[idx].Append(nil)
			} else {
				se[idx].Append(val)
			}
		}
	}

	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(fmt.Errorf("error reading CSV: %v", err))
		}

		for jdx, val := range record {
			if (settings[0].TreatEmptyAsNull && val == "") || (nullToken != "" && val == nullToken) {
				se[jdx].Append(nil)
			} else {
				se[jdx].Append(val)
			}
		}
	}

	df := golumn.New(se...)
	return &df
}

// ToCSV writes a DataFrame to a CSV file at the provided path using the
// same CSVSettings used by FromCSV. If header is true, column names are
// written as the first row.
func ToCSV(path string, df *golumn.DataFrame, settings ...CSVSettings) error {
	cfg := defaultCSVSettings
	if len(settings) > 1 {
		return fmt.Errorf("only one settings struct allowed")
	}
	if len(settings) == 1 {
		cfg = settings[0]
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.Comma = cfg.Separator

	// write header
	if cfg.Header {
		names := df.Names()
		if err := w.Write(names); err != nil {
			return fmt.Errorf("error writing header: %w", err)
		}
	}

	// write rows
	nrows, ncols := df.Shape()
	for i := range nrows {
		rec := make([]string, ncols)
		for j := range ncols {
			// write null token if value is null
			col := df.Column(df.Names()[j])
			if col.IsNull(i) {
				if cfg.NullToken != "" {
					rec[j] = cfg.NullToken
				} else if cfg.TreatEmptyAsNull {
					rec[j] = ""
				} else {
					rec[j] = ""
				}
			} else {
				rec[j] = fmt.Sprint(df.At(i, j))
			}
		}
		if err := w.Write(rec); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("csv writer error: %w", err)
	}

	return nil
}
