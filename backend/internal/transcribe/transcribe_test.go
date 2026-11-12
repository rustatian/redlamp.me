package transcribe

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming/types"
)

// fakeWriter is the audio side of a Transcribe stream. It records the size of each chunk.
type fakeWriter struct {
	chunks []int
	ended  chan struct{}
	once   sync.Once
}

func (w *fakeWriter) Send(_ context.Context, event types.AudioStream) error {
	w.chunks = append(w.chunks, len(event.(*types.AudioStreamMemberAudioEvent).Value.AudioChunk))
	return nil
}

func (w *fakeWriter) Close() error {
	w.once.Do(func() { close(w.ended) })
	return nil
}

func (w *fakeWriter) Err() error { return nil }

// fakeReader is the result side of a Transcribe stream.
type fakeReader struct {
	events chan types.TranscriptResultStream
	err    error
}

func (r *fakeReader) Events() <-chan types.TranscriptResultStream { return r.events }
func (r *fakeReader) Close() error                                { return nil }
func (r *fakeReader) Err() error                                  { return r.err }

// fakeStream makes a stream that gives the results after the end of the audio, as Transcribe gives the last results.
// The SDK has this constructor for tests: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/transcribestreaming#NewStartStreamTranscriptionEventStream
func fakeStream(results []types.TranscriptResultStream, err error) (*transcribestreaming.StartStreamTranscriptionEventStream, *fakeWriter) {
	w := &fakeWriter{ended: make(chan struct{})}
	r := &fakeReader{events: make(chan types.TranscriptResultStream), err: err}
	go func() {
		<-w.ended
		for _, result := range results {
			r.events <- result
		}
		close(r.events)
	}()
	return transcribestreaming.NewStartStreamTranscriptionEventStream(func(es *transcribestreaming.StartStreamTranscriptionEventStream) {
		es.Writer, es.Reader = w, r
	}), w
}

// result is a transcript event with one result.
func result(text string, partial bool) types.TranscriptResultStream {
	return &types.TranscriptResultStreamMemberTranscriptEvent{Value: types.TranscriptEvent{Transcript: &types.Transcript{Results: []types.Result{
		{IsPartial: partial, Alternatives: []types.Alternative{{Transcript: aws.String(text)}}},
	}}}}
}

// The audio goes to Transcribe in chunks of 100 ms, which are 3,200 bytes. The text is the final results with a space between them.
func TestExchange(t *testing.T) {
	tests := []struct {
		name       string
		audioBytes int
		results    []types.TranscriptResultStream
		streamErr  error
		wantChunks []int
		wantText   string
		wantErr    bool
	}{
		{name: "final results without the partial results", audioBytes: 7000,
			results:    []types.TranscriptResultStream{result("Основная", true), result("Основная идея", false), result("два", true), result("два указателя.", false)},
			wantChunks: []int{3200, 3200, 600}, wantText: "Основная идея два указателя."},
		{name: "silence of one chunk", audioBytes: 3200, wantChunks: []int{3200}, wantText: ""},
		{name: "event without a transcript", audioBytes: 2,
			results:    []types.TranscriptResultStream{&types.TranscriptResultStreamMemberTranscriptEvent{}, result("Да.", false)},
			wantChunks: []int{2}, wantText: "Да."},
		{name: "error of the stream", audioBytes: 2, results: []types.TranscriptResultStream{result("Да.", false)}, streamErr: errors.New("BadRequestException"),
			wantChunks: []int{2}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, writer := fakeStream(tt.results, tt.streamErr)
			text, err := exchange(context.Background(), stream, make([]byte, tt.audioBytes))
			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, want an error: %v", err, tt.wantErr)
			}
			if text != tt.wantText {
				t.Errorf("text = %q, want %q", text, tt.wantText)
			}
			if !slices.Equal(writer.chunks, tt.wantChunks) {
				t.Errorf("chunk sizes = %v, want %v", writer.chunks, tt.wantChunks)
			}
		})
	}
}
