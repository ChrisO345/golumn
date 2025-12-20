package main

import (
	"fmt"

	"github.com/chriso345/golumn"
	"github.com/chriso345/golumn/series"
)

func main() {
	// Create a new DataFrame
	df := golumn.New(
		series.New([]string{"Alice", "Bob", "Charlie"}, series.String, "Name"),
		series.New([]int{25, 30, 35}, series.Int, "Age"),
	)

	// Print the DataFrame
	fmt.Println(df)

	// Add a new column
	df.Append(series.New([]string{"New York", "Los Angeles", "New York"}, series.String, "City"))

	// Print the updated DataFrame
	fmt.Println(df)

	// Filter where Age is greater than 28
	filtered := df.Filter(func(row golumn.Row) bool {
		return row.Get("Age").(int) > 28
	})

	// Print the filtered DataFrame
	fmt.Println(filtered)
}
