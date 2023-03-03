package util

import (
	"strings"
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

	res := &strings.Builder{}
	res.Grow(len(val))

	for _, r := range val {
		if unicode.IsDigit(r) {
			res.WriteRune(r)
		}
	}
	return res.String()
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
