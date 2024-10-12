package history

type History struct {
	ID         int64  `gorm:"primaryKey;column:id"`
	Title      string `gorm:"column:title"`
	URL        string `gorm:"column:url"`
	UserID     int64  `gorm:"column:user_id"`
	DeviceName string `gorm:"column:device_name"`
	CreatedAt  int64  `gorm:"column:created_at"`
	UpdatedAt  int64  `gorm:"column:updated_at"`
}

func (History) TableName() string {
	return "histories"
}
