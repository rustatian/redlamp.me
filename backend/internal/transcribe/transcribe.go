// Package transcribe changes a recording to text with Amazon Transcribe streaming.
package transcribe

import (
	"context"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming"
	"github.com/aws/aws-sdk-go-v2/service/transcribestreaming/types"
)

// chunkBytes is 100 ms of mono audio at 16 kHz with samples of two bytes.
// The guide asks for chunks of 50 ms to 200 ms with an even number of bytes: https://docs.aws.amazon.com/transcribe/latest/dg/streaming.html
const chunkBytes = 3200

// Client is the speech service of the notebook.
type Client struct {
	api *transcribestreaming.Client
}

func New(cfg aws.Config) *Client {
	return &Client{api: transcribestreaming.NewFromConfig(cfg)}
}

// Transcribe returns the text of Russian speech. pcm is mono audio at 16 kHz with signed 16-bit little-endian samples.
func (c *Client) Transcribe(ctx context.Context, pcm []byte) (string, error) {
	output, err := c.api.StartStreamTranscription(ctx, &transcribestreaming.StartStreamTranscriptionInput{
		LanguageCode:         types.LanguageCodeRuRu,
		MediaEncoding:        types.MediaEncodingPcm,
		MediaSampleRateHertz: aws.Int32(16_000),
	})
	if err != nil {
		return "", err
	}
	return exchange(ctx, output.GetStream(), pcm)
}

// exchange sends the audio in chunks and joins the final results of the transcript.
func exchange(ctx context.Context, stream *transcribestreaming.StartStreamTranscriptionEventStream, pcm []byte) (string, error) {
	defer func() { _ = stream.Close() }()
	sent := make(chan error, 1)
	go func() {
		for chunk := range slices.Chunk(pcm, chunkBytes) {
			if err := stream.Send(ctx, &types.AudioStreamMemberAudioEvent{Value: types.AudioEvent{AudioChunk: chunk}}); err != nil {
				sent <- err
				return
			}
		}
		// The end of the audio makes Transcribe send the last results and end its stream.
		sent <- stream.Writer.Close()
	}()
	var parts []string
	for event := range stream.Events() {
		transcript, ok := event.(*types.TranscriptResultStreamMemberTranscriptEvent)
		if !ok || transcript.Value.Transcript == nil {
			continue
		}
		for _, result := range transcript.Value.Transcript.Results {
			// Transcribe sends a partial result again as a final result.
			if !result.IsPartial && len(result.Alternatives) > 0 {
				parts = append(parts, aws.ToString(result.Alternatives[0].Transcript))
			}
		}
	}
	if err := stream.Err(); err != nil {
		return "", err
	}
	if err := <-sent; err != nil {
		return "", err
	}
	return strings.Join(parts, " "), nil
}
