package password

import (
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	password := "mysecretpassword"

	hash, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}
	if hash == "" {
		t.Fatal("Hash returned empty string")
	}
	if hash == password {
		t.Error("hash should not equal plain password")
	}
}

func TestHash_DifferentHashesForSamePassword(t *testing.T) {
	password := "mysecretpassword"

	hash1, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	hash2, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if hash1 == hash2 {
		t.Error("different hashes should be generated for same password due to salt")
	}
}

func TestVerify_CorrectPassword(t *testing.T) {
	password := "mysecretpassword"

	hash, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if !Verify(password, hash) {
		t.Error("Verify returned false for correct password")
	}
}

func TestVerify_IncorrectPassword(t *testing.T) {
	password := "mysecretpassword"
	wrongPassword := "wrongpassword"

	hash, err := Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	if Verify(wrongPassword, hash) {
		t.Error("Verify returned true for incorrect password")
	}
}

func TestVerify_EmptyPassword(t *testing.T) {
	hash, err := Hash("")
	if err != nil {
		t.Fatalf("Hash failed for empty password: %v", err)
	}

	if !Verify("", hash) {
		t.Error("Verify returned false for correct empty password")
	}
}

func TestVerify_InvalidHash(t *testing.T) {
	if Verify("password", "invalid-hash") {
		t.Error("Verify should return false for invalid hash format")
	}
}

func TestVerify_EmptyHash(t *testing.T) {
	if Verify("password", "") {
		t.Error("Verify should return false for empty hash")
	}
}

func TestHash_LongPassword(t *testing.T) {
	longPassword := strings.Repeat("a", 72)

	hash, err := Hash(longPassword)
	if err != nil {
		t.Fatalf("Hash failed for long password: %v", err)
	}

	if !Verify(longPassword, hash) {
		t.Error("Verify returned false for correct long password")
	}
}

func TestHash_SpecialCharacters(t *testing.T) {
	passwords := []string{
		"p@ssw0rd!#$%",
		"密码测试",
		"パスワード",
		"emoji🎉password",
		"  spaces  ",
		"tab\tpassword",
		"newline\npassword",
	}

	for _, password := range passwords {
		hash, err := Hash(password)
		if err != nil {
			t.Fatalf("Hash failed for password %q: %v", password, err)
		}

		if !Verify(password, hash) {
			t.Errorf("Verify returned false for correct password %q", password)
		}
	}
}
