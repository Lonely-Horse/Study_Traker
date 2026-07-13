package main

type StudyLog struct {
	Date           string `json:"date"`
	Subject        string `json:"subject"`
	Theme          string `json:"theme"`
	Duration       int    `json:"duration"`
	Challenge      string `json:"challenge,omitempty"`
	ChallengeLevel int    `json:"challenge_level"`
}

type SkillResponse struct {
	Count            int        `json:"count"`
	TotalDurationMin int        `json:"total_duration_min"`
	Subjects         []string   `json:"subjects"`
	Logs             []StudyLog `json:"logs"`
}
