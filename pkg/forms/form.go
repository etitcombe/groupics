package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form represents an HTML form.
type Form struct {
	url.Values
	Errors errors
}

// New returns a new instance of a Form.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors{},
		//errors(map[string][]string{}),
	}
}

// Required checks each element of field and adds an error if its value is empty.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength checks field for length d. If it is longer it adds an error.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// PermittedValues checks that a specific field in the form matches one of a set
// of specific permitted values. If the check fails it adds an error.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Valid returns true if there are no errors, otherwise false.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
