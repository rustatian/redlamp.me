package notebook

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/jackc/pgx/v5"
)

// The system prompts of the tutor chat and of the evaluation of a solution.
var (
	//go:embed tutor-prompt.txt
	tutorPrompt string
	//go:embed evaluate-prompt.txt
	evaluatePrompt string
)

const (
	tutorModel = "claude-opus-5-5"
	// The thinking tokens count against this limit.
	tutorMaxTokens = 64_000
	// The request ends before the function timeout of 600 seconds.
	tutorTimeout = 9 * time.Minute
)

// TutorMessage is one message of a tutor chat.
type TutorMessage struct {
	ID        string  `json:"id"`
	Author    string  `json:"author"`
	Body      string  `json:"body"`
	Quote     *string `json:"quote"`
	CreatedAt string  `json:"createdAt"`
}

// tutorEvent is one server-sent event of an answer: a part of the text, the end, an error, or the sign that the model works.
type tutorEvent struct {
	Delta    string      `json:"delta,omitzero"`
	Done     bool        `json:"done,omitzero"`
	Error    *tutorError `json:"error,omitzero"`
	Thinking bool        `json:"thinking,omitzero"`
}

type tutorError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func tutorUnavailable() *HTTPError {
	return &HTTPError{Status: http.StatusServiceUnavailable, Code: "ai_unavailable", Message: "ИИ-репетитор сейчас недоступен. Попробуйте снова через несколько минут."}
}

func tutorInterrupted() *HTTPError {
	return &HTTPError{Status: http.StatusBadGateway, Code: "ai_interrupted", Message: "Ответ прервался. Отправьте сообщение ещё раз."}
}

func tutorSetup() *HTTPError {
	return &HTTPError{Status: http.StatusServiceUnavailable, Code: "setup_required", Message: "ИИ-репетитор не настроен. Попросите родителя проверить ключ API."}
}

// tutorFailure maps an error of the Anthropic API to the error of the notebook. A stream error event has the HTTP status 200, so the error type comes first.
// https://platform.claude.com/docs/en/api/errors
func tutorFailure(err error) *HTTPError {
	var failure *anthropic.Error
	if !errors.As(err, &failure) {
		// A network error or the timeout.
		return tutorUnavailable()
	}
	status := failure.StatusCode
	switch failure.Type() {
	case anthropic.ErrorTypeRateLimitError:
		status = http.StatusTooManyRequests
	case anthropic.ErrorTypeAuthenticationError, anthropic.ErrorTypeBillingError, anthropic.ErrorTypePermissionError:
		status = http.StatusUnauthorized
	case anthropic.ErrorTypeOverloadedError, anthropic.ErrorTypeAPIError, anthropic.ErrorTypeTimeoutError:
		status = http.StatusServiceUnavailable
	}
	switch {
	case status == http.StatusTooManyRequests:
		return &HTTPError{Status: http.StatusTooManyRequests, Code: "ai_usage_limit", Message: "Лимит запросов к ИИ-репетитору исчерпан. Попробуйте позже."}
	case status == http.StatusUnauthorized || status == http.StatusPaymentRequired || status == http.StatusForbidden:
		return tutorSetup()
	case status >= http.StatusInternalServerError:
		return tutorUnavailable()
	}
	return &HTTPError{Status: http.StatusBadGateway, Code: "ai_failed", Message: "ИИ-репетитор сейчас не может ответить. Попробуйте снова."}
}

// clip keeps the first max characters of a long text and marks the cut.
func clip(value string, max int) string {
	if start, cut := head(value, max); cut {
		return start + " [cut]"
	}
	return value
}

// said is the text of a student message for the model, with the selected text of the lesson.
func said(message string, quote *string) string {
	if quote == nil {
		return message
	}
	return "Selected text:\n> " + *quote + "\n\n" + message
}

var (
	// spaces matches the white space of a JavaScript regular expression.
	spaces        = regexp.MustCompile(`[\t\n\v\f\r \x{a0}\x{1680}\x{2000}-\x{200a}\x{2028}\x{2029}\x{202f}\x{205f}\x{3000}\x{feff}]+`)
	sectionStart  = regexp.MustCompile(`(?m)^## `)
	outlineHeader = regexp.MustCompile(`^#{1,4} `)
)

func squash(value string) string {
	return strings.Trim(spaces.ReplaceAllString(value, " "), " ")
}

// lessonText gives the outline of a lesson and its whole h2 sections from the one that contains the quote, while they fit the budget.
func lessonText(html string, quote *string) string {
	value := htmlText(html)
	var sections []string
	from := 0
	for _, at := range sectionStart.FindAllStringIndex(value, -1) {
		if at[0] > 0 {
			sections = append(sections, value[from:at[0]])
			from = at[0]
		}
	}
	sections = append(sections, value[from:])
	needle := ""
	if quote != nil {
		needle, _ = head(squash(*quote), 80)
	}
	start := max(0, slices.IndexFunc(sections, func(section string) bool { return strings.Contains(squash(section), needle) }))
	// Four lessons of up to 12,000 characters fit with their chapter headings.
	end, size := start+1, units(sections[start])
	for end < len(sections) && size+units(sections[end]) <= 50_000 {
		size += units(sections[end])
		end++
	}
	var outline []string
	for line := range strings.SplitSeq(value, "\n") {
		if outlineHeader.MatchString(line) {
			outline = append(outline, line)
		}
	}
	rest := ""
	if start > 0 || end < len(sections) {
		rest = " The rest is cut."
	}
	return fmt.Sprintf("Lesson outline:\n%s\n\nLesson text (sections %d-%d of %d).%s\n%s", clip(strings.Join(outline, "\n"), 3000), start+1, end, len(sections), rest, clip(strings.Join(sections[start:end], ""), 50_000))
}

func setNotFound() *HTTPError {
	return &HTTPError{Status: http.StatusNotFound, Code: "not_found", Message: "Такого набора нет."}
}

// tutorContext is the text that shows the model the page of the student: the learning set with its lessons, or the problem with its statement, homework, chapter and work.
func (s *Server) tutorContext(ctx context.Context, kind, id string, quote *string) (string, error) {
	if kind == "set" {
		set, ok := s.set(id)
		if !ok {
			return "", setNotFound()
		}
		var chapters, lessons []string
		for i, chapter := range set.Chapters {
			chapters = append(chapters, strconv.Itoa(i+1)+". "+chapter.Title+": "+chapter.Summary)
			for _, task := range chapter.Tasks {
				chapters = append(chapters, task.Title+" ("+task.Label()+", "+task.Difficulty+"). "+task.Note)
			}
			lessons = append(lessons, "<h1>Глава "+strconv.Itoa(i+1)+". "+chapter.Title+"</h1>\n"+chapter.Lesson)
		}
		return clip("Learning set: "+set.Title+"\nSummary: "+set.Summary+"\nTopics: "+strings.Join(set.Topics, ", ")+"\nChapters:\n"+strings.Join(chapters, "\n"), 4000) + "\n\n" + lessonText(strings.Join(lessons, "\n"), quote), nil
	}
	problem, err := getProblem(ctx, s.db, id)
	if err != nil {
		return "", err
	}
	codeforces := strings.HasPrefix(problem.Slug, "cf-")
	statement := ""
	if task, ok := s.setTask(problem.Slug); ok && task.Statement != "" {
		statement = task.Statement
	} else if !codeforces {
		// A problem for paid accounts has no content.
		statement = s.leetcodeContent(ctx, problem.Slug)
	}
	// The active homework task of the problem.
	var homework *struct {
		title, instructions, state string
		dueDate                    *string
	}
	var title, instructions, taskState string
	var dueDate *string
	err = s.db.QueryRow(ctx, "SELECT h.title,h.instructions,h.due_date,t.state FROM homework_tasks t JOIN homework h ON h.id=t.homework_id WHERE t.problem_id=$1 AND t.state NOT IN ('completed','cancelled') LIMIT 1", id).Scan(&title, &instructions, &dueDate, &taskState)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}
	if err == nil {
		homework = &struct {
			title, instructions, state string
			dueDate                    *string
		}{title, instructions, taskState, dueDate}
	}
	lines := []string{"Problem: ", "Topics: " + strings.Join(problem.Topics, ", "), "Summary: " + problem.Summary}
	if problem.Number != nil {
		lines[0] += strconv.FormatInt(*problem.Number, 10) + ". "
	}
	lines[0] += problem.Title + " (" + problem.Difficulty + ")"
	if homework != nil {
		due := "none"
		if homework.dueDate != nil {
			due = *homework.dueDate
		}
		lines = append(lines, "Homework: "+homework.title+". State: "+homework.state+". Due date: "+due+".")
	}
	parts := []string{clip(strings.Join(lines, "\n"), 4000)}
	if statement != "" {
		site := "LeetCode"
		if codeforces {
			site = "Codeforces"
		}
		parts = append(parts, site+" statement:\n"+clip(htmlText(statement), 6000))
	}
	if homework != nil {
		parts = append(parts, "Homework instructions:\n"+clip(homework.instructions, 2000))
	}
	if chapters := s.chaptersOf(problem.Slug); len(chapters) > 0 {
		ref := chapters[0]
		chapter, _ := s.chapter(ref.Set, ref.Number)
		parts = append(parts, "Chapter: "+ref.SetTitle+", глава "+strconv.FormatInt(ref.Number, 10)+". "+ref.Title+"\n"+clip(lessonText(chapter.Lesson, quote), 16_000))
	}
	// The draft comes before a saved attempt, and a newer saved attempt before an older one.
	var state, document string
	err = s.db.QueryRow(ctx, "SELECT state,document FROM attempts WHERE problem_id=$1 ORDER BY state,saved_at DESC LIMIT 1", id).Scan(&state, &document)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	}
	if err == nil {
		var work attemptDocument
		if err := json.Unmarshal([]byte(document), &work); err != nil {
			return "", err
		}
		texts := []string{"Notes: " + work.Notes}
		for i, a := range work.Approaches {
			texts = append(texts, "Approach "+strconv.Itoa(i+1)+": "+a.Label+"\nIdea: "+a.Idea+"\nWhy it works: "+a.Correctness+"\nTime complexity: "+a.TimeComplexity+"\nSpace complexity: "+a.SpaceComplexity+"\nEdge cases: "+a.EdgeCases+"\nMistakes: "+a.Mistakes+"\nLanguage: "+a.Language+"\nCode:\n"+a.Code)
		}
		name := "Newest saved attempt"
		if state == "draft" {
			name = "Current draft"
		}
		parts = append(parts, name+":\n"+clip(strings.Join(texts, "\n\n"), 12_000))
	}
	return strings.Join(parts, "\n\n"), nil
}

// tutorHistory returns the newest stored messages of a chat in the order of the chat: 20 messages or fewer, and 16,000 characters or fewer.
func (s *Server) tutorHistory(ctx context.Context, kind, id string) ([]anthropic.BetaMessageParam, error) {
	rows, err := s.db.Query(ctx, "SELECT author,body,quote FROM tutor_messages WHERE context_kind=$1 AND context_ref=$2 ORDER BY created_at DESC,id DESC LIMIT 20", kind, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var history []anthropic.BetaMessageParam
	student, size := false, 0
	for rows.Next() {
		var author, body string
		var quote *string
		if err := rows.Scan(&author, &body, &quote); err != nil {
			return nil, err
		}
		content := said(body, quote)
		if size += units(content); size > 16_000 {
			break
		}
		student = author == "student"
		if student {
			history = append(history, anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(content)))
		} else {
			history = append(history, anthropic.BetaMessageParam{OfAssistant: &anthropic.BetaAssistantMessageParam{Content: []anthropic.BetaContentBlockParamUnion{anthropic.NewBetaTextBlock(content)}}})
		}
	}
	// The Messages API takes a user turn first. The character limit can cut the history between a question and its answer.
	if len(history) > 0 && !student {
		history = history[:len(history)-1]
	}
	slices.Reverse(history)
	return history, rows.Err()
}

// thread checks the kind and the ID of a chat: a problem ID or the slug of a learning set.
func thread(kind, id any) (string, string, error) {
	value, err := choice(kind, "problem", "set")
	if err != nil {
		return "", "", err
	}
	if value == "set" {
		slug, err := text(id, 100, false)
		return value, slug, err
	}
	ref, err := uuid(id)
	return value, ref, err
}

func (s *Server) handleTutor(w http.ResponseWriter, r *http.Request, user *SessionUser) (bool, error) {
	if r.URL.Path != "/leetcode/api/tutor/messages" {
		return false, nil
	}
	return true, s.tutorRoute(w, r, user)
}

func (s *Server) tutorRoute(w http.ResponseWriter, r *http.Request, user *SessionUser) error {
	if err := methods(r, http.MethodGet, http.MethodPost); err != nil {
		return err
	}
	// The chats belong to the student: the parent can neither read nor write them.
	if err := requireRole(user, "student"); err != nil {
		return err
	}
	if r.Method == http.MethodPost {
		return s.ask(w, r)
	}
	ctx, query := r.Context(), r.URL.Query()
	limit, offset, err := pagination(query)
	if err != nil {
		return err
	}
	// A parameter that the query does not have is not a string.
	var kindValue, idValue any
	if query.Has("kind") {
		kindValue = query.Get("kind")
	}
	if query.Has("id") {
		idValue = query.Get("id")
	}
	kind, id, err := thread(kindValue, idValue)
	if err != nil {
		return err
	}
	if kind == "set" {
		if _, ok := s.set(id); !ok {
			return setNotFound()
		}
	} else if _, err := getProblem(ctx, s.db, id); err != nil {
		return err
	}
	rows, err := s.db.Query(ctx, "SELECT id,author,body,quote,created_at FROM tutor_messages WHERE context_kind=$1 AND context_ref=$2 ORDER BY created_at DESC,id DESC LIMIT $3 OFFSET $4", kind, id, limit+1, offset)
	if err != nil {
		return err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByPos[TutorMessage])
	if err != nil {
		return err
	}
	return writeJSON(w, r, http.StatusOK, page(items, limit, offset))
}

// ask sends a message of the student to the model and streams the answer to the browser as server-sent events.
func (s *Server) ask(w http.ResponseWriter, r *http.Request) error {
	data, err := readObject(r, "kind", "id", "message", "quote", "evaluate")
	if err != nil {
		return err
	}
	kind, id, err := thread(data["kind"], data["id"])
	if err != nil {
		return err
	}
	message, err := text(data["message"], 4000, true)
	if err != nil {
		return err
	}
	// The quote is null or a text. A request without the field is not valid.
	var quote *string
	if value, ok := data["quote"]; !ok || value != nil {
		selected, err := text(value, 2000, true)
		if err != nil {
			return err
		}
		quote = &selected
	}
	evaluate, asked := data["evaluate"]
	if asked && (evaluate != true || kind != "problem") {
		return invalid("Проверьте поля запроса.")
	}
	at := isoNow()
	block, err := s.tutorContext(r.Context(), kind, id, quote)
	if err != nil {
		return err
	}
	messages, err := s.tutorHistory(r.Context(), kind, id)
	if err != nil {
		return err
	}
	if s.tutor == nil {
		return tutorSetup()
	}
	system := tutorPrompt
	if asked {
		system = evaluatePrompt
	}
	// The context has text of the student and of other sites, so it has the user role: the rules of the system prompt stay above it.
	messages = append(messages, anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("<context>\n"+block+"\n</context>")), anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(said(message, quote))))
	// A browser disconnect cancels the request context. The answer still goes to the database, so the upstream request has its own context.
	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), tutorTimeout)
	defer cancel()

	// parts gets each part of the answer text. The empty text is the sign that the upstream stream is open.
	parts := make(chan string)
	var answer struct {
		text string
		// stop is the stop reason of a message that has its end event.
		stop   anthropic.BetaStopReason
		opened bool
		err    error
	}
	go func() {
		defer close(parts)
		// The thinking text can hold the full solution, so the display is "omitted" and only text deltas go to the browser.
		// A declined request gets an answer from the model that the API chooses for the decline.
		// https://platform.claude.com/docs/en/build-with-claude/refusals-and-fallback
		stream := s.tutor.Beta.Messages.NewStreaming(ctx, anthropic.BetaMessageNewParams{
			Model:        tutorModel,
			MaxTokens:    tutorMaxTokens,
			System:       []anthropic.BetaTextBlockParam{{Text: system}},
			Messages:     messages,
			Thinking:     anthropic.BetaThinkingConfigParamUnion{OfAdaptive: &anthropic.BetaThinkingConfigAdaptiveParam{Display: anthropic.BetaThinkingConfigAdaptiveDisplayOmitted}},
			OutputConfig: anthropic.BetaOutputConfigParam{Effort: anthropic.BetaOutputConfigEffortXhigh},
			Fallbacks:    anthropic.BetaFallbacksParamOfDefault(),
			Betas:        []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_07_01},
		})
		defer func() { _ = stream.Close() }()
		var text strings.Builder
		var stop anthropic.BetaStopReason
		for stream.Next() {
			if !answer.opened {
				answer.opened = true
				parts <- ""
			}
			switch event := stream.Current().AsAny().(type) {
			case anthropic.BetaRawContentBlockDeltaEvent:
				if delta, ok := event.Delta.AsAny().(anthropic.BetaTextDelta); ok && delta.Text != "" {
					text.WriteString(delta.Text)
					parts <- delta.Text
				}
			case anthropic.BetaRawMessageDeltaEvent:
				stop = event.Delta.StopReason
			case anthropic.BetaRawMessageStopEvent:
				answer.stop = stop
			}
		}
		answer.text, answer.err = text.String(), stream.Err()
	}()

	started, open := false, true
	start := func() {
		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_ = http.NewResponseController(w).Flush()
		started = true
	}
	// A browser disconnect fails the write. The loop reads the upstream to its end, so the turn is still stored.
	send := func(event tutorEvent) {
		if !started {
			start()
		}
		// An event always encodes.
		data, _ := stringify(event)
		if !open {
			return
		}
		if _, err := w.Write([]byte("data: " + string(data) + "\n\n")); err != nil {
			open = false
			return
		}
		_ = http.NewResponseController(w).Flush()
	}
	// CloudFront closes an origin response that is silent for 30 seconds, and the model can think longer before the first text.
	// https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/DownloadDistValuesOrigin.html
	timer := time.NewTimer(s.thinkingEvery)
	defer timer.Stop()
	for running := true; running; {
		select {
		case part, ok := <-parts:
			switch {
			case !ok:
				running = false
			case part == "":
				start()
			default:
				send(tutorEvent{Delta: part})
				timer.Reset(s.thinkingEvery)
			}
		case <-timer.C:
			send(tutorEvent{Thinking: true})
			timer.Reset(s.thinkingEvery)
		}
	}

	var failure *HTTPError
	var upstream *anthropic.Error
	switch {
	case answer.err == nil && answer.stop == anthropic.BetaStopReasonEndTurn:
		// A turn is the question and its complete answer.
		err := s.tx(ctx, func(tx pgx.Tx) error {
			const insert = "INSERT INTO tutor_messages (id,context_kind,context_ref,author,body,quote,created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
			if _, err := tx.Exec(ctx, insert, newID(), kind, id, "student", message, quote, at); err != nil {
				return err
			}
			_, err := tx.Exec(ctx, insert, newID(), kind, id, "assistant", answer.text, nil, isoNow())
			return err
		})
		if err == nil {
			send(tutorEvent{Done: true})
			return nil
		}
		log.Printf("tutor turn is not stored: %v", err)
		failure = tutorInterrupted()
	case errors.As(answer.err, &upstream), answer.err != nil && !answer.opened:
		failure = tutorFailure(answer.err)
	case answer.stop == anthropic.BetaStopReasonRefusal:
		failure = &HTTPError{Status: http.StatusBadGateway, Code: "ai_refused", Message: "ИИ-репетитор не может ответить на это сообщение. Измените вопрос и отправьте его снова."}
	default:
		// The stream ended before the end of the message, or the answer reached the token limit.
		failure = tutorInterrupted()
	}
	request := ""
	if upstream != nil {
		request = upstream.RequestID
	}
	log.Printf("tutor answer failed: code %s, stop reason %q, request %q", failure.Code, answer.stop, request)
	if !started {
		return failure
	}
	send(tutorEvent{Error: &tutorError{Code: failure.Code, Message: failure.Message}})
	return nil
}
