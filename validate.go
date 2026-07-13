package main

import (
	"errors"
)

func ValidateStudyLog(log StudyLog) error {
	if log.Date == "" {
		return errors.New("The Date is nil")
	}
	if log.Subject == "" {
		return errors.New("The Subject is nil")
	}
	if log.Duration <= 0 {
		return errors.New("The Duration isn't right")
	}
	if log.ChallengeLevel < 1 || log.ChallengeLevel > 5 {
		return errors.New("The ChallengeLavel isn't right")
	}

	return nil
}
