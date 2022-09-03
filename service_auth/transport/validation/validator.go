package validation

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/go-playground/validator"
	"reflect"
)

func NewValidator() *validator.Validate {
	// use a single instance of Validate, it caches struct info
	var validate *validator.Validate
	validate = validator.New()

	// register all sql.Null* types to use the ValidateValuer CustomTypeFunc
	validate.RegisterCustomTypeFunc(ValidateValuer, sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})

	return validate
}

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
		// handle the error how you want
		return errors.New("value not found")
	}
	return nil
}
