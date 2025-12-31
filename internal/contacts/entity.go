package contacts

type Contact struct {
	ID        int    `json:"id"`
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	// NOTE: gin will give bad response when value is false so we need to use pointer here
	IsPublished *bool `json:"isPublished" binding:"required"`
	AddedBy     int   `json:"addedBy" binding:"required"`
	ModifiedBy  int   `json:"modifiedBy" binding:"required"`
}
