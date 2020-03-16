package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// GetEnv attempts to retrieve env variable from k being the key
// or returns d being the default value.
func GetEnv(k, d string) string {
	if len(k) <= 0 {
		return d
	}

	if v, ok := os.LookupEnv(k); ok {
		return v
	}

	return d
}

// ValidateInputs uses validator on interface to check each validation rule.
// If any fail it will then run through each and return a custom error message
// depending on which field was wrong.
func ValidateInputs(set interface{}) (bool, map[string][]string) {
	validate := validator.New()
	err := validate.Struct(set)

	if err != nil {
		//Validation syntax is invalid
		if err, ok := err.(*validator.InvalidValidationError); ok {
			panic(err)
		}

		//Validation errors occurred
		errors := make(map[string][]string)

		// Use reflector to reverse engineer struct
		r := reflect.TypeOf(set).Elem()
		for _, err := range err.(validator.ValidationErrors) {
			// Attempt to find field by name and get json tag name
			field, _ := r.FieldByName(err.StructField())
			var name string

			// If json tag doesn't exist, use lower case of name
			if name = field.Tag.Get("json"); name == "" {
				name = strings.ToLower(err.StructField())
			}

			switch err.Tag() {
			case "required":
				errors[name] = append(errors[name], "The "+name+" is required")
				break
			case "email":
				errors[name] = append(errors[name], "The "+name+" should be a valid email")
				break
			case "eqfield":
				errors[name] = append(errors[name], "The "+name+" should be equal to the "+err.Param())
				break
			default:
				errors[name] = append(errors[name], "The "+name+" is invalid")
				break
			}
		}
		return false, errors
	}
	return true, nil
}

// CreateValidationResponse is a helper function that creates a http response
// using the a map of validation errors.
func CreateValidationResponse(fields map[string][]string) ([]byte, error) {
	r := make(map[string]interface{})
	r["status"] = "error"
	r["message"] = "validation error"
	r["errors"] = fields
	message, err := json.Marshal(r)
	if err != nil {
		// An error occurred processing the json
		return nil, err
	}

	return message, nil
}

// ParseBody unmarshal's json body into interface type
func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}
