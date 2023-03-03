package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
)

func main() {
	if os.Getenv("PPROF") == "cpu" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		// TODO:
		panic(err)
	}
	defer db.Close(context.Background())

	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS Customer (
			id SERIAL PRIMARY KEY,
			cpf VARCHAR(18) NOT NULL,
			private BOOLEAN NOT NULL,
			incomplete BOOLEAN NOT NULL,
			last_bought_at DATE,
			ticket_average NUMERIC(10, 2),
			ticket_last_purchase NUMERIC(10, 2),
			cnpj_most_frequent_store CHAR(18),
			cnpj_last_purchase_store CHAR(18)
		);
	`)
	if err != nil {
		// TODO:
		panic(err)
	}

	f, err := os.Open("base_teste.txt")
	if err != nil {
		// TODO:
		panic(err)
	}

	_, err = db.Prepare(context.Background(), "stmt-insert-customer",
		`INSERT INTO Customer (
			cpf,
			private,
			incomplete,
			last_bought_at,
			ticket_average,
			ticket_last_purchase,
			cnpj_most_frequent_store,
			cnpj_last_purchase_store
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`)

	if err != nil {
		// TODO:
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	// ignorar primeira linha
	scanner.Scan()

	for scanner.Scan() {
		cols := strings.Fields(scanner.Text())
		fmt.Println(cols)
		customer, err := CustomerFrom(cols)
		if err != nil {
			// TODO:
			panic(err)
		}

		_, err = db.Exec(context.Background(), "stmt-insert-customer", customer.ToArgs()...)
		if err != nil {
			// TODO:
			panic(err)
		}
	}
}

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
	args[0] = SanitizeCpfOrCnpj(args[0])
	args[3] = SanitizeNullable(args[3])
	args[4] = SanitizeTicket(args[4])
	args[5] = SanitizeTicket(args[5])
	args[6] = SanitizeCpfOrCnpj(args[6])
	args[7] = SanitizeCpfOrCnpj(args[7])

	customer := &Customer{}

	if err := ValidateCpfOrCnpj(args[0]); err != nil {
		return nil, err
	}

	customer.Cpf = SanitizeCpfOrCnpj(args[0])
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
		if err := ValidateCpfOrCnpj(args[6]); err != nil {
			return nil, err
		}

		customer.StoreLastPurchase = &args[6]
	}

	if args[7] != "" {
		if err := ValidateCpfOrCnpj(args[7]); err != nil {
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

// TODO:
func ValidateCpfOrCnpj(val string) error {
	if len(val) != 11 && len(val) != 14 {
		return fmt.Errorf("invalid CNPJ: %s", val)
	}
	return nil
}

func SanitizeNullable(val string) string {
	if val == "NULL" {
		return ""
	}
	return val
}

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

func SanitizeTicket(val string) string {
	if val == "NULL" {
		return ""
	}

	res := ""

	// remove `.` e substitui `,` por `.`
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
