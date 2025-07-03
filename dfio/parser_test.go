package dfio

import (
	"testing"

	"github.com/chriso345/golumn/internal/testutils/assert"
)

func TestFromCSV(t *testing.T) {
	expected := "   first_name  last_name  username\n0         Rob       Pike       rob\n1         Ken   Thompson       ken\n2      Robert  Griesemer       gri"

	df := FromCSV("testdata/test.csv")

	r, c := df.Shape()
	assert.AssertEqual(t, r, 3)
	assert.AssertEqual(t, c, 3)
	assert.AssertEqual(t, df.String(), expected)
}
