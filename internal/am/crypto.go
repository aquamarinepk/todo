package am

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEncryptionFailed = errors.New("encryption failed")
	ErrDecryptionFailed = errors.New("decryption failed")
)

type Crypto struct {
	key []byte
}

func NewCrypto(key ...[]byte) *Crypto {
	var k []byte
	if len(key) > 0 {
		k = key[0]
	}

	return &Crypto{key: k}
}

func (c *Crypto) SetKey(key []byte) {
	c.key = key
}

func (c *Crypto) EncryptEmail(email string) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, ErrEncryptionFailed
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, ErrEncryptionFailed
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, ErrEncryptionFailed
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(email), nil)
	return ciphertext, nil
}

func (c *Crypto) DecryptEmail(ciphertext []byte) (string, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", ErrDecryptionFailed
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

func (c *Crypto) HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (c *Crypto) CheckPassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
