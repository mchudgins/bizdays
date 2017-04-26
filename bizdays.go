package bizdays

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

var Logger log.FieldLogger

func init() {
	Logger = log.New().WithField("package", "bizdays")
}

func JulianDay(t time.Time) int {
	return t.Day() - 32075 +
		1461*(t.Year()+4800+(int(t.Month())-14)/12)/4 +
		367*(int(t.Month())-2-(int(t.Month())-14)/12*12)/12 -
		3*((t.Year()+4900+(int(t.Month())-14)/12)/100)/4
}

func CalDaysDiff(start time.Time, end time.Time) int {
	return JulianDay(end) - JulianDay(start)
}

// BizDaysDiff returns the number of intervening business days between the start and end dates.
// e.g., Between adjacent mondays and tuesdays there are 0 business days.
func BizDaysDiff(start time.Time, end time.Time, holidays []time.Time) int {

	caldays := CalDaysDiff(start, end)
	// see http://stackoverflow.com/questions/1617049/calculate-the-number-of-business-days-between-two-dates
	bizdays := (caldays*5-(int(start.Weekday()-end.Weekday())*2))/7 - 1

	// empirical fixes to the above formula

	if end.Weekday() == time.Sunday {
		Logger.Debug("end date is a Sunday")
		bizdays++
	}
	if start.Weekday() == time.Saturday {
		Logger.Debug("start date is a Saturday")
		bizdays++
	}

	if len(holidays) == 0 {
		return bizdays
	}

	for _, hday := range holidays {
		if hday.After(start) && hday.Before(end) {
			bizdays--
		}
	}

	return bizdays
}

func DateFromBizDays(now time.Time, bizdays int) time.Time {
	return now
}
