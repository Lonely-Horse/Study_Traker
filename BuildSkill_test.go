package main

import "testing"

func TestBuildSkill(t *testing.T) {
	tests := []struct {
		name         string
		logs         []StudyLog
		wantCount    int
		wantDuration int
		wantSubjects []string
	}{
		{
			name: "Test",
			logs: []StudyLog{
				{
					Date:           "2026-02-04",
					Subject:        "高数",
					Duration:       120,
					ChallengeLevel: 3,
				},
				{
					Date:           "2026-02-05",
					Subject:        "高数",
					Duration:       120,
					ChallengeLevel: 3,
				},
			},
			wantCount:    2,
			wantDuration: 240,
			wantSubjects: []string{"高数"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildSkill(tt.logs)
			if result.Count != tt.wantCount {
				t.Errorf("Count = %d ,want %d", result.Count, tt.wantCount)
			}
			if result.TotalDurationMin != tt.wantDuration {
				t.Errorf("TotalDurationMin = %d,want %d", result.TotalDurationMin, tt.wantDuration)
			}
			if len(result.Subjects) != len(tt.wantSubjects) {
				t.Fatalf("Subject length = %d ,want %d", len(result.Subjects), len(tt.wantSubjects))
			}
			for i := range tt.wantSubjects {
				if result.Subjects[i] != tt.wantSubjects[i] {
					t.Errorf("Subjects[%d] = %s ,want %s", i, result.Subjects[i], tt.wantSubjects[i])
				}
			}
		})
	}

}
