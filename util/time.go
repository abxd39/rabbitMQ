package util

import (
	"time"

	"sctek.com/typhoon/th-platform-gateway/common"
)

type DateTime time.Time

func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+common.DATETIME+`"`, string(data), time.Local)
	*t = DateTime(now)
	return
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(common.DATETIME)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, common.DATETIME)
	b = append(b, '"')
	return b, nil
}

func (t DateTime) String(formatter string) string {
	return time.Time(t).Format(formatter)
}

func (t DateTime) DefaultString() string {
	return time.Time(t).Format(common.DATETIME)
}

func (t DateTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func InitDefaultTime(day int) DateTime {
	return DateTime(time.Now().AddDate(0, 0, day))
}
