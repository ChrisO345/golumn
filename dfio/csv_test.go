package dfio

import (
	"os"
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestReadCSV(t *testing.T) {
	expected := "   first_name  last_name  username\n0         Rob       Pike       rob\n1         Ken   Thompson       ken\n2      Robert  Griesemer       gri"

	df := FromCSV("testdata/test.csv")

	r, c := df.Shape()
	assert.Equal(t, r, 3)
	assert.Equal(t, c, 3)
	assert.Equal(t, df.String(), expected)
}

func TestFromCSV_HeaderFalseTemp(t *testing.T) {
	f, err := os.CreateTemp(".", "csvh_test")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	defer os.Remove(name)
	f.WriteString("x,y\n1,2\n3,4\n")
	f.Close()
	cfg := CSVSettings{Header: false, Separator: ','}
	df := FromCSV(name, cfg)
	// columns should be named Column 0, Column 1 and have 2 rows + header included
	assert.Equal(t, df.Columns()[0].Len(), 3)
}

func TestWriteCSV(t *testing.T) {
	df := FromCSV("testdata/test.csv")

	f, err := os.CreateTemp(".", "csv_write_")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	f.Close()
	defer os.Remove(name)

	if err := ToCSV(name, df); err != nil {
		t.Fatalf("ToCSV error: %v", err)
	}

	df2 := FromCSV(name)
	r, c := df2.Shape()
	assert.Equal(t, r, 3)
	assert.Equal(t, c, 3)
	assert.Equal(t, df.String(), df2.String())
}
