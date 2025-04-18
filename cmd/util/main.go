package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/aquamarinepk/todo/internal/feat/auth"
)

type User struct {
	Username string
	Email    string
	Password string
}

func main() {
	encKey := []byte("8af0b8e0f14c4842b3e8f2dc41cf2872")

	users := []User{
		{
			Username: "johndoe",
			Email:    "john.doe@example.com",
			Password: "password123",
		},
		{
			Username: "janesmith",
			Email:    "jane.smith@example.com",
			Password: "password123",
		},
		{
			Username: "bobjohnson",
			Email:    "bob.johnson@example.com",
			Password: "password123",
		},
		{
			Username: "alicebrown",
			Email:    "alice.brown@example.com",
			Password: "password123",
		},
		{
			Username: "charliewilson",
			Email:    "charlie.wilson@example.com",
			Password: "password123",
		},
	}

	fmt.Println("-- Encrypted values for seed file --")
	fmt.Println("-- Copy these values into your seed file --")
	fmt.Println()

	for _, user := range users {
		encryptedEmail, err := auth.EncryptEmail(user.Email, encKey)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encrypting email for %s: %v\n", user.Username, err)
			continue
		}

		emailHex := hex.EncodeToString(encryptedEmail)
		passwordHex := hex.EncodeToString([]byte(user.Password))

		fmt.Printf("-- User: %s\n", user.Username)
		fmt.Printf("Email: X'%s'\n", emailHex)
		fmt.Printf("Password: X'%s'\n", passwordHex)
		fmt.Println()
	}
}
