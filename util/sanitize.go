package util

import (
	"unicode"
)

func SanitizeNullable(val string) string {
	if val == "NULL" {
		return ""
	}
	return val
}

// Remove todos os caracteres, exceto os numéricos.
// Optei por não usar regexp por ser um caso mais simples.
func SanitizeCpfOrCnpj(val string) string {
	if val == "NULL" {
		return ""
	}

	res := ""
	for _, r := range val {
		if unicode.IsDigit(r) {
			res += string(r)
		}
	}
	return res
}

// Remove ponto, se houver, e substitui vírgula por ponto
func SanitizeTicket(val string) string {
	if val == "NULL" {
		return ""
	}

	res := ""

	for _, r := range val {
		if r == '.' {
			continue
		} else if r == ',' {
			res += "."
		} else {
			res += string(r)
		}
	}

	return res
}
