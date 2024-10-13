package history

type History struct {
	ID           string `gorm:"primaryKey;column:id"`
	Title        string `gorm:"column:title"`
	URL          string `gorm:"column:url"`
	UserID       string `gorm:"column:user_id"`
	DeviceName   string `gorm:"column:device_name"`
	LastActiveAt int64  `gorm:"column:last_active_at"`
	CreatedAt    int64  `gorm:"column:created_at"`
	UpdatedAt    int64  `gorm:"column:updated_at"`
}

func (History) TableName() string {
	return "histories"
}
