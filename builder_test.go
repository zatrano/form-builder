package builder

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

type TestForm struct {
	Name  string `form:"name" validate:"required"`
	Email string `form:"email"`
}

func TestFormOpen(t *testing.T) {
	form := New(Config{Action: "/test", Method: "POST", CSRFToken: "abc"})
	html := form.Open()
	assert.Contains(t, string(html), `action="/test"`)
	assert.Contains(t, string(html), `method="POST"`)
	assert.Contains(t, string(html), `name="_csrf" value="abc"`)
}

func TestTextInputWithValueFromModel(t *testing.T) {
	model := TestForm{Name: "John Doe"}
	form := New(Config{Model: &model})
	html := form.Text("name")
	assert.Contains(t, string(html), `value="John Doe"`)
}

func TestTextInputWithValueFromOldInput(t *testing.T) {
	model := TestForm{Name: "John Doe"} // Model var ama old input öncelikli olmalı
	oldInput := url.Values{"name": {"Jane Doe"}}
	form := New(Config{Model: &model, OldInput: oldInput})
	html := form.Text("name")
	assert.Contains(t, string(html), `value="Jane Doe"`)
}

func TestTextInputWithError(t *testing.T) {
	errors := map[string]string{"name": "Name is required"}
	form := New(Config{Errors: errors})
	html := form.Text("name")
	assert.Contains(t, string(html), `class="form-control is-invalid"`)
}

func TestSelectWithOptions(t *testing.T) {
	options := []Option{
		{Value: "1", Text: "Admin"},
		{Value: "2", Text: "User"},
	}
	oldInput := url.Values{"role": {"2"}}
	form := New(Config{OldInput: oldInput})
	html := form.Select("role", options)
	assert.Contains(t, string(html), `<option value="2" selected>User</option>`)
}