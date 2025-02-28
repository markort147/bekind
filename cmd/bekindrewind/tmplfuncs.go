package main

type YearValidation struct {
	Year  string
	Valid bool
}

func WrapYearValidation(year string, valid bool) YearValidation {
	return YearValidation{
		Year:  year,
		Valid: valid,
	}
}
