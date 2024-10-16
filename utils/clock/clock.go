package clock

import "time"

func Now() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return time.Now().In(loc)
}
