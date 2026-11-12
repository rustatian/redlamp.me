package notebook

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/jackc/pgx/v5/pgxpool"
)

// The frames of a Messages API stream: https://platform.claude.com/docs/en/build-with-claude/streaming
const (
	tutorStart        = "event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"msg_1\",\"type\":\"message\",\"role\":\"assistant\",\"model\":\"claude-opus-5-5\",\"content\":[],\"stop_reason\":null,\"stop_sequence\":null,\"usage\":{\"input_tokens\":10,\"output_tokens\":1}}}\n\n"
	tutorPing         = "event: ping\ndata: {\"type\":\"ping\"}\n\n"
	tutorEnd          = "event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\",\"stop_sequence\":null},\"usage\":{\"output_tokens\":15}}\n\nevent: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"
	tutorTokenLimit   = "event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"max_tokens\",\"stop_sequence\":null},\"usage\":{\"output_tokens\":64000}}\n\nevent: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"
	tutorRefusal      = "event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"refusal\",\"stop_sequence\":null,\"stop_details\":{\"type\":\"refusal\",\"category\":\"cyber\",\"explanation\":\"Declined.\"}},\"usage\":{\"output_tokens\":1}}\n\nevent: message_stop\ndata: {\"type\":\"message_stop\"}\n\n"
	tutorOverloaded   = "event: error\ndata: {\"type\":\"error\",\"error\":{\"type\":\"overloaded_error\",\"message\":\"Overloaded\"}}\n\n"
	tutorRateLimited  = "event: error\ndata: {\"type\":\"error\",\"error\":{\"type\":\"rate_limit_error\",\"message\":\"Rate limited\"}}\n\n"
	tutorFallbackMark = "event: content_block_start\ndata: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"fallback\",\"from\":{\"model\":\"claude-opus-5-5\"},\"to\":{\"model\":\"claude-opus-5\"}}}\n\nevent: content_block_stop\ndata: {\"type\":\"content_block_stop\",\"index\":0}\n\n"
	// tutorThought is a thinking block with thinking text. The API sends no thinking text with the display "omitted". The server must not send it to the browser in any case.
	tutorThought = "event: content_block_start\ndata: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"thinking\",\"thinking\":\"\",\"signature\":\"\"}}\n\nevent: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"thinking_delta\",\"thinking\":\"The full solution is a hash map.\"}}\n\nevent: content_block_delta\ndata: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"signature_delta\",\"signature\":\"sig\"}}\n\nevent: content_block_stop\ndata: {\"type\":\"content_block_stop\",\"index\":0}\n\n"

	tutorSet = "permutations-and-combinations"
)

// tutorText is the frame of one part of the answer text.
func tutorText(text string) string {
	data, _ := json.Marshal(map[string]any{"type": "content_block_delta", "index": 1, "delta": map[string]any{"type": "text_delta", "text": text}})
	return "event: content_block_delta\ndata: " + string(data) + "\n\n"
}

// tutorReply is one answer of the fake upstream: an error with status and errorType, or a stream.
// A step of the stream is a string of frames or a time.Duration of silence. pause is the silence before the response starts.
type tutorReply struct {
	status    int
	errorType string
	pause     time.Duration
	steps     []any
}

// tutorSent is one request that the fake upstream got.
type tutorSent struct {
	path   string
	header http.Header
	body   map[string]any
}

// tutorUpstream is a fake of the Anthropic API. It gives each request the next reply and records the request.
type tutorUpstream struct {
	server  *httptest.Server
	mu      sync.Mutex
	replies []tutorReply
	sent    []tutorSent
}

func newTutorUpstream(t *testing.T, replies ...tutorReply) *tutorUpstream {
	t.Helper()
	u := &tutorUpstream{replies: replies}
	u.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		u.mu.Lock()
		u.sent = append(u.sent, tutorSent{path: r.URL.Path, header: r.Header.Clone(), body: body})
		var reply tutorReply
		if len(u.replies) > 0 {
			reply, u.replies = u.replies[0], u.replies[1:]
		} else {
			reply = tutorReply{status: http.StatusTeapot, errorType: "invalid_request_error"}
		}
		u.mu.Unlock()
		time.Sleep(reply.pause)
		w.Header().Set("request-id", "req_1")
		if reply.status != 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(reply.status)
			_ = json.NewEncoder(w).Encode(map[string]any{"type": "error", "error": map[string]any{"type": reply.errorType, "message": "Failed."}, "request_id": "req_1"})
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		for _, step := range reply.steps {
			switch value := step.(type) {
			case string:
				_, _ = io.WriteString(w, value)
				_ = http.NewResponseController(w).Flush()
			case time.Duration:
				time.Sleep(value)
			}
		}
	}))
	t.Cleanup(u.server.Close)
	return u
}

// requests returns the requests that the fake got.
func (u *tutorUpstream) requests() []tutorSent {
	u.mu.Lock()
	defer u.mu.Unlock()
	return slices.Clone(u.sent)
}

// tutorServer makes a notebook server with cfg and the client of the fake upstream. The client reads no credentials of the environment.
func tutorServer(t *testing.T, cfg Config, u *tutorUpstream) *Server {
	t.Helper()
	client := anthropic.NewClient(option.WithoutEnvironmentDefaults(), option.WithBaseURL(u.server.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0))
	cfg.PublicOrigin, cfg.Accounts, cfg.Anthropic = publicOrigin, testAccounts, &client
	if cfg.LeetCode == nil {
		cfg.LeetCode = &fakeLeetCode{}
	}
	s, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// tutorAnswer is the stream of the answer "Hi".
func tutorAnswer() tutorReply {
	return tutorReply{steps: []any{tutorStart + tutorText("Hi") + tutorEnd}}
}

// tutorSend is the body of a question about the learning set, with the fields of data in place of the defaults.
func tutorSend(data map[string]any) map[string]any {
	body := map[string]any{"kind": "set", "id": tutorSet, "message": "Explain this part.", "quote": nil}
	for key, value := range data {
		body[key] = value
	}
	return body
}

// tutorEvents returns the server-sent events of a response as decoded JSON values.
func tutorEvents(t *testing.T, w *httptest.ResponseRecorder) []any {
	t.Helper()
	events := []any{}
	for frame := range strings.SplitSeq(w.Body.String(), "\n\n") {
		if frame == "" {
			continue
		}
		var event any
		if err := json.Unmarshal([]byte(strings.TrimPrefix(frame, "data: ")), &event); err != nil {
			t.Fatalf("frame %q: %v", frame, err)
		}
		events = append(events, event)
	}
	return events
}

// tutorStored returns the stored messages, the message of the student first.
func tutorStored(t *testing.T, pool *pgxpool.Pool) []map[string]any {
	t.Helper()
	return problemsTableRows(t, pool, "SELECT context_kind AS kind,context_ref AS ref,author,body,quote FROM tutor_messages ORDER BY author DESC")
}

var (
	tutorInterruptedEvent = map[string]any{"error": map[string]any{"code": "ai_interrupted", "message": "Ответ прервался. Отправьте сообщение ещё раз."}}
	tutorUnavailableEvent = map[string]any{"error": map[string]any{"code": "ai_unavailable", "message": "ИИ-репетитор сейчас недоступен. Попробуйте снова через несколько минут."}}
	tutorUsageLimitEvent  = map[string]any{"error": map[string]any{"code": "ai_usage_limit", "message": "Лимит запросов к ИИ-репетитору исчерпан. Попробуйте позже."}}
	tutorRefusedEvent     = map[string]any{"error": map[string]any{"code": "ai_refused", "message": "ИИ-репетитор не может ответить на это сообщение. Измените вопрос и отправьте его снова."}}
)

// tutorTurn is the stored turn of the default question with the quote "Each circle is one call" and the answer.
func tutorTurn(answer string) []map[string]any {
	return []map[string]any{
		{"kind": "set", "ref": tutorSet, "author": "student", "body": "Explain this part.", "quote": "Each circle is one call"},
		{"kind": "set", "ref": tutorSet, "author": "assistant", "body": answer, "quote": nil},
	}
}

// The upstream stream becomes the events of the browser. Only a complete answer is stored.
func TestTutorAnswerStream(t *testing.T) {
	tests := []struct {
		name       string
		steps      []any
		wantEvents []any
		wantRows   []map[string]any
	}{
		{name: "complete answer", steps: []any{tutorStart + tutorPing, tutorText("Hi ") + tutorText("there") + tutorEnd},
			wantEvents: []any{map[string]any{"delta": "Hi "}, map[string]any{"delta": "there"}, map[string]any{"done": true}}, wantRows: tutorTurn("Hi there")},
		{name: "thinking text stays on the server", steps: []any{tutorStart + tutorThought + tutorText("Hi") + tutorEnd},
			wantEvents: []any{map[string]any{"delta": "Hi"}, map[string]any{"done": true}}, wantRows: tutorTurn("Hi")},
		{name: "answer of a fallback model", steps: []any{tutorStart + tutorFallbackMark + tutorText("Hi") + tutorEnd},
			wantEvents: []any{map[string]any{"delta": "Hi"}, map[string]any{"done": true}}, wantRows: tutorTurn("Hi")},
		{name: "stream that ends before the end of the message", steps: []any{tutorStart + tutorText("Hi ")},
			wantEvents: []any{map[string]any{"delta": "Hi "}, tutorInterruptedEvent}, wantRows: []map[string]any{}},
		{name: "answer at the token limit", steps: []any{tutorStart + tutorText("Hi ") + tutorTokenLimit},
			wantEvents: []any{map[string]any{"delta": "Hi "}, tutorInterruptedEvent}, wantRows: []map[string]any{}},
		{name: "refusal of the whole chain", steps: []any{tutorStart + tutorRefusal},
			wantEvents: []any{tutorRefusedEvent}, wantRows: []map[string]any{}},
		{name: "overload after text", steps: []any{tutorStart + tutorText("Hi ") + tutorOverloaded},
			wantEvents: []any{map[string]any{"delta": "Hi "}, tutorUnavailableEvent}, wantRows: []map[string]any{}},
		{name: "rate limit in the stream", steps: []any{tutorStart + tutorRateLimited},
			wantEvents: []any{tutorUsageLimitEvent}, wantRows: []map[string]any{}},
	}
	content := goldenSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			s := tutorServer(t, Config{DB: pool, Sets: content}, newTutorUpstream(t, tutorReply{steps: tt.steps}))
			w := request(t, s, "/tutor/messages", http.MethodPost, tutorSend(map[string]any{"quote": "Each circle is one call"}), loginAs(t, s, "student"))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			if got := w.Header().Get("Content-Type"); got != "text/event-stream; charset=utf-8" {
				t.Errorf("content type = %q, want text/event-stream; charset=utf-8", got)
			}
			if got := tutorEvents(t, w); !reflect.DeepEqual(got, tt.wantEvents) {
				t.Errorf("events = %v, want %v", got, tt.wantEvents)
			}
			reviewsEqualJSON(t, tutorStored(t, pool), tt.wantRows)
		})
	}
}

// An upstream failure before the stream gives the JSON error of the notebook. No message names the old upstream.
func TestTutorUpstreamErrors(t *testing.T) {
	tests := []struct {
		name string
		// reply is nil for an upstream that does not accept connections.
		reply *tutorReply
		// noKey is true for a server without an API key.
		noKey      bool
		wantStatus int
		wantCode   string
		wantCalls  int
	}{
		{name: "rate limit", reply: &tutorReply{status: 429, errorType: "rate_limit_error"}, wantStatus: 429, wantCode: "ai_usage_limit", wantCalls: 1},
		{name: "rejected key", reply: &tutorReply{status: 401, errorType: "authentication_error"}, wantStatus: 503, wantCode: "setup_required", wantCalls: 1},
		{name: "key without permission", reply: &tutorReply{status: 403, errorType: "permission_error"}, wantStatus: 503, wantCode: "setup_required", wantCalls: 1},
		{name: "account without credit", reply: &tutorReply{status: 402, errorType: "billing_error"}, wantStatus: 503, wantCode: "setup_required", wantCalls: 1},
		{name: "overload", reply: &tutorReply{status: 529, errorType: "overloaded_error"}, wantStatus: 503, wantCode: "ai_unavailable", wantCalls: 1},
		{name: "server error", reply: &tutorReply{status: 500, errorType: "api_error"}, wantStatus: 503, wantCode: "ai_unavailable", wantCalls: 1},
		{name: "request that the API does not accept", reply: &tutorReply{status: 400, errorType: "invalid_request_error"}, wantStatus: 502, wantCode: "ai_failed", wantCalls: 1},
		{name: "network error", wantStatus: 503, wantCode: "ai_unavailable", wantCalls: 0},
		{name: "server without an API key", reply: &tutorReply{}, noKey: true, wantStatus: 503, wantCode: "setup_required", wantCalls: 0},
	}
	content := goldenSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := database(t)
			var u *tutorUpstream
			if tt.reply != nil {
				u = newTutorUpstream(t, *tt.reply)
			} else {
				u = newTutorUpstream(t)
				u.server.Close()
			}
			s := tutorServer(t, Config{DB: pool, Sets: content}, u)
			if tt.noKey {
				s.tutor = nil
			}
			w := request(t, s, "/tutor/messages", http.MethodPost, tutorSend(nil), loginAs(t, s, "student"))
			body := decode[map[string]map[string]string](t, w)["error"]
			if w.Code != tt.wantStatus || body["code"] != tt.wantCode {
				t.Errorf("status %d with the code %q, want %d with %q; body %s", w.Code, body["code"], tt.wantStatus, tt.wantCode, w.Body.String())
			}
			if body["message"] == "" || strings.Contains(body["message"], "ChatGPT") {
				t.Errorf("message = %q, want a text without the name ChatGPT", body["message"])
			}
			if got := len(u.requests()); got != tt.wantCalls {
				t.Errorf("%d requests to the upstream, want %d", got, tt.wantCalls)
			}
			reviewsEqualJSON(t, tutorStored(t, pool), []map[string]any{})
		})
	}
}

// The request to the Messages API has the model, the effort, the thinking setup and the refusal fallback of the design.
func TestTutorRequest(t *testing.T) {
	pool := database(t)
	u := newTutorUpstream(t, tutorAnswer())
	s := tutorServer(t, Config{DB: pool, Sets: goldenSets(t)}, u)
	if w := request(t, s, "/tutor/messages", http.MethodPost, tutorSend(nil), loginAs(t, s, "student")); w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	sent := u.requests()
	if len(sent) != 1 {
		t.Fatalf("%d requests to the upstream, want 1", len(sent))
	}
	if sent[0].path != "/v1/messages" {
		t.Errorf("path = %q, want /v1/messages", sent[0].path)
	}
	// The SDK sends the API key as a bearer token: https://pkg.go.dev/github.com/anthropics/anthropic-sdk-go/v2/option#WithAPIKey
	if got := sent[0].header.Get("Authorization"); got != "Bearer test-key" {
		t.Errorf("authorization = %q, want the key of the client as a bearer token", got)
	}
	if got := sent[0].header.Values("Anthropic-Beta"); !slices.Contains(strings.Split(strings.Join(got, ","), ","), "server-side-fallback-2026-07-01") {
		t.Errorf("anthropic-beta = %q, want server-side-fallback-2026-07-01", got)
	}
	body := sent[0].body
	delete(body, "messages")
	reviewsEqualJSON(t, body, map[string]any{
		"model":         "claude-opus-5-5",
		"max_tokens":    64000,
		"stream":        true,
		"system":        []any{map[string]any{"type": "text", "text": tutorPrompt}},
		"thinking":      map[string]any{"type": "adaptive", "display": "omitted"},
		"output_config": map[string]any{"effort": "xhigh"},
		"fallbacks":     "default",
	})
}

// The browser gets a thinking event when the answer stream has no event for the time of the server. The times of the cases are short.
func TestTutorThinkingEvents(t *testing.T) {
	tests := []struct {
		name  string
		reply tutorReply
		// wantRest are the events without the thinking events.
		wantRest []any
	}{
		{name: "silence before the first text", reply: tutorReply{steps: []any{tutorStart, 200 * time.Millisecond, tutorText("Hi") + tutorEnd}},
			wantRest: []any{map[string]any{"delta": "Hi"}, map[string]any{"done": true}}},
		{name: "silence before the upstream answers, then an overload", reply: tutorReply{status: 529, errorType: "overloaded_error", pause: 200 * time.Millisecond},
			wantRest: []any{tutorUnavailableEvent}},
	}
	content := goldenSets(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tutorServer(t, Config{DB: database(t), Sets: content}, newTutorUpstream(t, tt.reply))
			s.thinkingEvery = 20 * time.Millisecond
			w := request(t, s, "/tutor/messages", http.MethodPost, tutorSend(nil), loginAs(t, s, "student"))
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200; body %s", w.Code, w.Body.String())
			}
			events := tutorEvents(t, w)
			thinking := map[string]any{"thinking": true}
			if len(events) == 0 || !reflect.DeepEqual(events[0], thinking) {
				t.Fatalf("events = %v, want a thinking event first", events)
			}
			rest := slices.DeleteFunc(events, func(event any) bool { return reflect.DeepEqual(event, thinking) })
			if !reflect.DeepEqual(rest, tt.wantRest) {
				t.Errorf("events without the thinking events = %v, want %v", rest, tt.wantRest)
			}
		})
	}
}

// tutorClosedBrowser is a response of a browser that closed the connection: each write fails.
type tutorClosedBrowser struct{ header http.Header }

func (b tutorClosedBrowser) Header() http.Header { return b.header }
func (b tutorClosedBrowser) WriteHeader(int)     {}
func (b tutorClosedBrowser) Write([]byte) (int, error) {
	return 0, errors.New("the browser disconnected")
}

func TestTutorStoresTheTurnAfterTheBrowserDisconnects(t *testing.T) {
	pool := database(t)
	s := tutorServer(t, Config{DB: pool, Sets: goldenSets(t)}, newTutorUpstream(t, tutorReply{steps: []any{tutorStart + tutorText("Hi "), tutorText("there") + tutorEnd}}))
	body, err := json.Marshal(tutorSend(map[string]any{"quote": "Each circle is one call"}))
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(http.MethodPost, functionURL+"/leetcode/api/tutor/messages", strings.NewReader(string(body)))
	r.Header.Set("Origin", publicOrigin)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Cookie", loginAs(t, s, "student"))
	s.ServeHTTP(tutorClosedBrowser{header: http.Header{}}, r)
	reviewsEqualJSON(t, tutorStored(t, pool), tutorTurn("Hi there"))
}

func TestTutorSendsTheInterruptedErrorWhenTheStoreFails(t *testing.T) {
	pool := database(t)
	s := tutorServer(t, Config{DB: pool, Sets: goldenSets(t)}, newTutorUpstream(t, tutorReply{steps: []any{tutorStart + tutorText("a") + tutorEnd}}))
	exec(t, pool, "CREATE FUNCTION fail_message() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'the store failed'; END $$")
	exec(t, pool, "CREATE TRIGGER fail_message BEFORE INSERT ON tutor_messages FOR EACH ROW EXECUTE FUNCTION fail_message()")
	w := request(t, s, "/tutor/messages", http.MethodPost, tutorSend(nil), loginAs(t, s, "student"))
	if got, want := tutorEvents(t, w), []any{map[string]any{"delta": "a"}, tutorInterruptedEvent}; !reflect.DeepEqual(got, want) {
		t.Errorf("events = %v, want %v", got, want)
	}
	reviewsEqualJSON(t, tutorStored(t, pool), []map[string]any{})
}
