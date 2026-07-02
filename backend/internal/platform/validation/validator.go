package validation

type Validator struct{}

func New() Validator {
	return Validator{}
}

func (Validator) ValidateStruct(value any) error {
	_ = value
	return nil
}
