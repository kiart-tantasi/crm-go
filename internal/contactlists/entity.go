package contactlists

type ContactList struct {
	ID         int    `json:"id"`
	Name       string `json:"name" binding:"required"`
	AddedBy    int    `json:"addedBy" binding:"required"`
	ModifiedBy int    `json:"modifiedBy" binding:"required"`
}

type BatchAddContactsRequest struct {
	ContactIDs []int `json:"contactIds" binding:"required"`
	AddedBy    int   `json:"addedBy" binding:"required"`
}

type BatchRemoveContactsRequest struct {
	ContactIDs []int `json:"contactIds" binding:"required"`
}
