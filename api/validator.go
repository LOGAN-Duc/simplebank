package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/techschool/simplebank/Utill"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return Utill.IsSupportedCurrency(currency)
	}
	return false
}
