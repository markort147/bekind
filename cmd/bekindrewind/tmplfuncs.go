package main

type StringValidation struct {
	Value   string
	Valid   bool
	Message string
}

func WrapStringValidation(value string, valid bool, message string) StringValidation {
	return StringValidation{
		Value:   value,
		Valid:   valid,
		Message: message,
	}
}

type Uint8Validation struct {
	Value   uint8
	Valid   bool
	Message string
}

func WrapUint8Validation(value uint8, valid bool, message string) Uint8Validation {
	return Uint8Validation{
		Value:   value,
		Valid:   valid,
		Message: message,
	}
}

type Uint16Validation struct {
	Value   uint16
	Valid   bool
	Message string
}

func WrapUint16Validation(value uint16, valid bool, message string) Uint16Validation {
	return Uint16Validation{
		Value:   value,
		Valid:   valid,
		Message: message,
	}
}
