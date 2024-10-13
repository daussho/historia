package history

type VisitRequest struct {
	Title      string `json:"title" validate:"required"`
	URL        string `json:"url" validate:"required"`
	Token      string `json:"token" validate:"required"`
	DeviceName string `json:"device_name" validate:"required"`
}
