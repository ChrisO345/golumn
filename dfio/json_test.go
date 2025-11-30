package dfio

import (
	"os"
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestReadJSON(t *testing.T) {
	df := FromJSON("testdata/test.json")
	r, c := df.Shape()
	assert.Equal(t, r, 3)
	assert.Equal(t, c, 3)

	assert.Equal(t, df.At(0, 0), "Rob")
	assert.Equal(t, df.At(0, 1), "Pike")
	assert.Equal(t, df.At(0, 2), "rob")

	assert.Equal(t, df.At(1, 0), "Ken")
	assert.Equal(t, df.At(1, 1), "Thompson")
	assert.Equal(t, df.At(1, 2), "ken")

	assert.Equal(t, df.At(2, 0), "Robert")
	assert.Equal(t, df.At(2, 1), "Griesemer")
	assert.Equal(t, df.At(2, 2), "gri")
}

func TestWriteJSON(t *testing.T) {
	df := FromCSV("testdata/test.csv")

	f, err := os.CreateTemp(".", "json_write_")
	if err != nil {
		t.Fatal(err)
	}
	name := f.Name()
	f.Close()
	defer os.Remove(name)

	if err := ToJSON(name, df); err != nil {
		t.Fatalf("ToJSON error: %v", err)
	}

	df2 := FromJSON(name)
	r, c := df2.Shape()
	assert.Equal(t, r, 3)
	assert.Equal(t, c, 3)

	for i := range r {
		for j := range c {
			assert.Equal(t, df.At(i, j), df2.At(i, j))
		}
	}
}
