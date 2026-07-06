package party

import (
	"testing"
)

func TestValidateInput(t *testing.T) {
	if err := ValidateInput("", nil); err == nil {
		t.Fatal("expected error for empty name")
	}
	if err := ValidateInput("   ", nil); err == nil {
		t.Fatal("expected error for whitespace name")
	}
	if err := ValidateInput("Acme", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	bad := "not-an-email"
	if err := ValidateInput("Acme", &bad); err == nil {
		t.Fatal("expected error for invalid email")
	}
	good := "user@example.com"
	if err := ValidateInput("Acme", &good); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	empty := ""
	if err := ValidateInput("Acme", &empty); err != nil {
		t.Fatalf("unexpected error for empty optional email: %v", err)
	}
}
