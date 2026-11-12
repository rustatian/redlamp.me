package notebook

// LearningSet is a learning set of the repository content: an ordered list of chapters.
type LearningSet struct {
	Slug     string
	Title    string
	Chapters []Chapter
}

// Chapter is one chapter of a learning set. Its number is its position in the set, from 1.
type Chapter struct {
	Title string
	Tasks []SetTask
}

// SetTask is a problem of a chapter. Number is nil for a problem without a number.
type SetTask struct {
	URL        string
	Number     *int64
	Title      string
	Difficulty string
	Topics     []string
	Note       string
}

// chapter returns the chapter number of the learning set setSlug.
func (s *Server) chapter(setSlug string, number int64) (Chapter, bool) {
	for _, set := range s.sets {
		if set.Slug == setSlug && number >= 1 && number <= int64(len(set.Chapters)) {
			return set.Chapters[number-1], true
		}
	}
	return Chapter{}, false
}
