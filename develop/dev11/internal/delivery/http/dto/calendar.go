package dto

type EventCreateRequest struct {
	ID          string `json:"id"`
	CreatorID   string `json:"creator_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
