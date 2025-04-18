package auth

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptEmail(t *testing.T) {
	key := bytes.Repeat([]byte("k"), 32) // 32 bytes key (not secure, just for test)
	cases := []struct {
		email string
	}{
		{"user@example.com"},
		{"test+alias@domain.org"},
		{""},
	}

	for _, c := range cases {
		enc, err := EncryptEmail(c.email, key)
		if err != nil {
			t.Fatalf("encryption failed for %q: %v", c.email, err)
		}

		dec, err := DecryptEmail(enc, key)
		if err != nil {
			t.Fatalf("decryption failed for %q: %v", c.email, err)
		}

		if dec != c.email {
			t.Errorf("expected %q, got %q", c.email, dec)
		}
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	cases := []struct {
		password string
	}{
		{"securePassword123"},
		{"password!@#"},
		{" "},
		{""},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.password)
		if err != nil {
			t.Fatalf("hashing failed for %q: %v", c.password, err)
		}

		if err := CheckPassword(hash, c.password); err != nil {
			t.Errorf("password check failed for %q: %v", c.password, err)
		}
	}
}
