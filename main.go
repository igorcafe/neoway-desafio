package main

import (
	"bufio"
	"context"
	"log"
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
		log.Println("running pprof for CPU")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Println("trying to connect to database", os.Getenv("DATABASE_URL"))
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		// TODO:
		panic(err)
	}
	defer db.Close(context.Background())

	log.Println("connection succeeded")

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

	log.Println("started processing data")

	for scanner.Scan() {
		cols := strings.Fields(scanner.Text())
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

	log.Println("finished processing data")
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
