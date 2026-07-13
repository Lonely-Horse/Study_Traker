package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func ReadLogs(filename string) ([]StudyLog, error) {
	var logs []StudyLog

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s:%w", filename, err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&logs)
	if err != nil {
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
	tmpfile := filename + ".tmp"
	tmp, err := os.Create(tmpfile)
	if err != nil {
		return fmt.Errorf("failed create the %s:%w", filename, err)
	}

	encoder := json.NewEncoder(tmp)
	encoder.SetIndent("", " ")
	err = encoder.Encode(logs)
	if err != nil {
		tmp.Close()
		err1 := os.Remove(tmpfile)
		if err1 != nil {
			return fmt.Errorf("failed remove the %s: %w", tmpfile, err1)
		}
		return fmt.Errorf("failed encode the %s: %w", filename, err)
	}

	err = tmp.Close()
	if err != nil {
		os.Remove(tmpfile)
		return fmt.Errorf("failed close the %s: %w", tmpfile, err)
	}

	err = os.Rename(tmpfile, filename)
	if err != nil {
		os.Remove(tmpfile)
		return fmt.Errorf("failed rename the %s to %s: %w", tmpfile, filename, err)
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
