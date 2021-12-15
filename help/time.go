package help

import "time"

func GetMS(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond()/1000000)
}
func GetNowMS() int64 {
	var addms int64
	return GetMS(time.Now()) + addms
}
