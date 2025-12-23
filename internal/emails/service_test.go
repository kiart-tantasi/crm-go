package emails

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	emailBody := "Hello {{.name}}!"
	data := map[string]any{"name": "World"}
	got, err := Render(emailBody, data)

	require.NoError(t, err)
	assert.Equal(t, "Hello World!", got)
}

func TestRenderInvalidEmail(t *testing.T) {
	emailBody := "Hello {{.name"
	data := map[string]any{"name": "World"}
	_, err := Render(emailBody, data)

	assert.Error(t, err)
}
