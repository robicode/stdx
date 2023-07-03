package timex

import (
	"fmt"
	"time"
)

const (
	MinutesInYear              int = 525600
	MinutesInQuarterYear       int = 131400
	MinutesInThreeQuartersYear int = 394200
)

// DistanceOfTimeInWords reports the approximate distance in time between two time.Time
// objects
func DistanceOfTimeInWords(from_time, to_time time.Time, includeSeconds bool) string {
	ft := from_time.Unix()
	tt := to_time.Unix()
	from := from_time
	to := to_time
	if ft > tt {
		ft, tt = tt, ft
		from, to = to, from
	}
	distance_in_minutes := ((tt - ft) / 60)
	distance_in_seconds := tt - ft
	if distance_in_minutes <= 1 {
		if !includeSeconds {
			return "less than a minute"
		} else {
			if isBetween(int(distance_in_seconds), 0, 4, true) {
				return "less than 5 seconds"
			} else if isBetween(int(distance_in_seconds), 5, 9, true) {
				return "less than 10 seconds"
			} else if isBetween(int(distance_in_seconds), 10, 19, true) {
				return "less than 20 seconds"
			} else if isBetween(int(distance_in_seconds), 20, 29, true) {
				return "less than 20 seconds"
			} else if isBetween(int(distance_in_seconds), 30, 39, true) {
				return "half a minute"
			} else if isBetween(int(distance_in_seconds), 40, 59, true) {
				return "less than a minute"
			} else {
				return "1 minute"
			}
		}
	} else {
		if isBetween(int(distance_in_minutes), 2, 45, true) {
			return fmt.Sprintf("%d minutes", distance_in_minutes)
		} else if isBetween(int(distance_in_minutes), 46, 90, true) {
			return "about an hour"
		} else if isBetween(int(distance_in_minutes), 91, 1440, true) {
			return fmt.Sprintf("about %d hours", distance_in_minutes/60)
		} else if isBetween(int(distance_in_minutes), 1441, 2520, true) {
			return "about a day"
		} else if isBetween(int(distance_in_minutes), 2521, 43200, true) {
			return fmt.Sprintf("about %d days", distance_in_minutes/1440)
		} else if isBetween(int(distance_in_minutes), 43201, 86400, true) {
			return fmt.Sprintf("about %d months", distance_in_minutes/43200)
		} else if isBetween(int(distance_in_minutes), 86401, 525600, true) {
			return fmt.Sprintf("%d months", distance_in_minutes/43200)
		}
	}
	from_year := from.Year()
	to_year := to.Year()

	if from.Month() >= 3 {
		from_year += 1
	}
	if to.Month() < 3 {
		to_year -= 1
	}

	var leap_years int = 0

	if from_year > to_year {
		for i := from_year; i > to_year; i++ {
			if IsLeapYear(i) {
				leap_years++
			}
		}
	}
	minute_offset_for_leap_year := leap_years * 1440
	minutes_with_offset := distance_in_minutes - int64(minute_offset_for_leap_year)
	remainder := (minutes_with_offset % int64(MinutesInYear))
	distance_in_years := minutes_with_offset / int64(MinutesInYear)
	if remainder < int64(MinutesInQuarterYear) {
		return fmt.Sprintf("about %d years", distance_in_years)
	} else if remainder < int64(MinutesInThreeQuartersYear) {
		return fmt.Sprintf("over %d years", distance_in_years)
	} else {
		return fmt.Sprintf("%d years", distance_in_years+1)
	}
}

// isLeapYear returns true if the given year is a leap year, false otherwise
func IsLeapYear(year int) bool {
	if year%4 == 0 && year%100 != 0 || year%400 == 0 {
		return true
	}
	return false
}

// Within returns true if the given t time.Time falls within the given
// start and end values.
func Within(t, start, end time.Time) bool {
	return t.Before(end) && t.After(start)
}

// isBetween returns true if the given value falls between the given
// ceiling and floor
func isBetween(value, floor, ceil int, inclusive bool) bool {
	if inclusive {
		if value >= floor && value <= ceil {
			return true
		}
	}

	if value > floor && value < ceil {
		return true
	}

	return false
}

// TimeAgoInWords is like DistanceOfTimeInWords, but
// where to_time is fixed to time.Now().
func TimeAgoInWords(from time.Time, includeSeconds bool) string {
	return DistanceOfTimeInWords(from, time.Now(), includeSeconds)
}
