package docgen

import (
	"os"
	"testing"
)

// TestMain is the root testing method
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestGen tests document generation
func TestGen(t *testing.T) {
	OutputExamplePdf("tests/test.pdf")
}
