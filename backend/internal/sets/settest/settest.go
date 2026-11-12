// Package settest gives tests the JSON copy of the learning sets and the recorded hidden tests.
package settest

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
)

// sets.json is the JSON copy of the learning sets. The test TestSetsEqualTheFile of package content keeps it equal to the Go content and writes it with -update.
//
//go:embed sets.json
var golden []byte

// Golden returns the learning sets of the file sets.json. A test uses them as the content of a server.
func Golden(t testing.TB) []sets.LearningSet {
	t.Helper()
	var content []sets.LearningSet
	decoder := json.NewDecoder(bytes.NewReader(golden))
	// A key that the model does not have is a gap of the model.
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&content); err != nil {
		t.Fatal(err)
	}
	for _, set := range content {
		for _, chapter := range set.Chapters {
			for i, card := range chapter.Cards {
				// The answer of a choice card and of a code card is the index of an option. JSON gives a float64.
				if index, ok := card.Answer.(float64); ok {
					chapter.Cards[i].Answer = int(index)
				}
			}
			for _, figure := range chapter.Figures {
				for _, frame := range figure.Frames {
					for _, row := range frame.Rows {
						for i, cell := range row.Cells {
							// A whole number in a cell is an int in the Go content.
							if number, ok := cell.(float64); ok && number == float64(int(number)) {
								row.Cells[i] = int(number)
							}
						}
					}
				}
			}
		}
	}
	return content
}
