package main

type StringValidation struct {
	Value string
	Valid bool
}

func WrapStringValidation(value string, valid bool) StringValidation {
	return StringValidation{
		Value: value,
		Valid: valid,
	}
}

type Uint8Validation struct {
	Value uint8
	Valid bool
}

func WrapUint8Validation(value uint8, valid bool) Uint8Validation {
	return Uint8Validation{
		Value: value,
		Valid: valid,
	}
}
