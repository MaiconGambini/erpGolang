package party

import "testing"

func TestNormalizeDocument(t *testing.T) {
	got := NormalizeDocument("12.345.678/0001-99")
	if got != "12345678000199" {
		t.Fatalf("expected digits only, got %q", got)
	}
}

func TestValidateCPF(t *testing.T) {
	if !ValidateCPF("529.982.247-25") {
		t.Fatal("expected valid CPF")
	}
	if ValidateCPF("111.111.111-11") {
		t.Fatal("expected invalid CPF")
	}
}

func TestValidateCNPJ(t *testing.T) {
	if !ValidateCNPJ("12.345.678/0001-95") {
		t.Fatal("expected valid CNPJ")
	}
	if ValidateCNPJ("11.111.111/1111-11") {
		t.Fatal("expected invalid CNPJ")
	}
}

func TestValidateDocument(t *testing.T) {
	if err := ValidateDocument("cpf", "529.982.247-25"); err != nil {
		t.Fatalf("expected valid cpf, got %v", err)
	}
	if err := ValidateDocument("cnpj", "12.345.678/0001-95"); err != nil {
		t.Fatalf("expected valid cnpj, got %v", err)
	}
	if err := ValidateDocument("cpf", "000.000.000-00"); err == nil {
		t.Fatal("expected invalid cpf")
	}
}
