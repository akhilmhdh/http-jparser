package jparser

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
)

var validate *validator.Validate

var (
	ErrValidatorError  = errors.New("VALIDATION ERROR")
	ErrJSONEncodingErr = ""
)

func init() {
	validate = validator.New()
}

type Response interface {
	GetStatusCode() int
}

type FieldError struct {
	Field     string `json:"field"`
	Err       string `json:"error"`
	FieldType string `json:"field_type"`
}

// ValidationErrors contains various validation errors happened to each field in an json object
type ValidationErrors []FieldError

func (v *ValidationErrors) Error() string {
	return ErrValidatorError.Error()
}

func validateFields(v interface{}) error {
	err := validate.Struct(v)

	if err == nil {
		return nil
	}

	var e ValidationErrors
	for _, vErr := range err.(validator.ValidationErrors) {
		e = append(e, FieldError{
			Field:     vErr.Field(),
			Err:       vErr.Error(),
			FieldType: vErr.Type().Name(),
		})
	}

	return &e
}

// Get function is used to get the JSON request body from http.Request  object
// After decoding request body validation is triggered using govalidator
func Get(r *http.Request, v interface{}) (err error) {
	defer io.Copy(ioutil.Discard, r.Body)
	err = json.NewDecoder(r.Body).Decode(&v)

	if err != nil {
		return err
	}

	err = validateFields(v)

	return
}

func encoderJson(w http.ResponseWriter, r *http.Request, v interface{}) (buf *bytes.Buffer, err error) {
	buf = &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err = enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// Send function is used to send JSON.
// The response must have a method that returns status code
// The struct can be something like a generalized return body
func Send(w http.ResponseWriter, r *http.Request, v Response) (err error) {
	buf, err := encoderJson(w, r, v)

	// write headers and content
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(v.GetStatusCode())
	w.Write(buf.Bytes())
	return err
}

// SendWithStatusCode function is used to send JSON with a status code instead of method.
func SendWithStatusCode(w http.ResponseWriter, r *http.Request, v interface{}, status int) (err error) {
	buf, err := encoderJson(w, r, v)

	// write headers and content
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
	return err
}
