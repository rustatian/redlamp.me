package notebook

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// The cases of test/homework.test.ts and test/homework-content.test.ts that read /dashboard or /homework/{id}/sets.

// homeworkDashboard is the part of the Dashboard of GET /dashboard that the homework tests read.
type homeworkDashboard struct {
	Homework []map[string]any `json:"homework"`
}

func TestHomeworkReturnCompletedTaskOfArchivedProblem(t *testing.T) {
	f := homeworkSetup(t)
	submission := homeworkSolve(t, f)
	if w := request(t, f.s, "/tasks/"+f.taskID+"/review", http.MethodPost, map[string]any{"version": 2, "submissionId": submission, "decision": "completed", "body": ""}, f.parent); w.Code != http.StatusOK {
		t.Errorf("completed review: status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	if w := request(t, f.s, "/problems/"+f.problemID+"/archive", http.MethodPost, map[string]any{"version": 1}, f.parent); w.Code != http.StatusOK {
		t.Errorf("archive: status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	review := func() *httptest.ResponseRecorder {
		return request(t, f.s, "/tasks/"+f.taskID+"/review", http.MethodPost, map[string]any{"version": 3, "submissionId": submission, "decision": "changes_requested", "body": "Try again."}, f.parent)
	}
	rejected := review()
	if rejected.Code != http.StatusConflict {
		t.Errorf("review of the archived problem: status = %d, want 409", rejected.Code)
	}
	problemsMatchBody(t, decode[any](t, rejected), map[string]any{"error": map[string]any{"code": "archived"}})
	if got, want := homeworkRow(t, f.pool, "SELECT state,version FROM homework_tasks WHERE id=$1", f.taskID), map[string]any{"state": "completed", "version": int64(3)}; !reflect.DeepEqual(got, want) {
		t.Errorf("task = %v, want %v", got, want)
	}
	if got := queryValue[int64](t, f.pool, "SELECT count(*) AS count FROM feedback"); got != 1 {
		t.Errorf("feedback count = %d, want 1", got)
	}
	if w := request(t, f.s, "/problems/"+f.problemID+"/restore", http.MethodPost, map[string]any{"version": 2}, f.parent); w.Code != http.StatusOK {
		t.Errorf("restore: status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	if w := review(); w.Code != http.StatusOK {
		t.Errorf("review of the restored problem: status = %d, want 200; body %s", w.Code, w.Body.String())
	}
	dashboard := decode[map[string]any](t, request(t, f.s, "/dashboard?today=2026-10-07", http.MethodGet, nil, f.student))
	problemsMatchBody(t, dashboard["homework"], []any{map[string]any{"id": f.taskID, "state": "changes_requested"}})
}

// readRole is "" for no read marker.
func TestHomeworkUnreadInEveryView(t *testing.T) {
	tests := []struct {
		name     string
		viewer   string
		author   string
		kind     string
		readRole string
		readAt   string
		unread   bool
	}{
		{name: "message from the other role without a read marker", viewer: "parent", author: "student", kind: "reply", readRole: "", readAt: "", unread: true},
		{name: "own message", viewer: "student", author: "student", kind: "reply", readRole: "", readAt: "", unread: false},
		{name: "read after the message", viewer: "parent", author: "student", kind: "reply", readRole: "parent", readAt: "2026-10-05T10:00:00.000Z", unread: false},
		{name: "message after the read", viewer: "parent", author: "student", kind: "reply", readRole: "parent", readAt: "2026-10-05T09:00:00.000Z", unread: true},
		{name: "read by the other role only", viewer: "parent", author: "student", kind: "reply", readRole: "student", readAt: "2026-10-05T10:00:00.000Z", unread: true},
		{name: "requested changes for the student", viewer: "student", author: "parent", kind: "changes_requested", readRole: "", readAt: "", unread: true},
		{name: "completed review for the student", viewer: "student", author: "parent", kind: "completed", readRole: "", readAt: "", unread: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := homeworkSetup(t)
			exec(t, f.pool, homeworkInsertFeedback, newID(), f.problemID, tt.author, tt.kind, "2026-10-05T10:00:00.000Z")
			if tt.readRole != "" {
				exec(t, f.pool, homeworkInsertRead, tt.readRole, f.problemID, tt.readAt)
			}
			cookie := f.as(tt.viewer)
			tasks := decode[homeworkPage](t, request(t, f.s, "/tasks", http.MethodGet, nil, cookie)).Items
			bundle := decode[homeworkDetails](t, request(t, f.s, "/homework/"+f.homeworkID, http.MethodGet, nil, cookie)).Tasks
			summary := decode[homeworkPage](t, request(t, f.s, "/homework", http.MethodGet, nil, cookie)).Items
			dashboard := decode[homeworkDashboard](t, request(t, f.s, "/dashboard?today=2026-10-07", http.MethodGet, nil, cookie)).Homework
			got := map[string]any{"tasks": homeworkFirst(t, tasks)["unread"], "bundle": homeworkFirst(t, bundle)["unread"], "summary": homeworkFirst(t, summary)["unread"], "dashboard": homeworkFirst(t, dashboard)["unread"]}
			want := map[string]any{"tasks": tt.unread, "bundle": tt.unread, "summary": tt.unread, "dashboard": tt.unread}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("unread = %v, want %v", got, want)
			}
		})
	}
}
