package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchData(t *testing.T) {
	// Mock API response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{"id": 1, "name": "John", "subData": map[string]string{"subkey1": "subvalue1"}}
		json.NewEncoder(w).Encode(data)
	}))
	defer ts.Close()

	// Fetch data
	client := NewClient()
	data, err := client.FetchData(ts.URL)

	// Check error (non 2xx status code)
	require.NoError(t, err)

	// Check root-level
	assert.Equal(t, "John", data["name"])
	assert.Equal(t, 1.0, data["id"])
	// Check sub-level
	subData := data["subData"].(map[string]interface{})
	assert.Equal(t, "subvalue1", subData["subkey1"])
}

func TestFetchDataError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewClient()
	_, err := client.FetchData(ts.URL)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
