//go:build race

package notebook

// raceDetector is true in a test binary with the race detector: https://go.dev/doc/articles/race_detector
// The detector makes the code many times slower, so a test does not check a run time with it.
const raceDetector = true
