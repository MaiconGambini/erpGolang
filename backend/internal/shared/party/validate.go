package party

import (
	"errors"
	"net/mail"
	"strings"
)

var ErrValidation = errors.New("validation error")

func ValidateInput(name string, email *string) error {
	if strings.TrimSpace(name) == "" {
		return ErrValidation
	}
	if email != nil && *email != "" {
		if _, err := mail.ParseAddress(*email); err != nil {
			return ErrValidation
		}
	}
	return nil
}

type PartyInput struct {
	Name         string
	DocumentType *string
	Document     *string
	Email        *string
	State        *string
}

func ValidatePartyInput(in PartyInput) error {
	if err := ValidateInput(in.Name, in.Email); err != nil {
		return err
	}
	if in.State != nil && *in.State != "" {
		state := strings.ToUpper(strings.TrimSpace(*in.State))
		if len(state) != 2 {
			return ErrValidation
		}
	}
	if in.DocumentType != nil && *in.DocumentType != "" && in.Document != nil && *in.Document != "" {
		if err := ValidateDocument(*in.DocumentType, *in.Document); err != nil {
			return ErrValidation
		}
	}
	return nil
}

func NormalizeDocumentPtr(doc *string) *string {
	if doc == nil || *doc == "" {
		return nil
	}
	normalized := NormalizeDocument(*doc)
	if normalized == "" {
		return nil
	}
	return &normalized
}

func NormalizeDocumentTypePtr(docType *string) *string {
	if docType == nil || *docType == "" {
		return nil
	}
	normalized := strings.ToLower(strings.TrimSpace(*docType))
	return &normalized
}

func NormalizeStatePtr(state *string) *string {
	if state == nil || *state == "" {
		return nil
	}
	normalized := strings.ToUpper(strings.TrimSpace(*state))
	return &normalized
}
