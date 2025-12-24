package contacts

type Contact struct {
	ID          int    `json:"id"`
	Firstname   string `json:"firstname" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Email       string `json:"email" binding:"required"`
	IsPublished *bool  `json:"is_published" binding:"required"`
	AddedBy     int    `json:"added_by" binding:"required"`
	ModifiedBy  int    `json:"modified_by" binding:"required"`
}
