package emails

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	// Test when fetch func map is called
	emailBody := "Hello world !"
	got, err := Render(emailBody)

	require.NoError(t, err)
	assert.Equal(t, "Hello world !", got)
}

func TestRenderFetch(t *testing.T) {
	// Mock API response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{"name": "FOOBAR", "extraData": map[string]string{"location": "The USA"}}
		json.NewEncoder(w).Encode(data)
	}))
	defer ts.Close()

	// Test when fetch func map is called
	emailBody := fmt.Sprintf("{{ $data := fetch \"%s\" }}Hello {{ $data.name }} ! You are from {{ $data.extraData.location }}", ts.URL)
	got, err := Render(emailBody)

	require.NoError(t, err)
	assert.Equal(t, "Hello FOOBAR ! You are from The USA", got)
}

func TestRenderInvalidEmail(t *testing.T) {
	emailBody := "Hello {{.name"
	_, err := Render(emailBody)
	assert.Error(t, err)
}
