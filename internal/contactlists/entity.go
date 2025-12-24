package contactlists

type ContactList struct {
	ID         int    `json:"id"`
	Name       string `json:"name" binding:"required"`
	AddedBy    int    `json:"added_by" binding:"required"`
	ModifiedBy int    `json:"modified_by" binding:"required"`
}

type AddContactItem struct {
	ContactID int `json:"contact_id" binding:"required"`
	AddedBy   int `json:"added_by" binding:"required"`
}

type BatchAddContactsRequest struct {
	Contacts []AddContactItem `json:"contacts" binding:"required"`
}

type BatchRemoveContactsRequest struct {
	ContactIDs []int `json:"contact_ids" binding:"required"`
}
