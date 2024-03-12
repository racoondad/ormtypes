package ormtypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Time is time data type.
type Time time.Duration

// NewTime is a constructor for Time and returns new Time.
// func NewTime(hour, min, sec, nsec int) Time {
// 	return newTime(hour, min, sec, nsec)
// }

func NewTime(hour, min, sec, nsec int) Time {
	return Time(
		time.Duration(hour)*time.Hour +
			time.Duration(min)*time.Minute +
			time.Duration(sec)*time.Second +
			time.Duration(nsec)*time.Nanosecond,
	)
}

func NowTime() Time {
	t := time.Now()
	return NewTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// GormDataType returns gorm common data type. This type is used for the field's column type.
func (Time) GormDataType() string {
	return "time"
}

// GormDBDataType returns gorm DB data type based on the current using database.
func (Time) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql":
		return "TIME"
	case "postgres":
		return "TIME"
	case "sqlserver":
		return "TIME"
	case "sqlite":
		return "TEXT"
	default:
		return ""
	}
}

// Scan implements sql.Scanner interface and scans value into Time,
func (t *Time) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		t.setFromString(string(v))
	case string:
		t.setFromString(v)
	case time.Time:
		t.setFromTime(v)
	default:
		return fmt.Errorf("failed to scan value: %v", v)
	}

	return nil
}

func (t *Time) setFromString(str string) {
	var h, m, s, n int
	fmt.Sscanf(str, "%02d:%02d:%02d.%09d", &h, &m, &s, &n)
	*t = NewTime(h, m, s, n)
}

func (t *Time) setFromTime(src time.Time) {
	*t = NewTime(src.Hour(), src.Minute(), src.Second(), src.Nanosecond())
}

// Value implements driver.Valuer interface and returns string format of Time.
func (t Time) Value() (driver.Value, error) {
	return t.String(), nil
}

// String implements fmt.Stringer interface.
func (t Time) String() string {
	if nsec := t.nanoseconds(); nsec > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%09d", t.hours(), t.minutes(), t.seconds(), nsec)
	} else {
		// omit nanoseconds unless any value is specified
		return fmt.Sprintf("%02d:%02d:%02d", t.hours(), t.minutes(), t.seconds())
	}
}

func (t Time) hours() int {
	return int(time.Duration(t).Truncate(time.Hour).Hours())
}

func (t Time) minutes() int {
	return int((time.Duration(t) % time.Hour).Truncate(time.Minute).Minutes())
}

func (t Time) seconds() int {
	return int((time.Duration(t) % time.Minute).Truncate(time.Second).Seconds())
}

func (t Time) nanoseconds() int {
	return int((time.Duration(t) % time.Second).Nanoseconds())
}

// MarshalJSON implements json.Marshaler to convert Time to json serialization.
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON implements json.Unmarshaler to deserialize json data.
func (t *Time) UnmarshalJSON(data []byte) error {
	// ignore null
	if string(data) == "null" {
		return nil
	}
	t.setFromString(strings.Trim(string(data), `"`))
	return nil
}

func (t Time) Ago(minute int) Time {
	h := minute / 60
	m := minute % 60
	return NewTime(t.hours()+h, t.minutes()+m, t.seconds(), t.nanoseconds())
}

func (t Time) After(in Time) bool {
	if t.hours() > in.hours() {
		return true
	}
	if t.minutes() > in.minutes() {
		return true
	}
	if t.seconds() > in.seconds() {
		return true
	}
	if t.nanoseconds() > in.nanoseconds() {
		return true
	}
	return false
}

func (t Time) Before(in Time) bool {
	if t.hours() < in.hours() {
		return true
	}
	if t.minutes() < in.minutes() {
		return true
	}
	if t.seconds() < in.seconds() {
		return true
	}
	if t.nanoseconds() < in.nanoseconds() {
		return true
	}
	return false
}

func (t Time) Equal(in Time) bool {
	if t.hours() == in.hours() && t.minutes() == in.minutes() && t.seconds() == in.seconds() && t.nanoseconds() == in.nanoseconds() {
		return true
	}
	return false
}

func (t Time) IsZero() bool {
	return t.seconds() == 0 && t.nanoseconds() == 0
}

func (t Time) SubMinute(in Time) (result int) {
	return (t.hours()-in.hours())*60 + (t.minutes() - in.minutes())
}

func (t Time) SubMinuteByString(in string) (result int) {
	tempTime := NewTime(0, 0, 0, 0)
	tempTime.setFromString(in)
	return t.SubMinute(tempTime)
}
