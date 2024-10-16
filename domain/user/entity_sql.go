package user

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP()"`
}

func (User) TableName() string {
	return "users"
}

type UserToken struct {
	UserID    string    `gorm:"column:user_id"`
	Token     string    `gorm:"column:token"`
	ExpiredAt time.Time `gorm:"column:expired_at;default:CURRENT_TIMESTAMP()"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP()"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}
