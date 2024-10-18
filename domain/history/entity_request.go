package history

type VisitRequest struct {
	UserID     string
	Title      string `json:"title" validate:"required"`
	URL        string `json:"url" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
}
