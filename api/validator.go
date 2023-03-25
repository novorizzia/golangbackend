package api

import (
	"backendmaster/utils/crv"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// to get value dari field dan konversi valuenya menjadi string
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		// jika ok true, maka currency adalah string yang valid

		// check currency supported or not
		return crv.IsSupportedCurrency(currency)

	}
	return false

}
