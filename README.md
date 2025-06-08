# golumn

`golumn` is a fast, type-safe, and easy-to-use dataframe library in Go. Built for data manipulation, transformation, and analysis, it provides a simple API for creating and manipulating columnar data structures.

> Golumn: sturdy columns for structured data.

---

## Features

- DataFrame and Series data structures
- Type-safe column operations
- Immutable-style transformations
- Fast, portable, and Go-native

---

## Installation

golumn is available on GitHub and can be installed using Go modules:

```bash
go get github.com/chriso345/golumn
```

---

## Usage

Here's a quick example of how to use `golumn`:

```go
package main

import (
  "fmt"
  "github.com/chriso345/golumn"
)

func main() {
	// Create a new DataFrame
	df := golumn.New(
		golumn.series.New([]string{"Alice", "Bob", "Charlie"}, golumn.series.String, "Name"),
		golumn.series.New([]int{25, 30, 35}, golumn.series.Int, "Age"),
	)

	// Print the DataFrame
	fmt.Println(df)

	// Add a new column
	df.Append(golumn.series.New([]string{"New York", "Los Angeles", "Chicago"}, golumn.series.String, "City"))

	// Print the updated DataFrame
	fmt.Println(df)
}
```

---

## Submodules

* **`golumn/series`**
  Core implementation of the `Series` type — a one-dimensional, type-safe, columnar data structure.

* **`golumn/io`** *(planned)*
  I/O utilities for loading and saving `DataFrame`s and `Series` in formats like CSV and JSON.

* **`golumn/math`** *(planned)*
  Statistical functions and numerical operations for both `Series` and `DataFrame` types.

* **`golumn/plot`** *(planned)*
  Minimal plotting tools for visual exploration of tabular data (e.g. line charts, histograms).

--- 

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
