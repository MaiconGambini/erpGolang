package party

import (
	"errors"
	"strings"
	"unicode"
)

const (
	DocumentTypeCPF  = "cpf"
	DocumentTypeCNPJ = "cnpj"
)

var ErrInvalidDocument = errors.New("invalid document")

func NormalizeDocument(doc string) string {
	var b strings.Builder
	for _, r := range doc {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func ValidateCPF(cpf string) bool {
	digits := NormalizeDocument(cpf)
	if len(digits) != 11 || allSameDigits(digits) {
		return false
	}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(digits[i]-'0') * (10 - i)
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if int(digits[9]-'0') != d1 {
		return false
	}
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(digits[i]-'0') * (11 - i)
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	return int(digits[10]-'0') == d2
}

func ValidateCNPJ(cnpj string) bool {
	digits := NormalizeDocument(cnpj)
	if len(digits) != 14 || allSameDigits(digits) {
		return false
	}
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i, w := range weights1 {
		sum += int(digits[i]-'0') * w
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if int(digits[12]-'0') != d1 {
		return false
	}
	sum = 0
	for i, w := range weights2 {
		sum += int(digits[i]-'0') * w
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	return int(digits[13]-'0') == d2
}

func ValidateDocument(docType, doc string) error {
	normalized := NormalizeDocument(doc)
	switch strings.ToLower(strings.TrimSpace(docType)) {
	case DocumentTypeCPF:
		if !ValidateCPF(normalized) {
			return ErrInvalidDocument
		}
	case DocumentTypeCNPJ:
		if !ValidateCNPJ(normalized) {
			return ErrInvalidDocument
		}
	default:
		return ErrInvalidDocument
	}
	return nil
}

func allSameDigits(digits string) bool {
	if digits == "" {
		return true
	}
	first := digits[0]
	for i := 1; i < len(digits); i++ {
		if digits[i] != first {
			return false
		}
	}
	return true
}
