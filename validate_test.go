package main

import "testing"

func TestValidateStudyLog(t *testing.T) {
	tests := []struct {
		name    string
		log     StudyLog
		wantErr bool
	}{
		{
			name: "合法日志",
			log: StudyLog{
				Date:           "2026-11-14",
				Subject:        "高数",
				Theme:          "微分方程",
				Duration:       20,
				ChallengeLevel: 3,
			},
			wantErr: false,
		},
		{
			name: "Date为空",
			log: StudyLog{
				Date:           "",
				Subject:        "大学英语",
				Theme:          "词汇",
				Duration:       100,
				ChallengeLevel: 4,
			},
			wantErr: true,
		},
		{
			name: "Subject为空",
			log: StudyLog{
				Date:           "2021-2-2",
				Subject:        "",
				Theme:          "函数解析",
				Duration:       120,
				ChallengeLevel: 5,
			},
			wantErr: true,
		},
		{
			name: "Duration为0",
			log: StudyLog{
				Date:           "2021-2-2",
				Subject:        "高数",
				Theme:          "函数解析",
				Duration:       0,
				ChallengeLevel: 5,
			},
			wantErr: true,
		},
		{
			name: "Duration为负",
			log: StudyLog{
				Date:           "2021-2-2",
				Subject:        "高数",
				Theme:          "函数解析",
				Duration:       -1,
				ChallengeLevel: 5,
			},
			wantErr: true,
		},
		{
			name: "ChallengeLevel为0",
			log: StudyLog{
				Date:           "2021-2-2",
				Subject:        "高数",
				Theme:          "函数解析",
				Duration:       120,
				ChallengeLevel: 0,
			},
			wantErr: true,
		},
		{
			name: "ChallengeLevel为6",
			log: StudyLog{
				Date:           "2021-2-2",
				Subject:        "高数",
				Theme:          "函数解析",
				Duration:       120,
				ChallengeLevel: 6,
			},
			wantErr: true,
		},
		{
			name: "ChallengeLevel为1",
			log: StudyLog{
				Date:           "2021-2-2",
				Subject:        "高数",
				Theme:          "函数解析",
				Duration:       120,
				ChallengeLevel: 1,
			},
			wantErr: false,
		},
		{
			name: "ChallengeLevel为5",
			log: StudyLog{
				Date:           "2021-2-2",
				Subject:        "高数",
				Theme:          "函数解析",
				Duration:       120,
				ChallengeLevel: 5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStudyLog(tt.log)
			if tt.wantErr && err == nil {
				t.Errorf("expected error,got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("expcted no error,got %v", err)
			}
		})
	}
}
