package history

type VisitRequest struct {
	Title      string `json:"title" validate:"required"`
	URL        string `json:"url" validate:"required"`
	UserID     int64  `json:"user_id" validate:"required,gt=0"`
	DeviceName string `json:"device_name" validate:"required"`
}
