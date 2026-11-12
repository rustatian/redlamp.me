package notebook

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakeTranscriber is the speech service of the tests. It records the audio of each call.
type fakeTranscriber struct {
	text  string
	err   error
	calls [][]byte
}

func (f *fakeTranscriber) Transcribe(_ context.Context, pcm []byte) (string, error) {
	f.calls = append(f.calls, pcm)
	return f.text, f.err
}

// dictationPost sends a recording from the public origin with the content type.
func dictationPost(s *Server, cookie, contentType string, audio []byte) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, functionURL+"/leetcode/api/dictation", bytes.NewReader(audio))
	r.Header.Set("Origin", publicOrigin)
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Cookie", cookie)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	return w
}

// dictationServer makes a server with the speech service. A nil service is a server without one.
func dictationServer(t *testing.T, transcriber Transcriber) *Server {
	t.Helper()
	s, err := New(Config{DB: database(t), PublicOrigin: publicOrigin, Accounts: testAccounts, LeetCode: &fakeLeetCode{}, Transcriber: transcriber})
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// The cases of test/dictation.test.ts for a request that the server rejects, with the body rules of PCM audio. The speech service gets no call.
func TestDictationRejects(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		contentType string
		audio       []byte
		wantStatus  int
		wantCode    string
	}{
		{name: "parent account", role: "parent", contentType: "application/octet-stream", audio: make([]byte, 4), wantStatus: 403, wantCode: "forbidden"},
		{name: "no audio", role: "student", contentType: "application/octet-stream", audio: nil, wantStatus: 400, wantCode: "invalid_input"},
		{name: "half of a sample", role: "student", contentType: "application/octet-stream", audio: make([]byte, 3), wantStatus: 400, wantCode: "invalid_input"},
		{name: "JSON body of the Worker", role: "student", contentType: "application/json", audio: []byte(`{"audio":"AAAA"}`), wantStatus: 400, wantCode: "invalid_input"},
		{name: "one sample over two minutes", role: "student", contentType: "application/octet-stream", audio: make([]byte, 3_840_002), wantStatus: 413, wantCode: "too_large"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			speech := &fakeTranscriber{text: "text"}
			s := dictationServer(t, speech)
			w := dictationPost(s, loginAs(t, s, tt.role), tt.contentType, tt.audio)
			if w.Code != tt.wantStatus || errorCode(w) != tt.wantCode {
				t.Errorf("status %d with the code %q, want %d with %q; body %s", w.Code, errorCode(w), tt.wantStatus, tt.wantCode, w.Body.String())
			}
			if len(speech.calls) != 0 {
				t.Errorf("%d calls of the speech service, want none", len(speech.calls))
			}
		})
	}
}

// The cases of test/dictation.test.ts for a recording that the speech service accepts. Two minutes of audio are 3,840,000 bytes.
func TestDictationReturnsTheText(t *testing.T) {
	tests := []struct {
		name     string
		audio    []byte
		heard    string
		wantText string
	}{
		{name: "speech", audio: []byte{1, 0, 2, 0}, heard: " Основная идея – два указателя. ", wantText: "Основная идея – два указателя."},
		{name: "silence", audio: make([]byte, 32_000), heard: "", wantText: ""},
		{name: "the longest recording", audio: make([]byte, 3_840_000), heard: "Два указателя.", wantText: "Два указателя."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			speech := &fakeTranscriber{text: tt.heard}
			s := dictationServer(t, speech)
			w := dictationPost(s, loginAs(t, s, "student"), "application/octet-stream", tt.audio)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			if got := decode[map[string]string](t, w); len(got) != 1 || got["text"] != tt.wantText {
				t.Errorf("body = %v, want the text %q only", got, tt.wantText)
			}
			if len(speech.calls) != 1 || !bytes.Equal(speech.calls[0], tt.audio) {
				t.Errorf("the speech service did not get the audio of the request in one call")
			}
		})
	}
}

func TestDictationReportsAFailedTranscription(t *testing.T) {
	tests := []struct {
		name string
		// speech is nil for a server without a speech service.
		speech Transcriber
	}{
		{name: "error of the speech service", speech: &fakeTranscriber{err: errors.New("BadRequestException: the stream timed out")}},
		{name: "server without a speech service", speech: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := dictationServer(t, tt.speech)
			w := dictationPost(s, loginAs(t, s, "student"), "application/octet-stream", []byte{1, 0})
			want := map[string]map[string]string{"error": {"code": "dictation_unavailable", "message": "Распознавание речи сейчас недоступно. Попробуйте позже."}}
			if w.Code != http.StatusServiceUnavailable {
				t.Errorf("status = %d, want 503", w.Code)
			}
			reviewsEqualJSON(t, decode[any](t, w), want)
		})
	}
}
