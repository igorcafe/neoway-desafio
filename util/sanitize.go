package util

import (
	"strings"
	"unicode"
)

func SanitizeColumns(cols []string, res []any) error {
	for i := 0; i < len(res); i++ {
		res[0] = nil
	}

	cpf := cols[0]
	private := cols[1] == "1"
	incomplete := cols[2] == "1"
	lastBoughtAt := cols[3]
	ticketAverage := cols[4]
	ticketLastPurchase := cols[5]
	storeLastPurchase := cols[6]
	storeMostFrequent := cols[7]

	if cpf != "NULL" {
		cpf, err := SanitizeCpfOrCnpj(cpf)
		if err != nil {
			return err
		}
		res[0] = cpf
	}

	res[1] = private
	res[2] = incomplete

	if lastBoughtAt != "NULL" {
		res[3] = lastBoughtAt
	}

	if ticketAverage != "NULL" {
		res[4] = SanitizeTicket(ticketAverage)
	}

	if ticketLastPurchase != "NULL" {
		res[5] = SanitizeTicket(ticketLastPurchase)
	}

	if storeLastPurchase != "NULL" {
		storeLastPurchase, err := SanitizeCpfOrCnpj(storeLastPurchase)
		if err != nil {
			return err
		}
		res[6] = storeLastPurchase
	}

	if storeMostFrequent != "NULL" {
		storeMostFrequent, err := SanitizeCpfOrCnpj(storeMostFrequent)
		if err != nil {
			return err
		}
		res[7] = storeMostFrequent
	}

	return nil
}

// Remove todos os caracteres, exceto os numéricos.
// Optei por não usar regexp por ser um caso mais simples.
func SanitizeCpfOrCnpj(val string) (string, error) {
	res := &strings.Builder{}
	res.Grow(len(val))

	for _, r := range val {
		if unicode.IsDigit(r) {
			res.WriteRune(r)
		}
	}

	s := res.String()
	err := ValidateCpfOrCnpj(s)
	return s, err
}

// Remove ponto, se houver, e substitui vírgula por ponto
func SanitizeTicket(val string) string {
	res := &strings.Builder{}
	res.Grow(len(val))

	for _, r := range val {
		if r == '.' {
			continue
		} else if r == ',' {
			res.WriteRune('.')
		} else {
			res.WriteRune(r)
		}
	}

	return res.String()
}
