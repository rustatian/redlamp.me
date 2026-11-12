// Package content has the learning sets of the notebook, in the order of the Бомжарня page.
package content

import (
	"github.com/rustatian/redlamp.me/backend/internal/sets"
	"github.com/rustatian/redlamp.me/backend/internal/sets/complexity"
)

// Sets are the learning sets.
var Sets = []sets.LearningSet{complexity.Set}
