// Package content has the learning sets of the notebook, in the order of the Бомжарня page.
package content

import (
	"maps"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
	"github.com/rustatian/redlamp.me/backend/internal/sets/complexity"
	"github.com/rustatian/redlamp.me/backend/internal/sets/heaps"
	"github.com/rustatian/redlamp.me/backend/internal/sets/permutations"
	"github.com/rustatian/redlamp.me/backend/internal/sets/twopointers"
)

// Sets are the learning sets.
var Sets = []sets.LearningSet{complexity.Set, permutations.Set, twopointers.Set, heaps.Set}

// References are the reference solutions of the tasks of all learning sets, by problem slug.
var References = merge(permutations.References, twopointers.References, heaps.References)

// Generators make the hidden tests of the tasks of all learning sets, by problem slug.
var Generators = merge(permutations.Generators, twopointers.Generators, heaps.Generators)

// merge joins the maps of the learning sets. A problem is in one learning set only.
func merge[V any](parts ...map[string]V) map[string]V {
	all := map[string]V{}
	for _, part := range parts {
		maps.Copy(all, part)
	}
	return all
}
