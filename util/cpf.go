package util

import (
	"fmt"
	"strconv"
)

// `ValidateCpfOrCnpj` espera uma string `val` já sanitizada e verifica se
// `val` representa um CPF ou CNPJ matematicamente válido
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
	// checa se o CPF possui algum dígito distinto ou se então é formado só
	// por um dígito repetido, por exemplo: 555.555.555-55
	hasDistinctDigits := false
	for i := 1; i < 11; i++ {
		if val[i-1] != val[i] {
			hasDistinctDigits = true
			break
		}
	}

	if !hasDistinctDigits {
		return fmt.Errorf("Invalid CPF with no distinct digit: %s", val)
	}

	sum := 0

	// Somando os 9 primeiros dígitos do CPF, aplicando os seguintes pesos:
	// (A * 10) + (B * 9) + (C * 8) + (D * 7) + (E * 6) + (F * 5) + (G * 4) + (H * 3) + (I * 2)
	// onde as letras repesentariam um CPF ABC.DEF.GHI-JK
	for i := 0; i < 9; i++ {
		digit, err := strconv.Atoi(val[i : i+1])
		if err != nil {
			return err
		}
		sum += (10 - i) * digit
	}

	wantDigit10 := strconv.Itoa(sum * 10 % 11 % 10) // dígito esperado
	gotDigit10 := val[9:10]                         // dígito obtido

	if wantDigit10 != gotDigit10 {
		return fmt.Errorf("expected digit at position 10 to be %s, but got %s", wantDigit10, gotDigit10)
	}

	sum = 0

	// Somando os 10 primeiros dígitos do CPF, aplicando os seguintes pesos:
	// (A * 11) + (B * 10) + (C * 9) + (D * 8) + (E * 7) + (F * 6) + (G * 5) + (H * 4) + (I * 3) + (J * 2)
	// onde as letras repesentariam um CPF ABC.DEF.GHI-JK
	for i := 0; i < 10; i++ {
		digit, err := strconv.Atoi(val[i : i+1])
		if err != nil {
			return err
		}
		sum += (11 - i) * digit
	}

	wantDigit11 := strconv.Itoa(sum * 10 % 11 % 10) // dígito esperado
	gotDigit11 := val[10:11]                        // dígito obtido

	if wantDigit11 != gotDigit11 {
		return fmt.Errorf("expected digit at position 11 to be %s, but got %s", wantDigit11, gotDigit11)
	}

	return nil
}

// TODO:
func ValidateCnpj(val string) error {
	return nil
}
