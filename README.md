# Go Form Builder

A fluent, secure, and component-friendly HTML form builder for Go, inspired by `laravelcollective/html`. This package eliminates the hassle of writing repetitive HTML, handling CSRF tokens, managing validation errors, and repopulating forms after errors.

## Features

- **Fluent API:** Build forms with clean Go code: `Form.Text("name")`, `Form.Select("country", ...)`
- **Automatic State Management:** Automatically populates form values from old input, a bound model, or default values.
- **Built-in Security:** Automatic CSRF token injection and `_method` spoofing for `PUT`/`DELETE` requests.
- **Integrated Validation:** Uses `go-playground/validator/v10` under the hood. Automatically adds `is-invalid` classes and displays error messages.
- **Fully Featured:** Supports all standard HTML form elements, including `select` with `optgroup`, multi-checkboxes, and more.

## Installation

```bash
go get github.com/zatrano/form-builder
```

## Usage

**1. In your Handler (e.g., with Fiber):**

```go
import builder "github.com/zatrano/form-builder"

func (h *UserHandler) Edit(c *fiber.Ctx) error {
    user, _ := h.userService.GetByID(1) // Get your model
    
    // Create the builder, binding the model to the form
    form := builder.New(builder.Config{
        Action:    "/users/1",
        Method:    "PUT",
        CSRFToken: c.Locals("csrf").(string),
        Model:     user,
    })

    return c.Render("users/edit", fiber.Map{"Form": form})
}
```

**2. In your Template (`html/template`):**

```html
{{.Form.Open}}

<div class="mb-3">
    {{.Form.Label "name" "Full Name"}}
    {{.Form.Text "name" (dict "class" "my-custom-class")}}
    {{.Form.FieldError "name"}}
</div>

<div class="mb-3">
    {{.Form.Label "email" "Email Address"}}
    {{.Form.Email "email"}}
    {{.Form.FieldError "email"}}
</div>

{{.Form.Submit "Update User"}}
{{.Form.Close}}
```
