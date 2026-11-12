package notebook

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LeetCode reads the data of one question from leetcode.com.
type LeetCode interface {
	// Question returns the JSON object with fields of the question titleSlug, or nil when LeetCode has no such question.
	// fields is a GraphQL selection, for example "title difficulty topicTags { name }".
	Question(ctx context.Context, titleSlug, fields string) (json.RawMessage, error)
}

// GraphQL is the LeetCode client for the GraphQL endpoint at URL, for example https://leetcode.com/graphql.
type GraphQL struct {
	URL string
}

// Question asks for one question and waits at most 5 seconds.
func (g GraphQL) Question(ctx context.Context, titleSlug, fields string) (json.RawMessage, error) {
	body, err := json.Marshal(map[string]any{
		"query":     "query question($titleSlug: String!) { question(titleSlug: $titleSlug) { " + fields + " } }",
		"variables": map[string]string{"titleSlug": titleSlug},
	})
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, g.URL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", "https://leetcode.com/problems/"+titleSlug+"/")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, fmt.Errorf("LeetCode returned %d", response.StatusCode)
	}
	var answer struct {
		Data *struct {
			Question json.RawMessage `json:"question"`
		} `json:"data"`
	}
	// A question with its statement and code snippets is below 1 MiB.
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&answer); err != nil {
		return nil, err
	}
	if answer.Data == nil {
		return nil, errors.New("LeetCode returned no data")
	}
	if string(answer.Data.Question) == "null" {
		return nil, nil
	}
	return answer.Data.Question, nil
}

// newID returns a random UUID of version 4: https://www.rfc-editor.org/rfc/rfc9562#name-uuid-version-4
func newID() string {
	var b [16]byte
	// rand.Read always fills b.
	_, _ = rand.Read(b[:])
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
