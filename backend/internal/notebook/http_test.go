package notebook

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

// The expected messages are the messages of src/leetcode/http.ts.
const (
	badNumber = "Проверьте числовые поля."
	badChoice = "Выберите допустимое значение."
	badID     = "Неверный идентификатор записи."
	badDate   = "Укажите правильную дату."
	badFields = "Проверьте поля запроса."
	badText   = "Проверьте текстовые поля."
	notJSON   = "Отправьте запрос в формате JSON."
	badJSON   = "Неверный JSON в запросе."
	tooLarge  = "Файл или запрос слишком большой."
)

// checkInvalid fails the test unless err is the HTTP 400 invalid_input error with message. An empty message means no error.
func checkInvalid(t *testing.T, err error, message string) {
	t.Helper()
	if message == "" {
		if err != nil {
			t.Errorf("error = %v, want none", err)
		}
		return
	}
	var failure *HTTPError
	if !errors.As(err, &failure) || failure.Status != http.StatusBadRequest || failure.Code != "invalid_input" || failure.Message != message {
		t.Errorf("error = %#v, want 400 invalid_input %q", err, message)
	}
}

func TestInteger(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		min     int64
		max     int64
		want    int64
		wantErr string
	}{
		{name: "accepts the minimum", value: 1.0, min: 1, max: 50, want: 1},
		{name: "accepts the maximum", value: 50.0, min: 1, max: 50, want: 50},
		{name: "rejects a number below the minimum", value: 0.0, min: 1, max: 50, wantErr: badNumber},
		{name: "rejects a number above the maximum", value: 51.0, min: 1, max: 50, wantErr: badNumber},
		{name: "rejects a fraction", value: 1.5, min: 1, max: 50, wantErr: badNumber},
		{name: "rejects a string of digits", value: "1", min: 1, max: 50, wantErr: badNumber},
		{name: "rejects null", value: nil, min: 0, max: 50, wantErr: badNumber},
		{name: "rejects a boolean", value: true, min: 0, max: 50, wantErr: badNumber},
		{name: "accepts the largest safe integer", value: 9007199254740991.0, min: 1, max: maxSafeInteger, want: 9007199254740991},
		{name: "rejects the first integer above the safe range", value: 9007199254740992.0, min: 1, max: maxSafeInteger, wantErr: badNumber},
		{name: "accepts a negative number in the range", value: -3.0, min: -5, max: 5, want: -3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := integer(tt.value, tt.min, tt.max)
			checkInvalid(t, err, tt.wantErr)
			if got != tt.want {
				t.Errorf("integer = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestChoice(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		want    string
		wantErr string
	}{
		{name: "accepts an option", value: "medium", want: "medium"},
		{name: "rejects another string", value: "Medium", wantErr: badChoice},
		{name: "rejects a number", value: 1.0, wantErr: badChoice},
		{name: "rejects null", value: nil, wantErr: badChoice},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := choice(tt.value, "easy", "medium", "hard")
			checkInvalid(t, err, tt.wantErr)
			if got != tt.want {
				t.Errorf("choice = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr string
	}{
		{name: "accepts a lower-case identifier", value: "11111111-1111-4111-8111-11111111abcd"},
		{name: "accepts an upper-case identifier", value: "11111111-1111-4111-8111-11111111ABCD"},
		{name: "rejects a short identifier", value: "11111111-1111-4111-8111-11111111abc", wantErr: badID},
		{name: "rejects a letter after f", value: "11111111-1111-4111-8111-11111111abcg", wantErr: badID},
		{name: "rejects an identifier without hyphens", value: "111111111111411181111111111abcd0", wantErr: badID},
		{name: "rejects a line feed after the identifier", value: "11111111-1111-4111-8111-11111111abcd\n", wantErr: badID},
		{name: "rejects a number", value: 1.0, wantErr: badID},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := uuid(tt.value)
			checkInvalid(t, err, tt.wantErr)
			if tt.wantErr == "" && got != tt.value {
				t.Errorf("uuid = %q, want %q", got, tt.value)
			}
		})
	}
}

func TestCalendarDate(t *testing.T) {
	date := func(value string) *string { return &value }
	tests := []struct {
		name    string
		value   any
		want    *string
		wantErr string
	}{
		{name: "accepts null", value: nil, want: nil},
		{name: "accepts a date", value: "2026-10-07", want: date("2026-10-07")},
		{name: "accepts 29 February of a leap year", value: "2028-02-29", want: date("2028-02-29")},
		{name: "rejects 29 February of a common year", value: "2026-02-29", wantErr: badDate},
		{name: "rejects month 13", value: "2026-13-01", wantErr: badDate},
		{name: "rejects day zero", value: "2026-10-00", wantErr: badDate},
		{name: "rejects a month of one digit", value: "2026-1-07", wantErr: badDate},
		{name: "rejects a date with a time", value: "2026-10-07T00:00:00.000Z", wantErr: badDate},
		{name: "rejects an empty string", value: "", wantErr: badDate},
		{name: "rejects a number", value: 20261007.0, wantErr: badDate},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calendarDate(tt.value)
			checkInvalid(t, err, tt.wantErr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calendarDate = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantLimit  int64
		wantOffset int64
		wantErr    string
	}{
		{name: "uses 50 and 0 without parameters", query: "", wantLimit: 50, wantOffset: 0},
		{name: "reads both parameters", query: "limit=5&offset=100", wantLimit: 5, wantOffset: 100},
		{name: "accepts the limit 50", query: "limit=50", wantLimit: 50, wantOffset: 0},
		{name: "rejects the limit 51", query: "limit=51", wantErr: badNumber},
		{name: "rejects the limit 0", query: "limit=0", wantErr: badNumber},
		{name: "rejects an empty limit", query: "limit=", wantErr: badNumber},
		{name: "rejects a negative offset", query: "offset=-1", wantErr: badNumber},
		{name: "rejects an offset that is not a number", query: "offset=abc", wantErr: badNumber},
		{name: "rejects an offset above the safe range", query: "offset=9007199254740992", wantErr: badNumber},
		{name: "uses the first value of a repeated parameter", query: "offset=5&offset=7", wantLimit: 50, wantOffset: 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, err := url.ParseQuery(tt.query)
			if err != nil {
				t.Fatal(err)
			}
			limit, offset, err := pagination(query)
			checkInvalid(t, err, tt.wantErr)
			if limit != tt.wantLimit || offset != tt.wantOffset {
				t.Errorf("pagination = %d, %d; want %d, %d", limit, offset, tt.wantLimit, tt.wantOffset)
			}
		})
	}
}

func TestPage(t *testing.T) {
	next := func(value int64) *int64 { return &value }
	tests := []struct {
		name   string
		items  []string
		limit  int64
		offset int64
		want   Page[string]
	}{
		{name: "an extra item gives the next offset", items: []string{"a", "b", "c"}, limit: 2, offset: 4, want: Page[string]{Items: []string{"a", "b"}, NextOffset: next(6)}},
		{name: "a full page without an extra item is the last page", items: []string{"a", "b"}, limit: 2, offset: 4, want: Page[string]{Items: []string{"a", "b"}}},
		{name: "no rows give an empty list", items: nil, limit: 2, offset: 0, want: Page[string]{Items: []string{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := page(tt.items, tt.limit, tt.offset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("page = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestObject(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr string
	}{
		{name: "accepts the allowed keys", value: map[string]any{"username": "a", "password": "b"}},
		{name: "accepts an object without a key", value: map[string]any{"username": "a"}},
		{name: "rejects another key", value: map[string]any{"username": "a", "role": "parent"}, wantErr: badFields},
		{name: "rejects an array", value: []any{}, wantErr: badFields},
		{name: "rejects null", value: nil, wantErr: badFields},
		{name: "rejects a string", value: "username", wantErr: badFields},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := object(tt.value, "username", "password")
			checkInvalid(t, err, tt.wantErr)
			if tt.wantErr == "" && !reflect.DeepEqual(got, tt.value) {
				t.Errorf("object = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		max      int
		required bool
		wantErr  string
	}{
		{name: "accepts a string of the maximum length", value: "абвгд", max: 5},
		{name: "rejects a string above the maximum length", value: "абвгде", max: 5, wantErr: badText},
		{name: "counts a character outside the basic plane as two units", value: "😀😀😀", max: 5, wantErr: badText},
		{name: "accepts an empty string that is not required", value: "", max: 5},
		{name: "rejects an empty required string", value: "", max: 5, required: true, wantErr: badText},
		{name: "rejects a required string of spaces", value: " \t\n", max: 5, required: true, wantErr: badText},
		{name: "rejects a number", value: 5.0, max: 5, wantErr: badText},
		{name: "rejects null", value: nil, max: 5, wantErr: badText},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := text(tt.value, tt.max, tt.required)
			checkInvalid(t, err, tt.wantErr)
			if tt.wantErr == "" && got != tt.value {
				t.Errorf("text = %q, want %q", got, tt.value)
			}
		})
	}
}

func TestReadJSON(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		body        string
		want        any
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{name: "reads an object", contentType: "application/json", body: `{"a":1,"b":[null,"x"]}`, want: map[string]any{"a": 1.0, "b": []any{nil, "x"}}},
		{name: "accepts a content type with a parameter", contentType: "application/json; charset=utf-8", body: `"x"`, want: "x"},
		{name: "rejects plain text", contentType: "text/plain", body: `{}`, wantStatus: 400, wantCode: "invalid_input", wantMessage: notJSON},
		{name: "rejects a request without a content type", contentType: "", body: `{}`, wantStatus: 400, wantCode: "invalid_input", wantMessage: notJSON},
		{name: "rejects a content type in upper case", contentType: "Application/JSON", body: `{}`, wantStatus: 400, wantCode: "invalid_input", wantMessage: notJSON},
		{name: "rejects text that is not JSON", contentType: "application/json", body: `{`, wantStatus: 400, wantCode: "invalid_input", wantMessage: badJSON},
		{name: "rejects an empty body", contentType: "application/json", body: ``, wantStatus: 400, wantCode: "invalid_input", wantMessage: badJSON},
		{name: "rejects text after the value", contentType: "application/json", body: `{} {}`, wantStatus: 400, wantCode: "invalid_input", wantMessage: badJSON},
		{name: "rejects bytes that are not UTF-8", contentType: "application/json", body: "\"\xff\"", wantStatus: 400, wantCode: "invalid_input", wantMessage: badJSON},
		{name: "accepts a body of the maximum size", contentType: "application/json", body: `"` + strings.Repeat("x", 30) + `"`, want: strings.Repeat("x", 30)},
		{name: "rejects a body above the maximum size", contentType: "application/json", body: `"` + strings.Repeat("x", 31) + `"`, wantStatus: 413, wantCode: "too_large", wantMessage: tooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/leetcode/api/session", strings.NewReader(tt.body))
			if tt.contentType != "" {
				r.Header.Set("Content-Type", tt.contentType)
			}
			got, err := readJSON(r, 32)
			if tt.wantStatus == 0 {
				if err != nil {
					t.Fatalf("readJSON: %v", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("readJSON = %#v, want %#v", got, tt.want)
				}
				return
			}
			var failure *HTTPError
			if !errors.As(err, &failure) || failure.Status != tt.wantStatus || failure.Code != tt.wantCode || failure.Message != tt.wantMessage {
				t.Errorf("error = %#v, want %d %s %q", err, tt.wantStatus, tt.wantCode, tt.wantMessage)
			}
		})
	}
}

func TestReadBodyChecksTheSizeOfAStream(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{name: "accepts a stream of the maximum size", length: 32},
		{name: "rejects a stream above the maximum size", length: 33, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/leetcode/api/dictation", strings.NewReader(strings.Repeat("x", tt.length)))
			// A request without Content-Length reaches the size check of the stream.
			r.ContentLength = -1
			body, err := readBody(r, 32)
			var failure *HTTPError
			if tt.wantErr {
				if !errors.As(err, &failure) || failure.Status != 413 || failure.Code != "too_large" || failure.Message != tooLarge {
					t.Errorf("error = %#v, want 413 too_large", err)
				}
				return
			}
			if err != nil || len(body) != tt.length {
				t.Errorf("readBody = %d bytes, %v; want %d bytes", len(body), err, tt.length)
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantAllow  string
		wantBody   string
	}{
		{name: "an HTTP error gives its status, code and message", err: &HTTPError{Status: 401, Code: "unauthorized", Message: "Войдите, чтобы продолжить."}, wantStatus: 401, wantBody: `{"error":{"code":"unauthorized","message":"Войдите, чтобы продолжить."}}`},
		{name: "an HTTP error gives its fields", err: &HTTPError{Status: 409, Code: "draft_exists", Message: "Откройте текущий черновик, прежде чем начинать новую попытку.", Fields: map[string]string{"attemptId": "a1"}}, wantStatus: 409, wantBody: `{"error":{"code":"draft_exists","message":"Откройте текущий черновик, прежде чем начинать новую попытку.","fields":{"attemptId":"a1"}}}`},
		{name: "a method error gives the Allow header and a fixed message", err: methods(httptest.NewRequest(http.MethodPost, "/leetcode/login", nil), http.MethodGet, http.MethodHead), wantStatus: 405, wantAllow: "GET, HEAD", wantBody: `{"error":{"code":"method_not_allowed","message":"Этот метод недоступен."}}`},
		{name: "a statement that changed no row gives a conflict", err: changed(pgconn.NewCommandTag("UPDATE 0")), wantStatus: 409, wantBody: `{"error":{"code":"conflict","message":"Эта запись изменилась. Обновите страницу и повторите."}}`},
		{name: "another error gives 503 without its text", err: errors.New("connection refused"), wantStatus: 503, wantBody: `{"error":{"code":"unavailable","message":"Тетрадь недоступна. Сохраните работу и попробуйте снова."}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeError(w, httptest.NewRequest(http.MethodGet, "/leetcode/api/problems", nil), tt.err)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
			if got := w.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			if got := w.Header().Get("Allow"); got != tt.wantAllow {
				t.Errorf("Allow = %q, want %q", got, tt.wantAllow)
			}
			if got := w.Body.String(); got != tt.wantBody {
				t.Errorf("body = %s, want %s", got, tt.wantBody)
			}
		})
	}
}

func TestChangedAcceptsAChangedRow(t *testing.T) {
	if err := changed(pgconn.NewCommandTag("UPDATE 1")); err != nil {
		t.Errorf("changed = %v, want no error for one changed row", err)
	}
}

func TestMethodsAcceptsAnAllowedMethod(t *testing.T) {
	if err := methods(httptest.NewRequest(http.MethodHead, "/leetcode/login", nil), http.MethodGet, http.MethodHead); err != nil {
		t.Errorf("methods = %v, want no error for HEAD", err)
	}
}
