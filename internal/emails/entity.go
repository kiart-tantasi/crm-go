package emails

import "time"

type Email struct {
	ID           int       `json:"id"`
	Alias        string    `json:"alias" binding:"required"`
	Content      string    `json:"template" binding:"required"`
	DateAdded    time.Time `json:"date_added"`
	AddedBy      int       `json:"added_by"`
	DateModified time.Time `json:"date_modified"`
	ModifiedBy   int       `json:"modified_by"`
}
