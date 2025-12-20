package dfio

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chriso345/golumn"
	"github.com/chriso345/golumn/series"
)

// FromJSON reads a JSON file containing an array of objects and returns a DataFrame.
// All values are converted to strings similarly to FromCSV.
func FromJSON(path string) *golumn.DataFrame {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("error reading file: %v", err))
	}

	var arr []map[string]any
	if err := json.Unmarshal(data, &arr); err != nil {
		panic(fmt.Errorf("error unmarshalling json: %v", err))
	}

	if len(arr) == 0 {
		panic(fmt.Errorf("empty JSON array"))
	}

	// Collect column names from first object and sort for deterministic order
	names := make([]string, 0, len(arr[0]))
	for k := range arr[0] {
		names = append(names, k)
	}
	// simple selection sort for deterministic order
	for i := 0; i < len(names); i++ {
		min := i
		for j := i + 1; j < len(names); j++ {
			if names[j] < names[min] {
				min = j
			}
		}
		names[i], names[min] = names[min], names[i]
	}

	se := make([]series.Series, len(names))
	for idx, val := range names {
		se[idx] = series.NewEmptySeries(series.String, 0, val)
	}

	for _, obj := range arr {
		for j, name := range names {
			v, ok := obj[name]
			if !ok || v == nil {
				se[j].Append(nil)
				continue
			}
			se[j].Append(fmt.Sprint(v))
		}
	}

	df := golumn.New(se...)
	return &df
}

// ToJSON writes a DataFrame to a JSON file as an array of objects using column order from Names().
func ToJSON(path string, df *golumn.DataFrame) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer f.Close()

	names := df.Names()
	nrows, ncols := df.Shape()
	arr := make([]map[string]any, nrows)
	for i := range nrows {
		rec := make(map[string]any, ncols)
		for j, name := range names {
			rec[name] = df.At(i, j)
		}
		arr[i] = rec
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(arr); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}
