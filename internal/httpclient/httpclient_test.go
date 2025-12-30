package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchDataAndMap(t *testing.T) {
	// Mock API response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{"id": 1, "name": "John", "extraData": map[string]string{"location": "Thailand"}}
		json.NewEncoder(w).Encode(data)
	}))
	defer ts.Close()

	// Fetch data
	client := NewClient()
	data, err := client.FetchDataAndMap(ts.URL)

	// Check error (non 2xx status code)
	require.NoError(t, err)

	// Check root-level
	assert.Equal(t, "John", data["name"])
	assert.Equal(t, 1.0, data["id"])
	// Check sub-level
	subData := data["extraData"].(map[string]interface{})
	assert.Equal(t, "Thailand", subData["location"])
}

func TestFetchDataAndMapError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := NewClient()
	_, err := client.FetchDataAndMap(ts.URL)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status code: 500")
}
