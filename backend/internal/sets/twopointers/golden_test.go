package twopointers

import (
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/sets/settest"
)

// The references and the generators give the answers and the hidden tests of the Worker.
func TestReferencesAndGeneratorsEqualTheWorker(t *testing.T) {
	settest.Compare(t, Set, References, Generators)
}
