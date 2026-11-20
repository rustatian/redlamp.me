package notebook

import (
	"math"
	"time"
)

// startRating is the rating before the first contest.
const startRating = 800

// wrongPenalty is the time added to the finish time for each wrong submission before the first accepted one.
const wrongPenalty = 5 * time.Minute

// problemResult is the result of one problem of a finished contest.
// Accepted is the time from the start to the first accepted submission. Wrong is the number of wrong submissions before it.
type problemResult struct {
	Rating, Points int
	Solved         bool
	Accepted       time.Duration
	Wrong          int
}

// contestScore returns the sum of the points of the solved problems and the finish time:
// the latest first accepted submission plus the penalty of the wrong submissions of the solved problems.
// The finish time of a contest without a solved problem is 0.
func contestScore(problems []problemResult) (score int, finish time.Duration) {
	var latest time.Duration
	var wrong int
	for _, p := range problems {
		if !p.Solved {
			continue
		}
		score += p.Points
		latest = max(latest, p.Accepted)
		wrong += p.Wrong
	}
	return score, latest + time.Duration(wrong)*wrongPenalty
}

// ratingChange returns round(32 × Σ(S − E)) over the problems.
// S is 1 for a solved problem, else 0. E = 1 / (1 + 10^((R_problem − R) / 400)) is the expected score at the rating before.
func ratingChange(before int, problems []problemResult) int {
	var sum float64
	for _, p := range problems {
		expected := 1 / (1 + math.Pow(10, float64(p.Rating-before)/400))
		if p.Solved {
			sum += 1 - expected
		} else {
			sum -= expected
		}
	}
	return int(math.Round(32 * sum))
}

// ratingHistory returns the rating after each contest. The contests are in the order of their start.
func ratingHistory(contests [][]problemResult) []int {
	history := make([]int, len(contests))
	rating := startRating
	for i, problems := range contests {
		rating += ratingChange(rating, problems)
		history[i] = rating
	}
	return history
}
