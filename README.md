# Go Form Builder (zatrano/form-builder)

<p align="center">
  <a href="https://pkg.go.dev/github.com/zatrano/form-builder"><img src="https://img.shields.io/badge/go.dev-reference-007d9c?style=flat-square&logo=go" alt="Go Doc"></a>
  <a href="https://github.com/zatrano/form-builder/actions/workflows/test.yml"><img src="https://img.shields.io/github/actions/workflow/status/zatrano/form-builder/test.yml?branch=main&style=flat-square" alt="Build Status"></a>
  <a href="https://goreportcard.com/report/github.com/zatrano/form-builder"><img src="https://goreportcard.com/badge/github.com/zatrano/form-builder?style=flat-square" alt="Go Report Card"></a>
  <a href="https://github.com/zatrano/form-builder/blob/main/LICENSE"><img src="https://img.shields.io/github/license/zatrano/form-builder?style=flat-square" alt="License"></a>
</p>

A fluent, secure, and component-friendly HTML form builder for Go, inspired by `laravelcollective/html`. This package eliminates the hassle of writing repetitive HTML, handling CSRF tokens, managing validation errors, and repopulating forms after failed submissions.

## ✨ Features

- **Fluent API:** Build forms with clean and expressive Go code in your templates: `.Form.Text("name")`, `.Form.Select("country", ...)`
- **Automatic State Management:** Automatically populates form values with the correct data, following this priority:
  1.  Old Input (after a validation error)
  2.  Bound Model (for editing forms)
  3.  Default/empty value
- **Built-in Security:** Automatic CSRF token injection via `Form.Open()` and seamless integration with any CSRF middleware.
- **Integrated Validation:** Designed to work with `go-playground/validator/v10`. Automatically adds `is-invalid` classes and displays error messages with `Form.FieldError("name")`.
- **Fully Featured:** Supports all standard HTML form elements, including `select` with `<optgroup>`, multi-checkboxes, radios, file inputs, and more.
- **Framework Agnostic Core:** The core builder logic is independent of any web framework, but it's exceptionally easy to use with frameworks like [Fiber](https://gofiber.io/).

## 🚀 Installation

```bash
go get github.com/zatrano/form-builder
```

## 📋 Quick Start Guide (with Fiber)

This guide demonstrates how to use `zatrano/form-builder` in a standard Go web application using the Fiber framework.

### 1. Define Your Form Struct

In your application, define a struct with `form` and `validate` tags.

```go
// app/forms/contact_form.go
package forms

type ContactForm struct {
    Name    string `form:"name" validate:"required"`
    Email   string `form:"email" validate:"required,email"`
    Message string `form:"message" validate:"required,min=10"`
}
```

### 2. Configure Your Handler

Your HTTP handler will be responsible for creating and processing the form builder.

```go
package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	builder "github.com/zatrano/form-builder"
	"your-app/app/forms"
)

// ShowContactForm handles GET requests to display the form.
func ShowContactForm(c *fiber.Ctx) error {
	// Create a new form builder for a new (empty) form.
	form := builder.New(builder.Config{
		Action:    "/contact",
		Method:    "POST",
		CSRFToken: c.Locals("csrf").(string), // Assuming CSRF middleware is set up
	})

	return c.Render("contact", fiber.Map{
		"Title": "Contact Us",
		"Form":  form,
	})
}

// HandleContactForm handles POST requests to process the form.
func HandleContactForm(c *fiber.Ctx) error {
	var formModel forms.ContactForm
	c.BodyParser(&formModel)

	// Validate the struct.
	errors, err := builder.Validate(&formModel)

	// If validation fails...
	if err != nil {
		formValues, _ := c.FormValues()

		// Re-create the builder with errors and old input.
		form := builder.New(builder.Config{
			Action:    "/contact",
			Method:    "POST",
			CSRFToken: c.Locals("csrf").(string),
			Errors:    errors,
			OldInput:  formValues,
		})

		// Render the form again with the errors and repopulated data.
		return c.Render("contact", fiber.Map{
			"Title": "Contact Us",
			"Form":  form,
		})
	}

	// Validation successful!
	// ... (process the data, send email, save to DB, etc.) ...

	return c.Redirect("/success")
}
```

### 3. Build Your Form in the View

In your `html/template` file, use the `Form` object to fluently build your HTML.

```html
<!-- views/contact.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css">
</head>
<body>
    <div class="container py-5">
        <h1>Contact Us</h1>

        {{.Form.Open}}
        
        <div class="mb-3">
            {{.Form.Label "name" "Your Name"}}
            {{.Form.Text "name" (dict "class" "form-control" "placeholder" "John Doe")}}
            {{.Form.FieldError "name"}}
        </div>

        <div class="mb-3">
            {{.Form.Label "email" "Your Email"}}
            {{.Form.Email "email" (dict "class" "form-control" "placeholder" "you@example.com")}}
            {{.Form.FieldError "email"}}
        </div>

        <div class="mb-3">
            {{.Form.Label "message" "Your Message"}}
            {{.Form.Textarea "message" (dict "class" "form-control" "rows" "5")}}
            {{.Form.FieldError "message"}}
        </div>
        
        {{.Form.Submit "Send Message"}}
        
        {{.Form.Close}}
    </div>
</body>
</html>
```

## API Reference

### Main Functions

- `builder.New(config Config) *Builder`: Creates a new form builder instance.
- `builder.Validate(s interface{}) (map[string]string, error)`: Validates any struct with `validate` tags.

### Builder Methods

- `.Open()`: Renders the opening `<form>` tag with CSRF and method spoofing.
- `.Close()`: Renders the closing `</form>` tag.
- `.Label(name, text, attrs...)`
- `.Text(name, attrs...)`
- `.Email(name, attrs...)`
- `.Password(name, attrs...)`
- `.Textarea(name, attrs...)`
- `.Select(name, options, attrs...)`
- `.Checkbox(name, value, attrs...)`
- `.Radio(name, value, attrs...)`
- `.File(name, attrs...)`
- `.Hidden(name, attrs...)`
- `.Submit(text, attrs...)`
- `.Button(text, attrs...)`
- `.FieldError(name)`: Renders the validation error message for a specific field.

All element methods accept an optional `map[string]string` to add custom HTML attributes.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the [issues page](https://github.com/zatrano/form-builder/issues).

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
