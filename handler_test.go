package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandlerGetlogs(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		filename string
		wantCode int
	}{
		{
			name:     "Use the Get method to curl the server",
			method:   http.MethodGet,
			filename: "testdata/test_logs.json",
			wantCode: 200,
		},
		{
			name:     "Use the Post method to curl the server",
			method:   http.MethodPost,
			filename: "testdata/test_logs.json",
			wantCode: 405,
		},
		{
			name:     "Use another file to try",
			method:   http.MethodGet,
			filename: "testdata/test_log.json",
			wantCode: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/logs", nil)
			w := httptest.NewRecorder()
			HandlerGetLogs(w, req, tt.filename)
			if w.Code != tt.wantCode {
				t.Errorf("The Code is %d,want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestHandlerPostlog(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		body     string
		wantCode int
	}{
		{
			name:     "Use the POST method to try",
			method:   http.MethodPost,
			body:     `{"date":"2026-06-26","subject":"高数","theme":"极限","duration":90,"challenge_level":3}`,
			wantCode: 201,
		},
		{
			name:     "Use the GET method to try",
			method:   http.MethodGet,
			body:     "",
			wantCode: 405,
		},
		{
			name:     "Use the bad json body",
			method:   http.MethodPost,
			body:     `"date":"2026-06-26","subject":"高数","theme":"极限","duration":90,"challenge_level":3`,
			wantCode: 400,
		},
		{
			name:     "The json need date",
			method:   http.MethodPost,
			body:     `{"date":"","subject":"高数","theme":"极限","duration":90,"challenge_level":3}`,
			wantCode: 400,
		},
		{
			name:     "The json need subject",
			method:   http.MethodPost,
			body:     `{"date":"2026-06-26","subject":"","theme":"极限","duration":90,"challenge_level":3}`,
			wantCode: 400,
		},
		{
			name:     "The Duration isn't eql zero",
			method:   http.MethodPost,
			body:     `{"date":"2026-06-26","subject":"高数","theme":"极限","duration":0,"challenge_level":3}`,
			wantCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, _ := os.ReadFile("testdata/test_logs.json")
			tmpfile := "testdata/tmp_test.json"
			os.WriteFile(tmpfile, src, 0644)
			defer os.Remove(tmpfile)
			defer os.Remove(tmpfile + ".tmp")

			var body io.Reader
			if tt.body != "" {
				body = strings.NewReader(tt.body)
			}
			req := httptest.NewRequest(tt.method, "/api/logs", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			HandlerPostlog(w, req, tmpfile)
			if w.Code != tt.wantCode {
				t.Errorf("test %s failed: got status %d, want %d, body: %s", tt.name, w.Code, tt.wantCode, w.Body.String())
			}
		})
	}
}

func TestFilterLogsBySubject(t *testing.T) {
	test_logs := []StudyLog{
		{Subject: "高数", Duration: 120},
		{Subject: "英语", Duration: 100},
		{Subject: "高数", Duration: 200},
	}
	tests := []struct {
		name         string
		subjectQuery string
		wantlen      int
	}{
		{
			name:         "Search the math",
			subjectQuery: "高数",
			wantlen:      2,
		},
		{
			name:         "Search the English",
			subjectQuery: "英语",
			wantlen:      1,
		},
		{
			name:         "Search the P.E",
			subjectQuery: "体育",
			wantlen:      0,
		},
		{
			name:         "Search empty string subject",
			subjectQuery: "",
			wantlen:      3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterLogsBySubject(test_logs, tt.subjectQuery)
			if len(result) != tt.wantlen {
				t.Errorf("got length %d, want %d", tt.wantlen, len(result))
			}
		})
	}
}
