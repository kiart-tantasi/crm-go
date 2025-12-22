package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	tpl := "Hello {{.name}}!"
	data := map[string]any{"name": "World"}
	got, err := Render(tpl, data)

	require.NoError(t, err)
	assert.Equal(t, "Hello World!", got)
}

func TestRenderInvalidTemplate(t *testing.T) {
	tpl := "Hello {{.name"
	data := map[string]any{"name": "World"}
	_, err := Render(tpl, data)

	assert.Error(t, err)
}
