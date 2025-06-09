package common

import "time"

func StrToStrPtr(s string) *string {
	return &s
}

func StrPtrToStr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func IntToIntPtr(i int) *int {
	return &i
}

func IntPtrToInt(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

func Float64ToFloat64Ptr(f float64) *float64 {
	return &f
}

func Float64PtrToFloat64(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0
}

func Int64ToInt64Ptr(i int64) *int64 {
	return &i
}

func Int64PtrToInt64(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

func BoolToBoolPtr(b bool) *bool {
	return &b
}

func BoolPtrToBool(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func TimePtrToTime(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}

func TimeToTimePtr(t time.Time) *time.Time {
	return &t
}

func TimePtrToDay(t *time.Time) int {
	if t != nil {
		_, _, d := t.Date()
		return d
	}
	return 0
}

func TimePtrToMonth(t *time.Time) int {
	if t != nil {
		_, m, _ := t.Date()
		return int(m)
	}
	return 0
}

func TimePtrToYear(t *time.Time) int {
	if t != nil {
		y, _, _ := t.Date()
		return y
	}
	return 0
}
