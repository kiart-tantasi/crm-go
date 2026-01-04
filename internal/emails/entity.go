package emails

import "database/sql"

// column date_added and date_modified are handled in mysql db
type Email struct {
	ID          int            `json:"id"`
	Alias       string         `json:"alias" binding:"required"`
	Subject     string         `json:"subject" binding:"required"`
	FromName    sql.NullString `json:"fromName"`
	FromAddress sql.NullString `json:"fromAddress"`
	Template    string         `json:"template" binding:"required"`
	AddedBy     int            `json:"addedBy" binding:"required"`
	ModifiedBy  int            `json:"modifiedBy" binding:"required"`
}

type BatchAddContactListsRequest struct {
	ContactListIDs []int `json:"contactListIds" binding:"required"`
	AddedBy        int   `json:"addedBy" binding:"required"`
}

type BatchRemoveContactListsRequest struct {
	ContactListIDs []int `json:"contactListIds" binding:"required"`
}
