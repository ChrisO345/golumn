package main

import (
	"fmt"
	"github.com/chriso345/golumn"
	"github.com/chriso345/golumn/series"
)

func main() {
	// Create a DataFrame with two columns
	df := golumn.New(
		series.New([]string{"Alice", "Bob", "Charlie"}, series.String, "Name"),
		series.New([]int{25, 30, 35}, series.Int, "Age"),
	)

	fmt.Println("Basic DataFrame:")
	fmt.Println(df)

	// Filter rows where Age > 28
	filtered := df.Filter(func(row golumn.Row) bool {
		return row.Get("Age").(int) > 28
	})

	fmt.Println("\nFiltered (Age > 28):")
	fmt.Println(filtered)
}
