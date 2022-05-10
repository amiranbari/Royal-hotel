package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(make(errors)),
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s cannot be blank", strings.Replace(field, "_", " ", -1)))
		}
	}
}

func (f *Form) Has(field string) bool {
	value := f.Get(field)
	if strings.TrimSpace(value) == "" {
		f.Errors.Add(field, fmt.Sprintf("%s cannot be blank", field))
		return false
	}
	return true
}

func (f *Form) MinLength(field string, length int) bool {
	value := f.Get(field)
	if len(strings.TrimSpace(value)) < length {
		f.Errors.Add(field, fmt.Sprintf("field %s must be at least %d characters long", field, length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) bool {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "This is not an email address.")
		return false
	}
	return true
}
