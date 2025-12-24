package emails

// column date_added and date_modified are handled in mysql db
type Email struct {
	ID         int    `json:"id"`
	Alias      string `json:"alias" binding:"required"`
	Template   string `json:"template" binding:"required"`
	AddedBy    int    `json:"addedBy" binding:"required"`
	ModifiedBy int    `json:"modifiedBy" binding:"required"`
}
