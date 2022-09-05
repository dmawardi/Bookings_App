package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// A custom form struct (embeds a url values object)
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if no errors present
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a new form struct
func NewForm(data url.Values) *Form {
	// Create a new form using data
	return &Form{
		data,
		// empty map
		errors(map[string][]string{}),
	}
}

// Required checks that field is not empty
// ... denotes multiple number of fields can be passed
func (f *Form) Required(fields ...string) {
	for _, fieldName := range fields {
		value := f.Get(fieldName)
		// If form has no value
		if strings.TrimSpace(value) == "" {
			f.Errors.AddErrorMessage(fieldName, "This field cannot be blank")
		}
	}
}

// Required checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	formField := f.Get(field)

	// If field empty
	if formField == "" {
		f.Errors.AddErrorMessage(field, "This field can not be blank")
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	formFieldData := f.Get(field)

	if len(formFieldData) < length {
		f.Errors.AddErrorMessage(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	// Validate whether the form field is an email
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.AddErrorMessage(field, "Invalid email address")
	}
}
