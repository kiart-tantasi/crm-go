package templates

import "time"

type Template struct {
	ID           int       `json:"id"`
	Name         string    `json:"name" binding:"required"`
	Subject      string    `json:"subject" binding:"required"`
	Body         string    `json:"body" binding:"required"`
	DateAdded    time.Time `json:"date_added"`
	DateModified time.Time `json:"date_modified"`
}
