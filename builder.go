package builder

import (
	"html/template"
	"net/url"
)

// Builder, bir HTML formu oluşturmak için gereken tüm durumu ve metodları içerir.
type Builder struct {
	model       interface{}
	oldInput    url.Values
	errors      map[string]string
	csrfToken   string
	csrfField   string
	action      string
	method      string
	isMultipart bool
}

// Config, yeni bir Builder oluşturmak için gerekli verileri taşır.
type Config struct {
	Action    string
	Method    string
	CSRFToken string
	CSRFField string
	Model     interface{}
	OldInput  url.Values
	Errors    map[string]string
	Multipart bool
}

// New, yeni bir form builder örneği oluşturur.
func New(config Config) *Builder {
	if config.OldInput == nil {
		config.OldInput = make(url.Values)
	}
	if config.Errors == nil {
		config.Errors = make(map[string]string)
	}
	if config.CSRFField == "" {
		config.CSRFField = "_csrf"
	}
	return &Builder{
		action:      config.Action,
		method:      config.Method,
		csrfToken:   config.CSRFToken,
		csrfField:   config.CSRFField,
		model:       config.Model,
		oldInput:    config.OldInput,
		errors:      config.Errors,
		isMultipart: config.Multipart,
	}
}