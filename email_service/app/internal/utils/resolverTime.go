package utils

import (
	"errors"
	"math"
	"strconv"
)

const (
	second          = 1
	secondsInMinute = 60 * second
	secondsInHour   = 60 * secondsInMinute
	secondsInDay    = 24 * secondsInHour
	secondsInWeek   = 7 * secondsInDay
)

func ResolveTime(timeInSeconds float64) (string, error) {
	if timeInSeconds < second {
		return "", errors.New("Time less than second")
	}

	if timeInSeconds == second {
		return "1 second", nil
	}

	if timeInSeconds <= 1.5*secondsInMinute {
		return strconv.Itoa(int(timeInSeconds)) + " seconds", nil
	}

	if timeInSeconds <= 1.5*secondsInHour {
		minutes := math.Round(timeInSeconds / secondsInMinute)
		return strconv.Itoa(int(minutes)) + " minutes", nil
	}

	if timeInSeconds < 24*secondsInHour {
		hours := math.Round(timeInSeconds / secondsInHour)
		return strconv.Itoa(int(hours)) + " hours", nil
	}

	if timeInSeconds == secondsInDay {
		return "1 day", nil
	}

	if timeInSeconds < secondsInWeek {
		days := math.Round(timeInSeconds / secondsInDay)
		return strconv.Itoa(int(days)) + " days", nil
	}

	if timeInSeconds > 10*secondsInWeek {
		return "", errors.New("More than 10 weeks")
	}

	weeks := math.Round(timeInSeconds / secondsInWeek)
	return strconv.Itoa(int(weeks)) + " weeks", nil
}
