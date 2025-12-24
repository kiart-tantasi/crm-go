package contactlists

type ContactList struct {
	ID         int    `json:"id"`
	Name       string `json:"name" binding:"required"`
	AddedBy    int    `json:"addedBy" binding:"required"`
	ModifiedBy int    `json:"modifiedBy" binding:"required"`
}

type AddContactItem struct {
	ContactID int `json:"contactId" binding:"required"`
	AddedBy   int `json:"addedBy" binding:"required"`
}

type BatchAddContactsRequest struct {
	Contacts []AddContactItem `json:"contacts" binding:"required"`
}

type BatchRemoveContactsRequest struct {
	ContactIDs []int `json:"contactIds" binding:"required"`
}
