package bizdays

import (
	"bytes"
	"testing"
	"time"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

var targetTests = []struct {
	date    string // initial date
	bizdays int    // +/- number of intervening bizdays
	target  string // expected result
}{
	{"20170101", 1, "20170103"},
}

func TestTargetDate(t *testing.T) {
	var pass, fail int

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	testLogger := createDebugLogger(buf)

	for _, tt := range targetTests {
		initial, err := time.Parse("20060102", tt.date)
		if err != nil {
			t.Errorf("Unable to parse initial date (%s) as time -- %s\n", tt.date, err)
		}

		Logger = testLogger.WithFields(log.Fields{"test": "TestBizIntervals",
			"initialDate": tt.date,
			"interval":    tt.bizdays})

		date := DateFromBizDays(initial, tt.bizdays)
		result := date.Format("20060102")

		if result != tt.target {
			t.Errorf("expected: %s, got %s for initial date %s (%s) and interval %d",
				tt.target,
				result,
				tt.date,
				initial.Weekday(),
				tt.bizdays)
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
