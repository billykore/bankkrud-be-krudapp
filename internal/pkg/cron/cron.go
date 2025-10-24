package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetWeekday parses the weekday field from a cron expression
// and returns the name of the day ("Sunday", "Monday", ...).
func GetWeekday(expr string) string {
	if expr == "" {
		return "undefined"
	}
	s := strings.Split(expr, " ") // ["0", "5", "*", "*", "[day]"]
	if len(s) < 5 {
		return "undefined"
	}
	day := s[4]
	intDay, err := strconv.Atoi(day)
	if err != nil {
		return "undefined"
	}
	if intDay < 0 || intDay > 6 {
		return "undefined"
	}
	return time.Weekday(intDay).String()
}

// GetDate parses the date field from a cron expression
// and returns it as an integer (1, 2, 3, ...).
func GetDate(expr string) int {
	if expr == "" {
		return 0
	}
	s := strings.Split(expr, " ") // ["0", "5", "[date]", "*", "*"]
	if len(s) < 5 {
		return 0
	}
	d := s[2]
	intDate, err := strconv.Atoi(d)
	if err != nil {
		return 0
	}
	return intDate
}

// DailyExpression returns a cron expression that triggers daily at 5:00 AM.
func DailyExpression() string {
	return "0 5 * * *"
}

// weekdaysValues maps weekday names to their corresponding cron values.
var weekdaysValues = map[string]string{
	"Sunday":    "0",
	"Monday":    "1",
	"Tuesday":   "2",
	"Wednesday": "3",
	"Thursday":  "4",
	"Friday":    "5",
	"Saturday":  "6",
}

// WeeklyExpression returns a cron expression that triggers weekly
// on the specified weekday at 5:00 AM.
func WeeklyExpression(weekday string) string {
	return fmt.Sprintf("0 5 * * %s", weekdaysValues[weekday])
}

// MonthlyExpression returns a cron expression that triggers monthly
// on the specified date at 5:00 AM.
func MonthlyExpression(date int) string {
	return fmt.Sprintf("0 5 %d * *", date)
}
