package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
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

	logs, err := ReadLogsSafe(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed read the logs"))
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(logs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Failed encode the logs"}`))
		return
	}
}

func HandlerPostlog(w http.ResponseWriter, r *http.Request, filename string) {
	var newlog StudyLog
	var extra any

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

	err = decoder.Decode(&extra)
	if err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Have some json not allowed"}`))
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
