package util

import "fmt"

// TODO:
func ValidateCpfOrCnpj(val string) error {
	if len(val) == 11 {
		return ValidateCpf(val)
	}
	if len(val) == 14 {
		return ValidateCnpj(val)
	}
	return fmt.Errorf("invalid CPF or CNPJ with length %d: %s", len(val), val)
}

func ValidateCpf(val string) error {
	return nil
}

func ValidateCnpj(val string) error {
	return nil
}
