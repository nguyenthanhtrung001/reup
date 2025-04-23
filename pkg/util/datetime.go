package util

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
)

func StrToDateTime(str string, tz *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation(DateTimeFormat, str, tz)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func StrToDateTimeInclMinutes(str string, tz *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation(DateTimeFormat, str, tz)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func DateTimeToStr(dt time.Time, ft string) string {
	if ft == "" {
		return dt.Format(DateTimeFormat)
	} else {
		return dt.Format(ft)
	}
}

func Now() time.Time {
	return time.Now().In(GetDefaultTimezone())
}

func GetDefaultTimezone() *time.Location {
	localTimeZone, _ := time.LoadLocation("Local")
	return localTimeZone
}

func GetAsiaHoChiMinhTimezone() string {
	asiaHoChiMinhTimeZone, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return asiaHoChiMinhTimeZone.String()
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, GetDefaultTimezone())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, GetDefaultTimezone())
}

func StartOfMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, GetDefaultTimezone())
}

func EndOfMinute(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 59, 0, GetDefaultTimezone())
}

func SetHour(t time.Time, hour int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), hour, t.Minute(), t.Second(), 0, GetDefaultTimezone())
}

func SetMinute(t time.Time, minute int) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), minute, t.Second(), 0, GetDefaultTimezone())
}

func DateTimeToInt(dt time.Time) int {
	return int(dt.Unix())
}

func FormatTime(t time.Time, format string) string {
	goFormat := convertPHPToGoTimeFormat(format)
	return t.Format(goFormat)
}

func convertPHPToGoTimeFormat(format string) string {
	replacements := map[string]string{
		"Y": "2006",
		"m": "01",
		"d": "02",
		"H": "15",
		"i": "04",
		"s": "05",
	}

	for php, goTime := range replacements {
		format = strings.Replace(format, php, goTime, -1)
	}

	return format
}

func BuildDateTimeStrFromDateStrAndHourMinute(date string, hour int, minute int) (string, error) {
	dateTime, err := StrToDateTime(fmt.Sprintf("%s 00:00:00", date), GetDefaultTimezone())
	if err != nil {
		return "", err
	}
	if hour > 23 {
		dateTime = dateTime.AddDate(0, 0, 1)
		hour = hour - 24
	}

	if hour < 0 {
		dateTime = dateTime.AddDate(0, 0, -1)
		hour = hour + 24
	}

	dateTime = SetHour(dateTime, hour)
	dateTime = SetMinute(dateTime, minute)

	return DateTimeToStr(dateTime, ""), nil
}

func MillisecondsToTime(ms int64) time.Time {
	seconds := ms / 1000
	nanoseconds := (ms % 1000) * 1000000
	return time.Unix(seconds, nanoseconds)
}

func MicrosecondsToTime(ms int64) time.Time {
	seconds := ms / 1000000
	nanoseconds := (ms % 1000000) * 1000
	return time.Unix(seconds, nanoseconds)
}

func GetWeekOfMonth(dt time.Time) int {
	// Ngày đầu tiên của tháng
	firstDay := time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())

	// Số ngày trong tháng
	dayOfMonth := dt.Day()

	// Điều chỉnh dựa trên ngày trong tuần của ngày đầu tiên
	adjustedDay := dayOfMonth + int(firstDay.Weekday())

	// Tính số tuần (sử dụng ceil)
	return int(math.Ceil(float64(adjustedDay) / 7.0))
}

func GetWeekOfYear(dt time.Time) int {
	_, week := dt.ISOWeek()
	return week
}

func strContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
