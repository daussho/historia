package history

import "time"

type History struct {
	ID           string    `gorm:"primaryKey;column:id"`
	Title        string    `gorm:"column:title"`
	URL          string    `gorm:"column:url"`
	UserID       string    `gorm:"column:user_id"`
	DeviceName   string    `gorm:"column:device_name"`
	LastActiveAt time.Time `gorm:"column:last_active_at;default:CURRENT_TIMESTAMP()"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP()"`
}

func (History) TableName() string {
	return "histories"
}
