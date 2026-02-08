package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/kiart-tantasi/crm-go/internal/contactlists"
	"github.com/kiart-tantasi/crm-go/internal/contacts"
	"github.com/kiart-tantasi/crm-go/internal/emails"
	"github.com/kiart-tantasi/crm-go/internal/env"
	"github.com/kiart-tantasi/crm-go/internal/httpclient"
	"github.com/kiart-tantasi/crm-go/internal/users"
)

/*
# Test with mock data

- go run cmd/api/main.go

- go run cmd/createmock/main.go
  - or CONTACT_COUNT=XXX go run cmd/createmock/main.go

- go run cmd/smtp-server/main.go

- Start sending email by `curl -X POST http://localhost:8080/emails/{email_id}/send` (replace {email_id} with email id below)
*/
func main() {
	// Env
	trueVal := true
	userId := -100
	emailID := -100
	contactListId := -100
	prefixContactId := 100_000
	contactCountStr := env.GetEnv("CONTACT_COUNT", "1000")

	// Cast string to int
	contactCount, err := strconv.Atoi(contactCountStr)
	if err != nil {
		log.Fatal(err)
	}

	// Shared HTTP client
	// TODO: experiment on singleton vs new client for every request
	client := httpclient.NewClient()

	// User
	user := users.User{
		ID:          userId,
		Username:    "testuser",
		Password:    "password123",
		Firstname:   sql.NullString{String: "Test", Valid: true},
		Lastname:    sql.NullString{String: "User", Valid: true},
		Email:       "test@example.com",
		IsPublished: &trueVal,
	}
	userBody, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	userReq, err := http.NewRequest("POST", "http://localhost:8080/users", bytes.NewBuffer(userBody))
	if err != nil {
		log.Fatal(err)
	}
	userReq.Header.Set("Content-Type", "application/json")
	userResp, err := client.Do(userReq)
	if err != nil {
		log.Fatal(err)
	} else {
		if userResp.StatusCode < 200 || userResp.StatusCode > 299 {
			printResponseError(userResp, user.Email)
			log.Fatal("Response is not 2xx")
		} else {
			log.Printf("User created: %d", userResp.StatusCode)
		}
	}
	userResp.Body.Close()

	// Email
	email := emails.Email{
		ID:          emailID,
		Alias:       "Default Email",
		Subject:     "Welcome to our service!",
		FromName:    sql.NullString{String: "Support Team", Valid: true},
		FromAddress: sql.NullString{String: "support@example.com", Valid: true},
		Template:    "Welcome to our service!",
		AddedBy:     userId,
		ModifiedBy:  userId,
	}
	body, err := json.Marshal(email)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8080/emails", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			printResponseError(resp, email.Alias)
			log.Fatal("Response is not 2xx")
		} else {
			log.Printf("Email created: %d", resp.StatusCode)
		}
	}
	defer resp.Body.Close()

	// Contacts
	var addedContactIds []int
	var mu sync.Mutex
	var wg sync.WaitGroup
	limit := make(chan interface{}, 10)

	for i := range contactCount {
		wg.Add(1)
		// simple semaphore
		limit <- i
		go func(i int) {
			defer func() {
				<-limit
				wg.Done()
			}()
			contactId := prefixContactId + i
			contact := contacts.Contact{
				ID:          contactId,
				Firstname:   fmt.Sprintf("Contact%d", i),
				Lastname:    "TestLastname",
				Email:       fmt.Sprintf("contact%d@example.com", i),
				IsPublished: &trueVal,
				AddedBy:     userId,
				ModifiedBy:  userId,
			}
			contactBody, err := json.Marshal(contact)
			if err != nil {
				log.Printf("Failed to marshal contact: %v", err)
				return
			}
			contactReq, err := http.NewRequest("POST", "http://localhost:8080/contacts", bytes.NewBuffer(contactBody))
			if err != nil {
				log.Printf("Failed to create request: %v", err)
				return
			}
			contactReq.Header.Set("Content-Type", "application/json")
			contactResp, err := client.Do(contactReq)
			if err != nil {
				log.Printf("Failed to send request: %v", err)
				return
			}
			defer contactResp.Body.Close()

			if contactResp.StatusCode < 200 || contactResp.StatusCode > 299 {
				printResponseError(contactResp, contact.Email)
				log.Fatal("Response is not 2xx")
			} else {
				log.Printf("Contact created: %s , %d", contact.Email, contactResp.StatusCode)
				mu.Lock()
				addedContactIds = append(addedContactIds, contactId)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	// Contact list
	contactList := contactlists.ContactList{
		ID:         contactListId,
		Name:       "Test List",
		AddedBy:    userId,
		ModifiedBy: userId,
	}
	contactListBody, err := json.Marshal(contactList)
	if err != nil {
		log.Fatal(err)
	}
	contactListReq, err := http.NewRequest("POST", "http://localhost:8080/contact-lists", bytes.NewBuffer(contactListBody))
	if err != nil {
		log.Fatal(err)
	}
	contactListReq.Header.Set("Content-Type", "application/json")
	contactListResp, err := client.Do(contactListReq)
	if err != nil {
		log.Fatal(err)
	} else {
		if contactListResp.StatusCode < 200 || contactListResp.StatusCode > 299 {
			printResponseError(contactListResp, contactList.Name)
			log.Fatal("Response is not 2xx")
		} else {
			log.Printf("Contact list created: %d", contactListResp.StatusCode)
		}
	}
	contactListResp.Body.Close()

	// Email - Contact list
	emailContactList := emails.BatchAddContactListsRequest{
		ContactListIDs: []int{contactListId},
		AddedBy:        userId,
	}
	emailContactListBody, err := json.Marshal(emailContactList)
	if err != nil {
		log.Fatal(err)
	}
	emailContactListReq, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/emails/%d/contact-lists", emailID), bytes.NewBuffer(emailContactListBody))
	if err != nil {
		log.Fatal(err)
	}
	emailContactListReq.Header.Set("Content-Type", "application/json")
	emailContactListResp, err := client.Do(emailContactListReq)
	if err != nil {
		log.Fatal(err)
	} else {
		if emailContactListResp.StatusCode < 200 || emailContactListResp.StatusCode > 299 {
			printResponseError(emailContactListResp, "Email-Contact List")
			log.Fatal("Response is not 2xx")
		} else {
			log.Printf("Email-Contact list relationship created: %d", emailContactListResp.StatusCode)
		}
	}
	emailContactListResp.Body.Close()

	// Contact list - Contacts
	contactListContacts := contactlists.BatchAddContactsRequest{
		ContactIDs: addedContactIds,
		AddedBy:    userId,
	}
	contactListContactsBody, err := json.Marshal(contactListContacts)
	if err != nil {
		log.Fatal(err)
	}
	contactListContactsReq, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/contact-lists/%d/contacts", contactListId), bytes.NewBuffer(contactListContactsBody))
	if err != nil {
		log.Fatal(err)
	}
	contactListContactsReq.Header.Set("Content-Type", "application/json")
	contactListContactsResp, err := client.Do(contactListContactsReq)
	if err != nil {
		log.Fatal(err)
	} else {
		if contactListContactsResp.StatusCode < 200 || contactListContactsResp.StatusCode > 299 {
			printResponseError(contactListContactsResp, "Contact List-Contacts")
			log.Fatal("Response is not 2xx")
		} else {
			log.Printf("Contact list-contacts relationship created: %d", contactListContactsResp.StatusCode)
		}
	}
	contactListContactsResp.Body.Close()
}

func printResponseError(resp *http.Response, label string) {
	var response map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
		log.Printf("Response error for %s: %v", label, response)
	}
}
