package history

import "time"

type History struct {
	ID           string    `db:"id"`
	Title        string    `db:"title"`
	URL          string    `db:"url"`
	UserID       string    `db:"user_id"`
	DeviceName   string    `db:"device_name"`
	LastActiveAt time.Time `db:"last_active_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
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
