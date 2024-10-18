package user

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (User) Columns() []string {
	return []string{
		"id",
		"name",
		"email",
		"password",
		"created_at",
		"updated_at",
	}
}

type UserToken struct {
	UserID    string    `gorm:"column:user_id"`
	Token     string    `gorm:"column:token"`
	ExpiredAt time.Time `gorm:"column:expired_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}
