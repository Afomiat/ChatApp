package userUtil

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func GenerateOTP() string {
	rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator with int64
	otp := rand.Intn(900000) + 100000        // Generate a random number between 100000 and 999999
	return fmt.Sprintf("%06d", otp)          // Format the number as a 6-digit string
}

func HassPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword, password string) error {
    fmt.Println("Hashed Password:", hashedPassword)
    fmt.Println("Input Password:", password)

    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        fmt.Println("Error comparing passwords:", err)
        return fmt.Errorf("invalid password")
    }
    return nil
}
