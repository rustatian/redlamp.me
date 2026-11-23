package notebook

import (
	"encoding/json"
	"maps"
	"slices"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/rustatian/redlamp.me/backend/internal/bank"
	"github.com/rustatian/redlamp.me/backend/internal/editorials"
	"github.com/rustatian/redlamp.me/backend/internal/sets/content"
)

// The content checks of the Editorials.

// checkSolution is the start of each code block of an Editorial. A Sample runs with empty stdin, so a solution has no run button.
const checkSolution = `<pre class="no-run"><code class="language-python">`

// checkEditorial is the check of one Editorial. scope is the slugs of the learning set tasks and the bank problems.
func checkEditorial(t *testing.T, slug string, editorial editorials.Editorial, scope map[string]bool) {
	t.Helper()
	html := editorial.HTML
	if !scope[slug] {
		t.Error("no learning set task and no bank problem has this slug")
	}
	// The panel has the heading h2, so the sections of an Editorial are h3.
	var tags []string
	for _, m := range checkLessonTag.FindAllStringSubmatch(html, -1) {
		if !checkFragmentTag.MatchString(m[1]) || m[1] == "h2" {
			tags = append(tags, m[1])
		}
	}
	checkNone(t, "tags outside the list", tags)
	if checkBadFragment(html) {
		t.Error("the Editorial has a style or event attribute, or a div that is not a rule")
	}
	checkNone(t, "figure problems", checkFigures(html, editorial.Figures))
	var ids []string
	for _, m := range checkFigureID.FindAllStringSubmatch(html, -1) {
		ids = append(ids, m[1])
	}
	checkNone(t, "repeated figure IDs", checkRepeated(ids))
	if len(editorial.Figures) == 0 {
		t.Error("the Editorial has no figure")
	}
	solutions := strings.Count(html, checkSolution)
	if solutions == 0 || strings.Count(html, "<pre") != solutions {
		t.Errorf("%d code blocks start with %s, want all %d code blocks and at least 1", solutions, checkSolution, strings.Count(html, "<pre"))
	}
	if strings.Count(html, "Время:") < solutions || strings.Count(html, "Дополнительная память:") < solutions || !strings.Contains(html, "O(") {
		t.Errorf("want Время:, Дополнительная память: and O( for each of the %d solutions", solutions)
	}
	if length := len(utf16.Encode([]rune(htmlText(html)))); length >= 12000 {
		t.Errorf("the text has %d characters, want less than 12000", length)
	}
	// The web has no math renderer, and the lessons have no dashes. The rule also applies to each string of the figures.
	figures, err := json.Marshal(editorial.Figures)
	if err != nil {
		t.Fatal(err)
	}
	if strings.ContainsAny(html+string(figures), "$–—") {
		t.Error("the Editorial or its figures have a $ of LaTeX, an en dash or an em dash")
	}
	checkEnglishNames(t, editorial)
}

func TestEditorialContent(t *testing.T) {
	scope := map[string]bool{}
	for _, set := range content.Sets {
		for _, task := range set.Tasks() {
			scope[task.Slug()] = true
		}
	}
	for _, problem := range bank.Problems {
		scope[problem.Slug] = true
	}
	for _, slug := range slices.Sorted(maps.Keys(editorials.All)) {
		t.Run(slug, func(t *testing.T) { checkEditorial(t, slug, editorials.All[slug], scope) })
	}
}
