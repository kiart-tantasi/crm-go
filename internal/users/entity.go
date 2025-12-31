package users

import "database/sql"

type User struct {
	ID          int            `json:"id"`
	Username    string         `json:"username" binding:"required"`
	Password    string         `json:"password" binding:"required"`
	Firstname   sql.NullString `json:"firstname" binding:"required"`
	Lastname    sql.NullString `json:"lastname" binding:"required"`
	Email       string         `json:"email" binding:"required"`
	IsPublished *bool          `json:"isPublished" binding:"required"`
	AddedBy     sql.NullInt64  `json:"addedBy"`
	ModifiedBy  sql.NullInt64  `json:"modifiedBy"`
}
