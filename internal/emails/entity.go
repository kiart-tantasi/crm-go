package emails

// column date_added and date_modified are handled in mysql db
type Email struct {
	ID         int    `json:"id"`
	Alias      string `json:"alias" binding:"required"`
	Template   string `json:"template" binding:"required"`
	AddedBy    int    `json:"added_by" binding:"required"`
	ModifiedBy int    `json:"modified_by" binding:"required"`
}
