package main

import (
	"fmt"

	"github.com/chriso345/golumn"
	"github.com/chriso345/golumn/series"
)

func main() {
	left := golumn.New(
		series.New([]int{1, 2}, series.Int, "id"),
		series.New([]string{"A", "B"}, series.String, "name"),
	)
	right := golumn.New(
		series.New([]int{2, 3}, series.Int, "id"),
		series.New([]string{"X", "Y"}, series.String, "val"),
	)
	joined := left.Join(right, "id")
	fmt.Println("Join result:")
	fmt.Println(joined)

	df := golumn.New(
		series.New([]int{1, 1, 2}, series.Int, "id"),
		series.New([]string{"a", "b", "a"}, series.String, "key"),
		series.New([]int{10, 20, 30}, series.Int, "val"),
	)
	pivot := df.Pivot("id", "key", "val")
	fmt.Println("\nPivot result:")
	fmt.Println(pivot)

	unpivot := pivot.Unpivot([]string{"id"}, "key", "val")
	fmt.Println("\nUnpivot result:")
	fmt.Println(unpivot)
}
