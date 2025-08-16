package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUpHandler(t *testing.T) {
	//mockdb := &mockDB{}
	//h := New(zap.NewNop(), &mockdb.DB)
	//req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(`{}`))
	//w := httptest.NewRecorder()

	//h.SignUp(w, req)

	t.Skip("Integration test pending")
}

func TestHashPassword(t *testing.T) {
	zapLogger := zap.NewNop()
	mockdb := &mockDB{}

	h := New(zapLogger, &mockdb.DB)

	password := "securepassword"

	hash, err := h.HashPassword(password)

	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	if hash == "" {
		t.Fatal("Error empty hash string")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		t.Errorf("bcrypt.CompareHashAndPassword failed: %v", err)
	}

}

func TestCheckHashPassword(t *testing.T) {
	password := "securepassword"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	falsePassword := "wrongpassword"
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}
	assert.True(t, CheckPasswordHash(string(hashedPassword), password))
	assert.False(t, CheckPasswordHash(string(hashedPassword), falsePassword))

}
