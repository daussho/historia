package history

import "time"

type History struct {
	ID           string    `db:"id" gorm:"primaryKey;column:id"`
	Title        string    `db:"title" gorm:"column:title"`
	URL          string    `db:"url" gorm:"column:url"`
	UserID       string    `db:"user_id" gorm:"column:user_id"`
	DeviceName   string    `db:"device_name" gorm:"column:device_name"`
	LastActiveAt time.Time `db:"last_active_at" gorm:"column:last_active_at"`
	CreatedAt    time.Time `db:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `db:"updated_at" gorm:"column:updated_at"`
}

func (History) TableName() string {
	return "histories"
}

func (History) Columns() []string {
	return []string{
		"id",
		"title",
		"url",
		"user_id",
		"device_name",
		"last_active_at",
		"created_at",
		"updated_at",
	}
}
