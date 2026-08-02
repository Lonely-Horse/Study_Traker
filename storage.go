package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

func ReadLogs(filename string) ([]StudyLog, error) {
	var logs []StudyLog

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []StudyLog{}, nil
		}
		return nil, fmt.Errorf("failed to open %s:%w", filename, err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&logs)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return []StudyLog{}, nil
		}
		return nil, fmt.Errorf("failed decode the %s:%w", filename, err)
	}
	return logs, nil
}

func ReadLogsSafe(filename string) (logs []StudyLog, err error) {
	mu.RLock()
	defer mu.RUnlock()
	logs, err = ReadLogs(filename)
	return logs, err
}

func WriteLogs(filename string, logs []StudyLog) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed open/create %s:%w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err = encoder.Encode(logs)
	if err != nil {
		return fmt.Errorf("failed encode the %s: %w", filename, err)
	}

	return nil
}

func BuildSkill(logs []StudyLog) SkillResponse {
	tdm := 0
	c := len(logs)

	subjectSet := make(map[string]bool)
	for _, log := range logs {
		tdm = tdm + log.Duration
		if log.Subject != "" {
			subjectSet[log.Subject] = true
		}
	}

	subjects := make([]string, 0, len(subjectSet))
	for subject := range subjectSet {
		subjects = append(subjects, subject)
	}

	sort.Strings(subjects)

	skillresponse := SkillResponse{
		Count:            c,
		TotalDurationMin: tdm,
		Subjects:         subjects,
		Logs:             logs,
	}

	return skillresponse
}

func AppendLog(filename string, newlog StudyLog) error {
	mu.Lock()
	defer mu.Unlock()
	logs, err := ReadLogs(filename)
	if err != nil {
		return err
	}
	logs = append(logs, newlog)
	err = WriteLogs(filename, logs)
	if err != nil {
		return err
	}

	return nil
}

func FilterLogsBySubject(logs []StudyLog, subject string) []StudyLog {

	if logs == nil || subject == "" {
		return logs
	}

	tmplog := make([]StudyLog, 0, len(logs))

	for _, log := range logs {
		if log.Subject == subject {
			tmplog = append(tmplog, log)
		}
	}
	return tmplog
}
