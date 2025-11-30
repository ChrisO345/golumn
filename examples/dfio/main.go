package main

import (
	"fmt"
	"os"

	"github.com/chriso345/golumn/dfio"
)

func main() {
	// Read from CSV (testdata included in dfio package)
	df := dfio.FromCSV("../dfio/testdata/test.csv")
	fmt.Println("Read CSV:")
	fmt.Println(df)

	// Write to temp CSV and read back
	f, _ := os.CreateTemp(".", "example_csv_")
	name := f.Name()
	f.Close()
	defer os.Remove(name)
	dfio.ToCSV(name, df)
	df2 := dfio.FromCSV(name)
	fmt.Println("\nRoundtrip CSV matches:", df.String() == df2.String())

	// JSON roundtrip
	jf, _ := os.CreateTemp(".", "example_json_")
	jname := jf.Name()
	jf.Close()
	defer os.Remove(jname)
	fmt.Println("Writing JSON to:", jname)
	if err := dfio.ToJSON(jname, df); err != nil {
		fmt.Println("ToJSON error:", err)
	}
	dfj := dfio.FromJSON(jname)
	fmt.Println("JSON roundtrip matches:", df.String() == dfj.String())
}
