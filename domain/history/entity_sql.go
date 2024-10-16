package history

import "time"

type History struct {
	ID           string    `gorm:"primaryKey;column:id"`
	Title        string    `gorm:"column:title"`
	URL          string    `gorm:"column:url"`
	UserID       string    `gorm:"column:user_id"`
	DeviceName   string    `gorm:"column:device_name"`
	LastActiveAt time.Time `gorm:"column:last_active_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (History) TableName() string {
	return "histories"
}
