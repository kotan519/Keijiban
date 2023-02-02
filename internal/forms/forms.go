package forms

import (
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

// Form creates ac custom struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns tre if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (f *Form) Required(field ...string) {
	for _, field := range field {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if form field is in pot and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// MinLength checks for string minimum length

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "メールアドレスを入力してください")
	}
}

func (f *Form)EmailKosen(field string) {
	fl := f.Get(field)
	str := "@kurume.kosen-ac.jp"
	ok := strings.Contains(fl, str)
	if !ok {
		f.Errors.Add(field, "久留米高専のメールを入力してください")
	}
}
