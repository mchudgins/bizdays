package bizdays

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

var holidays = []string{"20170102", "20170116"}

var intervalTests = []struct {
	start                      string
	end                        string
	caldays, bizdays, observed int
}{
	{"20170101", "20170102", 1, 0, 0}, // Sunday to Monday
	{"20170101", "20170103", 2, 1, 0}, // Sunday to Tuesday
	{"20170101", "20170104", 3, 2, 1},
	{"20170101", "20170105", 4, 3, 2},
	{"20170101", "20170106", 5, 4, 3},   // Sunday to Friday
	{"20170101", "20170107", 6, 5, 4},   // Sunday to Saturday
	{"20170101", "20170108", 7, 5, 4},   // Sunday to Sunday
	{"20170101", "20170109", 8, 5, 4},   // Sunday to Monday
	{"20170101", "20170110", 9, 6, 5},   // Sunday to Tuesday
	{"20170101", "20170111", 10, 7, 6},  // Sunday to Wednesday
	{"20170101", "20170113", 12, 9, 8},  // Sunday to Friday
	{"20170101", "20170114", 13, 10, 9}, // Sunday to Saturday
	{"20170101", "20170115", 14, 10, 9}, // Sunday to Sunday
	{"20170101", "20170116", 15, 10, 9}, // Sunday to Monday
	{"20170101", "20170117", 16, 11, 9}, // Sunday to Tuesday

	{"20161231", "20170101", 1, 0, 0}, // Saturday to Sunday
	{"20161231", "20170102", 2, 0, 0}, // Saturday to Monday
	{"20161231", "20170103", 3, 1, 0}, // Saturday to Tuesday
	{"20161231", "20170104", 4, 2, 1},
	{"20161231", "20170106", 6, 4, 3},   // Saturday to Friday
	{"20161231", "20170107", 7, 5, 4},   // Saturday to Saturday
	{"20161231", "20170108", 8, 5, 4},   // Saturday to Sunday
	{"20161231", "20170109", 9, 5, 4},   // Saturday to Monday
	{"20161231", "20170110", 10, 6, 5},  // Saturday to Tuesday
	{"20161231", "20170111", 11, 7, 6},  // Saturday to Wednesday
	{"20161231", "20170113", 13, 9, 8},  // Saturday to Friday
	{"20161231", "20170114", 14, 10, 9}, // Saturday to Saturday
	{"20161231", "20170115", 15, 10, 9}, // Saturday to Sunday
	{"20161231", "20170116", 16, 10, 9}, // Saturday to Monday
	{"20161231", "20170117", 17, 11, 9}, // Saturday to Tuesday

	{"20161230", "20161231", 1, 0, 0}, // Friday to Saturday
	{"20161230", "20170101", 2, 0, 0}, // Friday to Sunday
	{"20161230", "20170102", 3, 0, 0}, // Friday to Monday
	{"20161230", "20170103", 4, 1, 0}, // Friday to Tuesday
	{"20161230", "20170104", 5, 2, 1},
	{"20161230", "20170106", 7, 4, 3},    // Friday to Friday
	{"20161230", "20170107", 8, 5, 4},    // Friday to Saturday
	{"20161230", "20170108", 9, 5, 4},    // Friday to Sunday
	{"20161230", "20170109", 10, 5, 4},   // Friday to Monday
	{"20161230", "20170110", 11, 6, 5},   // Friday to Tuesday
	{"20161230", "20170111", 12, 7, 6},   // Friday to Wednesday
	{"20161230", "20170113", 14, 9, 8},   // Friday to Friday
	{"20161230", "20170114", 15, 10, 9},  // Friday to Saturday
	{"20161230", "20170115", 16, 10, 9},  // Friday to Sunday
	{"20161230", "20170116", 17, 10, 9},  // Friday to Monday
	{"20161230", "20170117", 18, 11, 9},  // Friday to Tuesday
	{"20161230", "20170118", 19, 12, 10}, // Friday to Wednesday

	{"20170103", "20170104", 1, 0, 0},   // Tuesday to Wednesday
	{"20170103", "20170105", 2, 1, 1},   // Tuesday to Thursday
	{"20170103", "20170106", 3, 2, 2},   // Tuesday to Friday
	{"20170103", "20170107", 4, 3, 3},   // Tuesday to Saturday
	{"20170103", "20170108", 5, 3, 3},   // Tuesday to Sunday
	{"20170103", "20170109", 6, 3, 3},   // Tuesday to Monday
	{"20170103", "20170110", 7, 4, 4},   // Tuesday to Tuesday
	{"20170103", "20170111", 8, 5, 5},   // Tuesday to Wednesday
	{"20170103", "20170112", 9, 6, 6},   // Tuesday to Thursday
	{"20170103", "20170113", 10, 7, 7},  // Tuesday to Friday
	{"20170103", "20170114", 11, 8, 8},  // Tuesday to Saturday
	{"20170103", "20170115", 12, 8, 8},  // Tuesday to Sunday
	{"20170103", "20170116", 13, 8, 8},  // Tuesday to Monday
	{"20170103", "20170117", 14, 9, 8},  // Tuesday to Tuesday
	{"20170103", "20170118", 15, 10, 9}, // Tuesday to Wednesday
}

func createDebugLogger(w io.Writer) *log.Logger {
	tmpLogger := log.New() //(&log.Logger{Out: os.Stderr, Level: log.DebugLevel}).
	tmpLogger.Level = log.DebugLevel
	tmpLogger.Out = w

	return tmpLogger
}

func parseHolidays(t []string) []time.Time {
	i := len(t)
	results := make([]time.Time, i, i)

	for i, h := range t {
		ht, err := time.Parse("20060102", h)
		if err != nil {
			panic(fmt.Sprintf("Unable to parse start time (%s) as time -- %s\n", h, err))
		}
		results[i] = ht
	}

	return results
}

func TestCalIntervals(t *testing.T) {
	for _, tt := range intervalTests {

		start, err := time.Parse("20060102", tt.start)
		if err != nil {
			t.Errorf("Unable to parse start time (%s) as time -- %s\n", tt.start, err)
		}
		end, err := time.Parse("20060102", tt.end)
		if err != nil {
			t.Errorf("Unable to parse start time (%s) as time -- %s\n", tt.end, err)
		}
		result := CalDaysDiff(start, end)

		if result != tt.caldays {
			t.Errorf("expected: %d, got %d for start %s and end %s", tt.caldays, result, tt.start, tt.end)
			t.Fail()
		}
	}
}

func TestBizIntervals(t *testing.T) {
	var pass, fail int

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	testLogger := createDebugLogger(buf)

	holidayTable := []time.Time{}

	for _, tt := range intervalTests {

		start, err := time.Parse("20060102", tt.start)
		if err != nil {
			t.Errorf("Unable to parse start time (%s) as time -- %s\n", tt.start, err)
		}
		end, err := time.Parse("20060102", tt.end)
		if err != nil {
			t.Errorf("Unable to parse start time (%s) as time -- %s\n", tt.end, err)
		}

		Logger = testLogger.WithFields(log.Fields{"test": "TestBizIntervals",
			"start": tt.start,
			"end":   tt.end})

		result := BizDaysDiff(start, end, holidayTable)
		if result != tt.bizdays {
			t.Errorf("expected: %d, got %d for start %s (%s) and end %s (%s)",
				tt.bizdays,
				result,
				tt.start,
				start.Weekday(),
				tt.end,
				end.Weekday())
			if buf.Len() > 0 {
				t.Errorf("Messages logged during test execution:\n%s[%d/%d]\n",
					buf, buf.Len(), buf.Cap())
			}
			t.Fail()
			fail++
		} else {
			pass++
		}

		buf.Reset()
	}
	fmt.Printf("Pass:  %d, Fail:  %d\n", pass, fail)
}

func TestBizIntervalsWithHolidays(t *testing.T) {
	var pass, fail int

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	testLogger := createDebugLogger(buf)

	holidayTable := parseHolidays(holidays)

	for _, tt := range intervalTests {

		start, err := time.Parse("20060102", tt.start)
		if err != nil {
			t.Errorf("Unable to parse start time (%s) as time -- %s\n", tt.start, err)
		}
		end, err := time.Parse("20060102", tt.end)
		if err != nil {
			t.Errorf("Unable to parse start time (%s) as time -- %s\n", tt.end, err)
		}

		Logger = testLogger.WithFields(log.Fields{"test": "TestBizIntervalsWithHolidays",
			"start": tt.start,
			"end":   tt.end})

		result := BizDaysDiff(start, end, holidayTable)
		if result != tt.observed {
			t.Errorf("expected: %d, got %d for start %s (%s) and end %s (%s)",
				tt.observed,
				result,
				tt.start,
				start.Weekday(),
				tt.end,
				end.Weekday())
			if buf.Len() > 0 {
				t.Errorf("Messages logged during test execution:\n%s[%d/%d]\n",
					buf, buf.Len(), buf.Cap())
			}
			t.Fail()
			fail++
		} else {
			pass++
		}

		buf.Reset()
	}
	fmt.Printf("Pass:  %d, Fail:  %d\n", pass, fail)
}
