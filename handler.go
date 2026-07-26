package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
)

func HandlerDashboard(w http.ResponseWriter, r *http.Request, filename string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("The method error"))
		return
	}

	logs, err := ReadLogsSafe(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed read the logs"))
		return
	}

	tmpl, err := template.ParseFiles("templates/dashboard.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed find the template file"))
		return
	}

	err = tmpl.Execute(w, logs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed execute the logs"))
		return
	}
}

func HandlerGetLogs(w http.ResponseWriter, r *http.Request, filename string) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"method not allowed"}`))
		return
	}

	subjectQuery := r.URL.Query().Get("subject")

	logs, err := ReadLogsSafe(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed read the logs"))
		return
	}

	logs_Query := FilterLogsBySubject(logs, subjectQuery)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(logs_Query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed encode the logs"}`))
		return
	}
}

func HandlerPostlog(w http.ResponseWriter, r *http.Request, filename string) {
	var newlog StudyLog

	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"The method isn't Post"}`))
		return
	}

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 5<<20))
	err := decoder.Decode(&newlog)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"The newlog can't decoder"}`))
		return
	}

	err = ValidateStudyLog(newlog)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = AppendLog(filename, newlog)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"The newlog can't append"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"Status":"OK"}`))
}

func HandlerLogs(w http.ResponseWriter, r *http.Request, filename string) {
	switch r.Method {
	case http.MethodGet:
		HandlerGetLogs(w, r, filename)
	case http.MethodPost:
		HandlerPostlog(w, r, filename)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Don't have the method"}`))
	}
}

func HandlerSkill(w http.ResponseWriter, r *http.Request, filename string) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"The method isn't MethodGet"}`))
		return
	}

	logs, err := ReadLogsSafe(filename)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Don't read the logs"}`))
		return
	}
	skillresponse := BuildSkill(logs)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(skillresponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed encode the logs"}`))
		return
	}
}

func AuthMiddleware(next http.HandlerFunc, token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if token == "" {
			next(w, r)
			return
		}
		authhead := r.Header.Get("Authorization")
		if !strings.HasPrefix(authhead, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
		authtoken := strings.TrimPrefix(authhead, "Bearer ")
		if authtoken != token {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}
		next(w, r)
	}
}
