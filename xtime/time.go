package xtime

import (
	"strconv"
	"time"
)

// DateTime 获取当前时间
func DateTime() string {
	return time.Now().Format(time.DateTime)
}

// Timestamp 获取当前时间的时间戳
func Timestamp() int64 {
	return time.Now().Unix()
}

var _units = []string{"ns", "us", "ms", "s"}

// ManualDuration 优化 time.Duration 显示
func ManualDuration(d time.Duration) string {
	f := float64(d)
	i := 0
	for f > 1000 && i < len(_units)-1 {
		f = f / 1000
		i++
	}
	return strconv.FormatFloat(f, 'f', 2, 64) + _units[i]
}

var _daysOfMonth = []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// DaysOfMonth 获取指定月份天数
func DaysOfMonth(year int, month time.Month) int {
	if month < 1 || month > 12 {
		return 0
	}
	if IsLeap(year) && month == time.February {
		return 29
	}
	return _daysOfMonth[month-1]
}

// IsLeap 判断是否是闰年
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DayBegin 获取指定时间当天的开始时间
func DayBegin(t time.Time, loc *time.Location) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

// DayEnd 获取指定时间当天的结束时间
func DayEnd(t time.Time, loc *time.Location) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, loc)
}

// DayRange 获取指定时间当天的开始时间和结束时间
func DayRange(t time.Time, loc *time.Location) (time.Time, time.Time) {
	return DayBegin(t, loc), DayEnd(t, loc)
}

// WeekBegin 获取指定时间当周的开始时间
func WeekBegin(t time.Time, loc *time.Location) time.Time {
	n := int(t.Weekday())
	if n == 0 {
		n = 7
	}
	return DayBegin(t.AddDate(0, 0, 1-n), loc)
}

// WeekEnd 获取指定时间当周的结束时间
func WeekEnd(t time.Time, loc *time.Location) time.Time {
	n := int(t.Weekday())
	if n == 0 {
		n = 7
	}
	return DayEnd(t.AddDate(0, 0, 7-n), loc)
}

// WeekRange 获取指定时间当周的开始时间和结束时间
func WeekRange(t time.Time, loc *time.Location) (time.Time, time.Time) {
	return WeekBegin(t, loc), WeekEnd(t, loc)
}

// MonthBegin 获取指定时间当月的开始时间
func MonthBegin(t time.Time, loc *time.Location) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, loc)
}

// MonthEnd 获取指定时间当月的结束时间
func MonthEnd(t time.Time, loc *time.Location) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, DaysOfMonth(year, month), 23, 59, 59, 0, loc)
}

// MonthRange 获取指定时间当月的开始时间和结束时间
func MonthRange(t time.Time, loc *time.Location) (time.Time, time.Time) {
	return MonthBegin(t, loc), MonthEnd(t, loc)
}

// YearBegin 获取指定时间当年的开始时间
func YearBegin(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, loc)
}

// YearEnd 获取指定时间当年的结束时间
func YearEnd(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), 12, 31, 23, 59, 59, 0, loc)
}

// YearRange 获取指定时间当年的开始时间和结束时间
func YearRange(t time.Time, loc *time.Location) (time.Time, time.Time) {
	return YearBegin(t, loc), YearEnd(t, loc)
}
