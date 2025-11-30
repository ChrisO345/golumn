package dfio

import (
	"os"
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestFromCSV(t *testing.T) {
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
