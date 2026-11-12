package notebook

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgconn"
)

// HTTPError is an error with the status, the code and the message of its JSON response.
type HTTPError struct {
	Status  int
	Code    string
	Message string
	Fields  map[string]string
}

func (e *HTTPError) Error() string { return e.Message }

func invalid(message string) error {
	return &HTTPError{Status: http.StatusBadRequest, Code: "invalid_input", Message: message}
}

// maxSafeInteger is the largest integer that a JSON number of the web client holds exactly.
const maxSafeInteger = 1<<53 - 1

// integer checks a decoded JSON value: a whole number from min to max.
func integer(value any, min, max int64) (int64, error) {
	number, ok := value.(float64)
	if !ok || number != math.Trunc(number) || math.Abs(number) > maxSafeInteger || number < float64(min) || number > float64(max) {
		return 0, invalid("Проверьте числовые поля.")
	}
	return int64(number), nil
}

func choice(value any, options ...string) (string, error) {
	s, ok := value.(string)
	if !ok || !slices.Contains(options, s) {
		return "", invalid("Выберите допустимое значение.")
	}
	return s, nil
}

var uuidPattern = regexp.MustCompile(`(?i)^[a-f0-9]{8}-(?:[a-f0-9]{4}-){3}[a-f0-9]{12}$`)

func uuid(value any) (string, error) {
	s, ok := value.(string)
	if !ok || !uuidPattern.MatchString(s) {
		return "", invalid("Неверный идентификатор записи.")
	}
	return s, nil
}

// calendarDate checks a decoded JSON value: null, or a date of the calendar in the form 2026-10-07.
func calendarDate(value any) (*string, error) {
	if value == nil {
		return nil, nil
	}
	s, ok := value.(string)
	if !ok {
		return nil, invalid("Укажите правильную дату.")
	}
	if _, err := time.Parse(time.DateOnly, s); err != nil {
		return nil, invalid("Укажите правильную дату.")
	}
	return &s, nil
}

// pagination reads the limit (1 to 50, default 50) and the offset (default 0) of a list request. Each is a decimal integer.
func pagination(query url.Values) (limit, offset int64, err error) {
	if limit, err = queryInteger(query, "limit", 50, 1, 50); err != nil {
		return 0, 0, err
	}
	if offset, err = queryInteger(query, "offset", 0, 0, maxSafeInteger); err != nil {
		return 0, 0, err
	}
	return limit, offset, nil
}

func queryInteger(query url.Values, name string, fallback, min, max int64) (int64, error) {
	if !query.Has(name) {
		return fallback, nil
	}
	number, err := strconv.ParseInt(query.Get(name), 10, 64)
	if err != nil {
		return 0, invalid("Проверьте числовые поля.")
	}
	return integer(float64(number), min, max)
}

// Page is one page of a list. NextOffset is null on the last page.
type Page[T any] struct {
	Items      []T    `json:"items"`
	NextOffset *int64 `json:"nextOffset"`
}

// page makes a page from the rows of a query that asked for limit+1 rows.
func page[T any](items []T, limit, offset int64) Page[T] {
	if int64(len(items)) > limit {
		next := offset + limit
		return Page[T]{Items: items[:limit], NextOffset: &next}
	}
	if items == nil {
		items = []T{}
	}
	return Page[T]{Items: items}
}

// changed returns the conflict error when a statement with a version check changed no row.
func changed(tag pgconn.CommandTag) error {
	if tag.RowsAffected() == 0 {
		return &HTTPError{Status: http.StatusConflict, Code: "conflict", Message: "Эта запись изменилась. Обновите страницу и повторите."}
	}
	return nil
}

// object checks a decoded JSON value: an object that has no key outside keys.
func object(value any, keys ...string) (map[string]any, error) {
	fields, ok := value.(map[string]any)
	if !ok {
		return nil, invalid("Проверьте поля запроса.")
	}
	for key := range fields {
		if !slices.Contains(keys, key) {
			return nil, invalid("Проверьте поля запроса.")
		}
	}
	return fields, nil
}

// text checks a decoded JSON value: a string of at most max UTF-16 code units, the length unit of the web client.
func text(value any, max int, required bool) (string, error) {
	s, ok := value.(string)
	if !ok || len(utf16.Encode([]rune(s))) > max || (required && strings.TrimSpace(s) == "") {
		return "", invalid("Проверьте текстовые поля.")
	}
	return s, nil
}

// methods returns the 405 error for a method outside allowed. Its message is the value of the Allow header.
func methods(r *http.Request, allowed ...string) error {
	if !slices.Contains(allowed, r.Method) {
		return &HTTPError{Status: http.StatusMethodNotAllowed, Code: "method_not_allowed", Message: strings.Join(allowed, ", ")}
	}
	return nil
}

func readBody(r *http.Request, maxBytes int64) ([]byte, error) {
	tooLarge := &HTTPError{Status: http.StatusRequestEntityTooLarge, Code: "too_large", Message: "Файл или запрос слишком большой."}
	if r.ContentLength > maxBytes {
		return nil, tooLarge
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBytes+1))
	if err != nil {
		return nil, err
	}
	if int64(len(body)) > maxBytes {
		return nil, tooLarge
	}
	return body, nil
}

// readJSON decodes a JSON request body into the types of encoding/json: map[string]any, []any, float64, string, bool and nil.
func readJSON(r *http.Request, maxBytes int64) (any, error) {
	mediaType, _, _ := strings.Cut(r.Header.Get("Content-Type"), ";")
	if strings.TrimSpace(mediaType) != "application/json" {
		return nil, invalid("Отправьте запрос в формате JSON.")
	}
	body, err := readBody(r, maxBytes)
	if err != nil {
		return nil, err
	}
	var value any
	if !utf8.Valid(body) || json.Unmarshal(body, &value) != nil {
		return nil, invalid("Неверный JSON в запросе.")
	}
	return value, nil
}

// send writes the status and the body. A HEAD response has no body.
func send(w http.ResponseWriter, r *http.Request, status int, body []byte) {
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_, _ = w.Write(body)
	}
}

// writeError writes err as the error JSON of the notebook. An error that is not an HTTPError is logged and gives 503.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	var failure *HTTPError
	if !errors.As(err, &failure) {
		log.Printf("notebook request failed: %s %s: %v", r.Method, r.URL.Path, err)
		failure = &HTTPError{Status: http.StatusServiceUnavailable, Code: "unavailable", Message: "Тетрадь недоступна. Сохраните работу и попробуйте снова."}
	}
	message := failure.Message
	if failure.Status == http.StatusMethodNotAllowed {
		w.Header().Set("Allow", failure.Message)
		message = "Этот метод недоступен."
	}
	type body struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Fields  map[string]string `json:"fields,omitzero"`
	}
	// A value of strings always encodes.
	data, _ := json.Marshal(map[string]body{"error": {Code: failure.Code, Message: message, Fields: failure.Fields}})
	w.Header().Set("Content-Type", "application/json")
	send(w, r, failure.Status, data)
}
