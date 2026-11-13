package notebook

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/rustatian/redlamp.me/backend/internal/sets"
	"github.com/rustatian/redlamp.me/backend/internal/sets/content"
)

// The content checks of the learning sets.

// checkSource is all learning sets of one source: "go".
type checkSource struct {
	name string
	all  []sets.LearningSet
}

func checkSources(t *testing.T) []checkSource {
	t.Helper()
	return []checkSource{{name: "go", all: content.Sets}}
}

// checkSet is one learning set of a source. name is the source and the slug, for example go/two-pointers.
type checkSet struct {
	name string
	all  []sets.LearningSet
	set  sets.LearningSet
}

func checkSets(t *testing.T) []checkSet {
	t.Helper()
	var cases []checkSet
	for _, source := range checkSources(t) {
		for _, set := range source.all {
			cases = append(cases, checkSet{name: source.name + "/" + set.Slug, all: source.all, set: set})
		}
	}
	return cases
}

// checkNone fails the test when the list of bad items is not empty.
func checkNone(t *testing.T, what string, bad []string) {
	t.Helper()
	if len(bad) > 0 {
		t.Errorf("%s: %q", what, bad)
	}
}

// checkRepeated returns each ID that an earlier item of the list already has.
func checkRepeated(ids []string) []string {
	var repeated []string
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			repeated = append(repeated, id)
		}
		seen[id] = true
	}
	return repeated
}

// checkSpace is the class of the JavaScript \s. The Go \s has only the ASCII spaces.
const checkSpace = `\t\n\v\f\r\p{Z}\x{FEFF}`

var (
	checkSlug         = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	checkMethod       = regexp.MustCompile(`^[a-z][A-Za-z]*$`)
	checkFragmentTag  = regexp.MustCompile(`^(?:h2|h3|p|ul|ol|li|table|caption|thead|tbody|tr|th|td|pre|code|strong|em|sup|blockquote|details|summary|a|br|div|figure)$`)
	checkStatementTag = regexp.MustCompile(`^(?:h3|p|ul|ol|li|table|thead|tbody|tr|th|td|pre|code|strong|em|sup|sub)$`)
	checkTag          = regexp.MustCompile(`<([^` + checkSpace + `/>]+)([^>]*)>`)
	// The JavaScript pattern is /<([a-z][a-z0-9]*)/gi.
	checkLessonTag = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)`)
	checkAttribute = regexp.MustCompile(`(?i)[` + checkSpace + `](?:style|on[a-z]+)=`)
	// The JavaScript pattern is /<div(?! class="rule"(?:>| data-exercise=))/i. Go has no lookahead, so a match without group 1 is a bad div.
	checkDiv         = regexp.MustCompile(`(?i)<div( class="rule"(?:>| data-exercise=))?`)
	checkFigure      = regexp.MustCompile(`<figure\b`)
	checkPlaceholder = regexp.MustCompile(`(?m)^<figure data-figure="[a-z0-9-]+"></figure>$`)
	checkFigureID    = regexp.MustCompile(`data-figure="([^"]+)"`)
	checkPythonToken = regexp.MustCompile(`"(?:[^"\\]|\\.)*"|\b(?:True|False|true|false|null)\b`)
	checkCyrillic    = regexp.MustCompile(`[а-яёА-ЯЁ]`)
	checkInteger     = regexp.MustCompile(`^-?\d+$`)
	// The JavaScript pattern has (?<!\p{L}) before the names and (?!\p{L}) after them.
	// Go has no lookaround, so the match takes the letter boundaries, and group 1 is the name.
	checkRussianName = regexp.MustCompile(`(?:^|\P{L})((?i:` +
		`куч(?:а|и|е|у|ей|ею|ам|ами|ах)?|дек(?:а|е|у|ом|и|ов|ам|ами|ах)?|стек(?:а|е|у|ом|и|ов|ам|ами|ах)?|массив\p{L}*|х[еэ]ш[- ]?таблиц\p{L}*|` +
		`(?:двоичн|бинарн)\p{L}*[` + checkSpace + `]+поиск\p{L}*|(?:два|двух|двумя)[` + checkSpace + `]+указател\p{L}*|скользящ\p{L}*[` + checkSpace + `]+окн\p{L}*|` +
		`динамическ\p{L}*[` + checkSpace + `]+программирован\p{L}*|сортировк\p{L}*|перестановк\p{L}*|подмножеств\p{L}*|подмассив\p{L}*|подпоследовательност\p{L}*|факториал\p{L}*` +
		`))(?:\P{L}|$)`)
)

// checkPython parses a Python literal of a test case, as pythonLiteral of test/sets.test.ts.
// The set files write Python literals in the JSON form with True and False. JSON true, false and null fail the parse.
func checkPython(text string) (any, error) {
	text = checkPythonToken.ReplaceAllStringFunc(text, func(token string) string {
		switch {
		case strings.HasPrefix(token, `"`):
			return token
		case token == "True":
			return "true"
		case token == "False":
			return "false"
		}
		return "?"
	})
	var value any
	err := json.Unmarshal([]byte(text), &value)
	return value, err
}

// checkBadTests reports a task without tests, with a wrong number of cases, with a mode that does not match its source, or with a wrong method name.
func checkBadTests(task sets.Task) bool {
	tests := task.Tests
	return tests == nil || len(tests.Cases) < 4 || len(tests.Cases) > 6 || (tests.Mode == "stdio") != strings.HasPrefix(task.Slug(), "cf-") || tests.Mode == "call" && !checkMethod.MatchString(tests.Method)
}

// checkBadCase reports a case without an expected answer, or with an input or an answer that the runner cannot read.
func checkBadCase(mode string, c sets.TestCase) bool {
	if c.Expected == nil {
		return true
	}
	if mode == "call" {
		input, err := checkPython(c.Input)
		if _, list := input.([]any); err != nil || !list {
			return true
		}
		_, err = checkPython(*c.Expected)
		return err != nil
	}
	return !strings.HasSuffix(c.Input, "\n") || strings.TrimSpace(*c.Expected) == ""
}

// checkBadStatement reports a statement without its example and constraints headings.
// The work window renders the statement as HTML, so a tag outside the list or any attribute fails.
func checkBadStatement(statement string) bool {
	if !strings.Contains(statement, "<h3>Пример 1</h3>") || !strings.Contains(statement, "<h3>Ограничения</h3>") {
		return true
	}
	return slices.ContainsFunc(checkTag.FindAllStringSubmatch(statement, -1), func(m []string) bool {
		return !checkStatementTag.MatchString(m[1]) || m[2] != ""
	})
}

// checkBadLesson reports a lesson without the lead paragraph first, with a style or event attribute, or with a div that is not a rule.
func checkBadLesson(lesson string) bool {
	if !strings.HasPrefix(lesson, `<p class="lead">`) || checkAttribute.MatchString(lesson) {
		return true
	}
	return slices.ContainsFunc(checkDiv.FindAllStringSubmatchIndex(lesson, -1), func(m []int) bool { return m[2] < 0 })
}

// checkBadRow reports a row without cells, a tree row with pointers or a window, a reversed window, or an index outside the cells.
func checkBadRow(row sets.FigureRow) bool {
	if len(row.Cells) == 0 || row.Tree && (row.Pointers != nil || row.Window != nil) || row.Window != nil && row.Window[0] > row.Window[1] {
		return true
	}
	indexes := slices.Concat(row.Marked, row.Dimmed)
	if row.Pointers != nil {
		for _, pointer := range *row.Pointers {
			indexes = append(indexes, pointer.Cell)
		}
	}
	if row.Window != nil {
		indexes = append(indexes, row.Window[:]...)
	}
	return slices.ContainsFunc(indexes, func(index int) bool { return index < 0 || index >= len(row.Cells) })
}

// checkBadFigure reports a figure with a wrong ID, an empty title, no frames, or a frame without a caption or rows or with a bad row.
func checkBadFigure(figure sets.Figure) bool {
	if !checkSlug.MatchString(figure.ID) || strings.TrimSpace(figure.Title) == "" || len(figure.Frames) == 0 {
		return true
	}
	return slices.ContainsFunc(figure.Frames, func(frame sets.FigureFrame) bool {
		return strings.TrimSpace(frame.Caption) == "" || len(frame.Rows) == 0 || slices.ContainsFunc(frame.Rows, checkBadRow)
	})
}

// checkSetContent is the check "$name content". all is the learning sets of the source of the set.
func checkSetContent(t *testing.T, all []sets.LearningSet, set sets.LearningSet) {
	t.Helper()
	same := 0
	for _, other := range all {
		if other.Slug == set.Slug {
			same++
		}
	}
	if same != 1 {
		t.Errorf("%d learning sets have the slug %q, want 1", same, set.Slug)
	}
	if !checkSlug.MatchString(set.Slug) {
		t.Errorf("slug %q has a wrong format", set.Slug)
	}
	if want := normalizeTopics(set.Topics, ""); !slices.Equal(set.Topics, want) {
		t.Errorf("topics = %q, want %q", set.Topics, want)
	}
	var slugs, topics, badSlugs, specs, cases, statements []string
	for _, task := range set.Tasks() {
		slug := task.Slug()
		slugs = append(slugs, slug)
		if !slices.Equal(task.Topics, normalizeTopics(task.Topics, slug)) {
			topics = append(topics, slug)
		}
		if !checkSlug.MatchString(slug) {
			badSlugs = append(badSlugs, slug)
		}
		if checkBadTests(task) {
			specs = append(specs, slug)
		}
		if task.Tests != nil {
			for _, c := range task.Tests.Cases {
				if checkBadCase(task.Tests.Mode, c) {
					cases = append(cases, slug+": "+c.Input)
				}
			}
		}
		if checkBadStatement(task.Statement) {
			statements = append(statements, slug)
		}
	}
	checkNone(t, "tasks with topics that are not normalized", topics)
	checkNone(t, "task slugs with a wrong format", badSlugs)
	checkNone(t, "repeated tasks", checkRepeated(slugs))
	checkNone(t, "tasks with wrong tests", specs)
	checkNone(t, "test cases with a wrong input or answer", cases)
	checkNone(t, "tasks with a wrong statement", statements)
	var blank, sizes, linkCounts, links, insecure, long, tags, lessons, placeholders, figures []string
	for i, chapter := range set.Chapters {
		number := fmt.Sprint(i + 1)
		if strings.TrimSpace(chapter.Title) == "" || strings.TrimSpace(chapter.Summary) == "" {
			blank = append(blank, chapter.Title)
		}
		least := 3
		if strings.Contains(chapter.Lesson, "data-exercise=") {
			least = 0
		}
		if len(chapter.Tasks) < least || len(chapter.Tasks) > 5 || len(chapter.Cards) < 6 || len(chapter.Cards) > 8 {
			sizes = append(sizes, chapter.Title)
		}
		if len(chapter.Links) < 3 || len(chapter.Links) > 5 {
			linkCounts = append(linkCounts, chapter.Title)
		}
		for _, link := range chapter.Links {
			links = append(links, link.URL)
			if u, err := url.Parse(link.URL); err != nil || u.Scheme != "https" {
				insecure = append(insecure, link.URL)
			}
		}
		// The length of a JavaScript string is its count of UTF-16 code units.
		if length := len(utf16.Encode([]rune(htmlText(chapter.Lesson)))); length >= 12000 {
			long = append(long, fmt.Sprintf("chapter %d: %d", i+1, length))
		}
		for _, m := range checkLessonTag.FindAllStringSubmatch(chapter.Lesson, -1) {
			if !checkFragmentTag.MatchString(m[1]) {
				tags = append(tags, m[1])
			}
		}
		if checkBadLesson(chapter.Lesson) {
			lessons = append(lessons, number)
		}
		// A figure is an empty placeholder on its own line. The placeholders follow the order of chapter.figures.
		if len(checkFigure.FindAllString(chapter.Lesson, -1)) != len(checkPlaceholder.FindAllString(chapter.Lesson, -1)) {
			placeholders = append(placeholders, number)
		}
		got, want := []string{}, []string{}
		for _, m := range checkFigureID.FindAllStringSubmatch(chapter.Lesson, -1) {
			got = append(got, m[1])
		}
		for _, figure := range chapter.Figures {
			want = append(want, figure.ID)
			if checkBadFigure(figure) {
				figures = append(figures, figure.ID)
			}
		}
		if !slices.Equal(got, want) {
			t.Errorf("chapter %d: figure placeholders %q, want %q", i+1, got, want)
		}
	}
	checkNone(t, "chapters without a title or a summary", blank)
	checkNone(t, "chapters with a wrong number of tasks or cards", sizes)
	checkNone(t, "chapters with a wrong number of links", linkCounts)
	checkNone(t, "links that are not https", insecure)
	checkNone(t, "repeated links", checkRepeated(links))
	checkNone(t, "lessons with 12000 or more characters of text", long)
	checkNone(t, "lesson tags outside the list", tags)
	checkNone(t, "chapters with a wrong lesson start, attribute or div", lessons)
	checkNone(t, "chapters with a figure that is not a placeholder", placeholders)
	checkNone(t, "bad figures", figures)
}

// checkChapterCosts is the check "$name states problem costs with English topics".
func checkChapterCosts(t *testing.T, chapter sets.Chapter) {
	t.Helper()
	var costs, topics []string
	for _, task := range chapter.Tasks {
		if !strings.Contains(task.Note, "Время:") || !strings.Contains(task.Note, "Память:") || !strings.Contains(task.Note, "O(") {
			costs = append(costs, task.Slug())
		}
		if slices.ContainsFunc(task.Topics, checkCyrillic.MatchString) || len(checkRepeated(task.Topics)) > 0 {
			topics = append(topics, task.Slug())
		}
	}
	checkNone(t, "tasks without the time and memory costs", costs)
	checkNone(t, "tasks with Russian or repeated topics", topics)
	if !strings.Contains(chapter.Lesson, "Дополнительная память") {
		t.Error("the lesson does not contain Дополнительная память")
	}
}

// checkStrings returns every string of a decoded JSON value, except the values of the url keys.
func checkStrings(value any) []string {
	var texts []string
	switch value := value.(type) {
	case string:
		texts = append(texts, value)
	case []any:
		for _, item := range value {
			texts = append(texts, checkStrings(item)...)
		}
	case map[string]any:
		for _, key := range slices.Sorted(maps.Keys(value)) {
			if key != "url" {
				texts = append(texts, checkStrings(value[key])...)
			}
		}
	}
	return texts
}

// checkRussianNames returns each match of checkRussianName in the text, as text.match of JavaScript with the g flag.
// The next search starts at the end of the name, so the letter boundary after a name can be the boundary before the next name.
func checkRussianNames(text string) []string {
	var names []string
	for start := 0; ; {
		m := checkRussianName.FindStringSubmatchIndex(text[start:])
		if m == nil {
			return names
		}
		names = append(names, text[start+m[2]:start+m[3]])
		start += m[3]
	}
}

// checkEnglishNames is the check "$name names the data structures in English".
// Technical names stay English. Ordinary Russian words such as a table row or a line of people remain valid.
func checkEnglishNames(t *testing.T, set sets.LearningSet) {
	t.Helper()
	data, err := json.Marshal(set)
	if err != nil {
		t.Fatal(err)
	}
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, text := range checkStrings(value) {
		names = append(names, checkRussianNames(text)...)
	}
	checkNone(t, "Russian names of data structures", names)
}

// checkDeck is the check "$name has consistent cards".
func checkDeck(t *testing.T, cards []sets.Card) {
	t.Helper()
	if len(cards) == 0 {
		t.Error("the deck has no cards")
	}
	var explain, answers, blanks, short []string
	for _, card := range cards {
		if strings.TrimSpace(card.Explain) == "" {
			explain = append(explain, card.ID)
		}
		if index, ok := card.Answer.(int); (card.Type == "choice" || card.Type == "code") && (!ok || index < 0 || index >= len(card.Options)) {
			answers = append(answers, card.ID)
		}
		// A code card shows code to read, or code with one blank to fill.
		if card.Type == "code" && strings.Count(card.Code, "___") > 1 {
			blanks = append(blanks, card.ID)
		}
		// The grader compares the whole answer string, so a short answer is one integer: "3, 1, 4" typed as "3,1,4" would be wrong.
		if answer, ok := card.Answer.(string); card.Type == "short" && (!ok || !checkInteger.MatchString(answer)) {
			short = append(short, card.ID)
		}
	}
	checkNone(t, "cards without an explanation", explain)
	checkNone(t, "cards with an answer outside the options", answers)
	checkNone(t, "code cards with more than one blank", blanks)
	checkNone(t, "short cards with an answer that is not one integer", short)
}

func TestSetChapterStatesCostsWithEnglishTopics(t *testing.T) {
	type chapterCase struct {
		name    string
		chapter sets.Chapter
	}
	var tests []chapterCase
	for _, c := range checkSets(t) {
		for i, chapter := range c.set.Chapters {
			tests = append(tests, chapterCase{name: fmt.Sprintf("%s/%d", c.name, i+1), chapter: chapter})
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { checkChapterCosts(t, tt.chapter) })
	}
}

func TestSetContent(t *testing.T) {
	for _, tt := range checkSets(t) {
		t.Run(tt.name, func(t *testing.T) { checkSetContent(t, tt.all, tt.set) })
	}
}

func TestSetNamesDataStructuresInEnglish(t *testing.T) {
	for _, tt := range checkSets(t) {
		t.Run(tt.name, func(t *testing.T) { checkEnglishNames(t, tt.set) })
	}
}

func TestSetDeckHasConsistentCards(t *testing.T) {
	type deckCase struct {
		name  string
		cards []sets.Card
	}
	var tests []deckCase
	for _, c := range checkSets(t) {
		for i, chapter := range c.set.Chapters {
			tests = append(tests, deckCase{name: fmt.Sprintf("%s-%d deck", c.name, i+1), cards: chapter.Cards})
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) { checkDeck(t, tt.cards) })
	}
}

func TestSetsUseEachProblemInOnlyOneChapter(t *testing.T) {
	for _, tt := range checkSources(t) {
		t.Run(tt.name, func(t *testing.T) {
			chapters := map[string]string{}
			var repeated []string
			for _, set := range tt.all {
				for i, chapter := range set.Chapters {
					place := fmt.Sprintf("%s/%d", set.Slug, i+1)
					for _, task := range chapter.Tasks {
						if first, ok := chapters[task.Slug()]; ok {
							repeated = append(repeated, fmt.Sprintf("%s in %s and %s", task.Slug(), first, place))
							continue
						}
						chapters[task.Slug()] = place
					}
				}
			}
			checkNone(t, "problems in more than one chapter", repeated)
		})
	}
}

func TestSetsHaveUniqueCardIDs(t *testing.T) {
	for _, tt := range checkSources(t) {
		t.Run(tt.name, func(t *testing.T) {
			var ids []string
			for _, set := range tt.all {
				for _, chapter := range set.Chapters {
					for _, card := range chapter.Cards {
						ids = append(ids, card.ID)
					}
				}
			}
			checkNone(t, "repeated card IDs", checkRepeated(ids))
		})
	}
}

func TestSetsHaveUniqueFigureIDs(t *testing.T) {
	for _, tt := range checkSources(t) {
		t.Run(tt.name, func(t *testing.T) {
			var ids []string
			for _, set := range tt.all {
				for _, chapter := range set.Chapters {
					for _, figure := range chapter.Figures {
						ids = append(ids, figure.ID)
					}
				}
			}
			checkNone(t, "repeated figure IDs", checkRepeated(ids))
		})
	}
}

// The cases of taskSlug and taskLabel of test/sets.test.ts. number is nil for a task without a number.
func TestTaskSlugAndLabel(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		number *int64
		slug   string
		label  string
	}{
		{name: "LeetCode task with a number", url: "https://leetcode.com/problems/subsets/", number: new(int64(78)), slug: "subsets", label: "LeetCode 78"},
		{name: "LeetCode task without a number", url: "https://leetcode.com/problems/subsets/", number: nil, slug: "subsets", label: "LeetCode"},
		{name: "Codeforces task", url: "https://codeforces.com/contest/1234/problem/B1", number: nil, slug: "cf-1234-b1", label: "Codeforces 1234B1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := sets.Task{URL: tt.url, Number: tt.number, Title: "Task", Difficulty: "easy", Topics: []string{}}
			if got := task.Slug(); got != tt.slug {
				t.Errorf("slug = %q, want %q", got, tt.slug)
			}
			if got := task.Label(); got != tt.label {
				t.Errorf("label = %q, want %q", got, tt.label)
			}
		})
	}
}
