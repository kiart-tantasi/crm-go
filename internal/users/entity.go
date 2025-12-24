package users

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Firstname   string `json:"firstname" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Email       string `json:"email" binding:"required"`
	IsPublished *bool  `json:"isPublished" binding:"required"`
	AddedBy     int    `json:"addedBy" binding:"required"`
	ModifiedBy  int    `json:"modifiedBy" binding:"required"`
}
