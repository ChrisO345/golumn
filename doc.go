// Package golumn provides high-performance, column-oriented data structures
// and tools for working with tabular data in Go.
//
// At its core, golumn is built around the concept of a DataFrame: a
// two-dimensional, columnar data structure where each column is a Series with a
// single, consistent data type.
//
// Each column in a DataFrame is a Series, and all Series in a DataFrame must be
// of equal length. Supported data types include int, float, string, bool, rune,
// and more.
//
// This package is suitable for applications requiring in-memory data
// manipulation, data analysis, or data cleaning in a statically typed language.
package golumn
