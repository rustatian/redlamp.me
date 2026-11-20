package notebook

import (
	"slices"
	"testing"
	"time"
)

// The expected values are worked out from the rules in the design of the contests:
// the score is the sum of the points, the finish time is the latest first accepted submission plus 5 minutes for each wrong submission before it,
// and the rating change is round(32 × Σ(S − E)) with E = 1 / (1 + 10^((R_problem − R) / 400)).

func TestContestScore(t *testing.T) {
	const minute = time.Minute
	tests := []struct {
		name       string
		problems   []problemResult
		wantScore  int
		wantFinish time.Duration
	}{
		{
			name: "two solved with penalty",
			problems: []problemResult{
				{Rating: 800, Points: 3, Solved: true, Accepted: 12 * minute, Wrong: 1},
				{Rating: 1000, Points: 4, Solved: true, Accepted: 40 * minute},
				{Rating: 1000, Points: 4, Wrong: 3},
				{Rating: 1200, Points: 6},
			},
			wantScore:  7,
			wantFinish: 45 * minute,
		},
		{
			name: "nothing solved",
			problems: []problemResult{
				{Rating: 800, Points: 3, Wrong: 2},
				{Rating: 1000, Points: 4},
				{Rating: 1000, Points: 4},
				{Rating: 1200, Points: 6},
			},
		},
		{
			name: "hard solved first",
			problems: []problemResult{
				{Rating: 1200, Points: 6, Solved: true, Accepted: 30 * minute, Wrong: 2},
				{Rating: 800, Points: 3, Solved: true, Accepted: 10 * minute},
			},
			wantScore:  9,
			wantFinish: 40 * minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, finish := contestScore(tt.problems)
			if score != tt.wantScore || finish != tt.wantFinish {
				t.Errorf("contestScore = %d, %v, want %d, %v", score, finish, tt.wantScore, tt.wantFinish)
			}
		})
	}
}

// ratingProblems are the problems of a contest at the rating 800: the ratings 800, 1000, 1000 and 1200.
func ratingProblems(solved ...bool) []problemResult {
	ratings := []int{800, 1000, 1000, 1200}
	problems := make([]problemResult, len(solved))
	for i, s := range solved {
		problems[i] = problemResult{Rating: ratings[i], Solved: s}
	}
	return problems
}

func TestRatingChange(t *testing.T) {
	tests := []struct {
		name     string
		before   int
		problems []problemResult
		want     int
	}{
		{name: "easy and medium solved at 800", before: 800, problems: ratingProblems(true, true, false, false), want: 30},
		{name: "nothing solved at 800", before: 800, problems: ratingProblems(false, false, false, false), want: -34},
		{name: "one easy solved at 800", before: 800, problems: ratingProblems(true), want: 16},
		{name: "one easy unsolved at 800", before: 800, problems: ratingProblems(false), want: -16},
		{name: "one hard solved at 900", before: 900, problems: []problemResult{{Rating: 1200, Solved: true}}, want: 27},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ratingChange(tt.before, tt.problems); got != tt.want {
				t.Errorf("ratingChange = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestRatingHistory(t *testing.T) {
	tests := []struct {
		name     string
		contests [][]problemResult
		want     []int
	}{
		{name: "two contests", contests: [][]problemResult{ratingProblems(true, true, false, false), ratingProblems(false)}, want: []int{830, 813}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ratingHistory(tt.contests); !slices.Equal(got, tt.want) {
				t.Errorf("ratingHistory = %v, want %v", got, tt.want)
			}
		})
	}
}
