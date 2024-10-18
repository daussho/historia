package user

import "time"

type UserSession struct {
	User
	ExpiredAt time.Time
}
