package user

type User struct {
	ID        string `gorm:"primaryKey;column:id"`
	Name      string `gorm:"column:name"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

type UserToken struct {
	UserID    string `gorm:"column:user_id"`
	Token     string `gorm:"column:token"`
	ExpiredAt int64  `gorm:"column:expired_at"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}
