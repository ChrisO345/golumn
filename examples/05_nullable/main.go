package main

import (
	"fmt"

	"github.com/chriso345/golumn"
	"github.com/chriso345/golumn/series"
)

func main() {
	// Demonstrate nulls and cleaning
	s := series.New([]int{1, 2, 3}, series.Int, "vals")
	s.Elem(1).Set(nil)

	df := golumn.New(s, series.New([]string{"a", "b", "c"}, series.String, "letters"))

	fmt.Println("Original with null:")
	fmt.Println(df)

	df2 := df.Copy()

	// fill the "vals" column in-place
	df2.Column("vals").MutFillNA(99)
	fmt.Println("DataFrame after FillNA (mutating on column):")
	fmt.Println(df2)

	df3 := df.Copy()

	// drop rows with any nulls
	df3 = df3.DropNA()
	fmt.Println("DataFrame after DropNA (dropping rows with any nulls):")
	fmt.Println(df3)

}
