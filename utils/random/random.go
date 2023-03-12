package random

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	// jika kita tidak menggukanakan seed maka akan dianggap bahwa random kita seed nya 1
	rand.Seed(time.Now().UnixNano()) // unix nano memastikan bahwa setiap kita menjalankan kode, nilai yang digen akan berbeda
}

// random angka diantara min dan max
func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // mengembalikan int antara min dan max
}

const alphabet = "abcedfghijklmnopqrstuvwxyz"

// random string dengan length sepanjang n
func randomString(n int) string {
	var sb strings.Builder

	lenAlpha := len(alphabet)

	for i := 0; i < n; i++ {
		letter := alphabet[rand.Intn(lenAlpha)] // random posisi antara 0 - 26
		sb.WriteByte(letter)                    // append huruf ke sb string builder
	}

	return sb.String()

}

// generate random owner name
func RandomOwner() string {
	return randomString(6)
}

func RandomMoney() int64 {
	return randomInt(100, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "RP", "RUB"}
	lengthCurrencies := len(currencies)
	randomCurrencies := currencies[rand.Intn(lengthCurrencies)]
	return randomCurrencies
}

func RandomDescription() string {
	var desc string

	for i := 0; i <= int(randomInt(1, 4)); i++ {
		desc += randomString(5)
		desc += " "
	}

	return desc
}
