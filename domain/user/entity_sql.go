package user

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;column:id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiredAt time.Time `db:"expired_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}

func (UserToken) Columns() []string {
	return []string{
		"user_id",
		"token",
		"expired_at",
		"created_at",
		"updated_at",
	}
}
