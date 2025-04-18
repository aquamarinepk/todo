package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// EncryptEmail encrypts an email using AES-GCM
func EncryptEmail(email string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(email), nil)
	return ciphertext, nil
}

// DecryptEmail decrypts an email using AES-GCM
func DecryptEmail(ciphertext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashPassword returns the bcrypt hash of the password.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CheckPassword compares the bcrypt-hashed password with the provided password.
func CheckPassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
