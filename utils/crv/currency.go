package crv

// currency yang ingin kita support
const (
	USD = "USD"
	RUB = "RUB"
	RP  = "RP"
)

// IsSupportedCurrency mengembalikan true jika currency disupport
func IsSupportedCurrency(currency string) bool {
	switch currency {
	// jika currency nilainya salah satu dari yang dibawah ini maka kembalikan true
	case USD, RUB, RP:
		return true
	}
	return false
}
