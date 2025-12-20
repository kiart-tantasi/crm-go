package template

import (
	"bytes"
	"fmt"
	"html/template"
)

func Render(bodyTemplate string, data any) (string, error) {
	// Create template
	tmpl, err := template.New("").Parse(bodyTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template with map-data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}
