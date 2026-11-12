package notebook

import (
	"context"
	"log"
	"mime"
	"net/http"
	"strings"
)

// Transcriber changes the speech of a recording to text.
type Transcriber interface {
	// Transcribe returns the text of Russian speech. pcm is mono audio at 16 kHz with signed 16-bit little-endian samples. Silence gives no text.
	Transcribe(ctx context.Context, pcm []byte) (string, error)
}

// maxDictationBytes is two minutes of mono audio at 16 kHz with samples of two bytes.
const maxDictationBytes = 2 * 60 * 16_000 * 2

func (s *Server) handleDictation(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	if r.URL.Path != "/leetcode/api/dictation" {
		return false, nil
	}
	return true, s.dictation(w, r, user)
}

// dictation returns the text of a recording of the student. The server does not store the audio.
func (s *Server) dictation(w http.ResponseWriter, r *http.Request, user *SessionUser) error {
	if err := methods(r, http.MethodPost); err != nil {
		return err
	}
	if err := requireRole(user, "student"); err != nil {
		return err
	}
	if kind, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type")); kind != "application/octet-stream" {
		return invalid("Проверьте поля запроса.")
	}
	audio, err := readBody(r, maxDictationBytes)
	if err != nil {
		return err
	}
	// A recording has one sample or more, and a sample has two bytes.
	if len(audio) == 0 || len(audio)%2 != 0 {
		return invalid("Проверьте поля запроса.")
	}
	unavailable := &HTTPError{Status: http.StatusServiceUnavailable, Code: "dictation_unavailable", Message: "Распознавание речи сейчас недоступно. Попробуйте позже."}
	if s.transcriber == nil {
		return unavailable
	}
	text, err := s.transcriber.Transcribe(r.Context(), audio)
	if err != nil {
		log.Printf("dictation failed: %v", err)
		return unavailable
	}
	return writeJSON(w, r, http.StatusOK, map[string]string{"text": strings.TrimSpace(text)})
}
