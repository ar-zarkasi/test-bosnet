package utils

import (
	"os"
	"time"

	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func ContainString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func IntDateToString(date int) string {
	timezone := os.Getenv("TZ")
	dateString := time.Unix(int64(date), 0).In(time.FixedZone(timezone, 0))
	return dateString.Format(time.RFC3339)
}

func DateToStringFormat(date time.Time, format string) string {
	timezone := os.Getenv("TZ")
	dateString := date.In(time.FixedZone(timezone, 0))
	return dateString.Format(format)
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func PickRandomInterface(arr []interface{}) interface{} {
	return arr[rand.Intn(len(arr))]
}

func IsInArithmeticSequence(number int) bool {
	// Suku pertama (a) dan beda (d) 
	a := 1 
	d := 3 
	// Jika (number - a) habis dibagi d, maka number ada dalam deret 
	if (number-a)%d == 0 { 
		return true 
	} 
	return false 
}