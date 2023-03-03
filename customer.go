package main

import "github.com/igoracmelo/neoway-desafio/util"

type Customer struct {
	Cpf                string
	Private            bool
	Incomplete         bool
	LastBoughtAt       *string
	TicketAverage      *string
	TicketLastPurchase *string
	StoreLastPurchase  *string
	StoreMostFrequent  *string
}

func CustomerFrom(args []string) (*Customer, error) {
	args[0] = util.SanitizeCpfOrCnpj(args[0])
	args[3] = util.SanitizeNullable(args[3])
	args[4] = util.SanitizeTicket(args[4])
	args[5] = util.SanitizeTicket(args[5])
	args[6] = util.SanitizeCpfOrCnpj(args[6])
	args[7] = util.SanitizeCpfOrCnpj(args[7])

	customer := &Customer{}

	if err := util.ValidateCpfOrCnpj(args[0]); err != nil {
		return nil, err
	}

	customer.Cpf = util.SanitizeCpfOrCnpj(args[0])
	customer.Private = args[1] == "1"
	customer.Incomplete = args[2] == "1"

	if args[3] != "" {
		customer.LastBoughtAt = &args[3]
	}

	if args[4] != "" {
		customer.TicketAverage = &args[4]
	}

	if args[5] != "" {
		customer.TicketLastPurchase = &args[5]
	}

	if args[6] != "" {
		if err := util.ValidateCpfOrCnpj(args[6]); err != nil {
			return nil, err
		}

		customer.StoreLastPurchase = &args[6]
	}

	if args[7] != "" {
		if err := util.ValidateCpfOrCnpj(args[7]); err != nil {
			return nil, err
		}

		customer.StoreMostFrequent = &args[7]
	}

	return customer, nil
}

func (c *Customer) ToArgs() []any {
	args := make([]any, 8)
	args[0] = c.Cpf
	args[1] = c.Private
	args[2] = c.Incomplete
	args[3] = c.LastBoughtAt
	args[4] = c.TicketAverage
	args[5] = c.TicketLastPurchase
	args[6] = c.StoreLastPurchase
	args[7] = c.StoreMostFrequent

	return args
}
