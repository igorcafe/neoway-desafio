package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/igoracmelo/neoway-desafio/util"
	"github.com/jackc/pgx/v5"
)

func main() {
	if os.Getenv("PPROF") == "cpu" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Panicf("failed to create cpu.prof: %v", err)
		}
		defer f.Close()
		log.Println("pprof: profiling CPU")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if len(os.Args) != 2 {
		log.Panicf("Missing positional arg FILE. Try %s some_file.txt", os.Args[0])
	}

	log.Println("trying to connect to database", os.Getenv("DATABASE_URL"))
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panicf("failed to connect to postgresql: %v", err)
	}
	defer db.Close(context.Background())

	log.Println("connection succeeded")

	_, err = db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS customer (
			id                        SERIAL PRIMARY KEY,
			cpf                       TEXT,
			private                   BOOLEAN,
			incomplete                BOOLEAN,
			last_bought_at            DATE,
			ticket_average            DECIMAL,
			ticket_last_purchase      DECIMAL,
			cnpj_most_frequent_store  TEXT,
			cnpj_last_purchase_store  TEXT
		);
	`)
	if err != nil {
		log.Panicf("failed to create table customer: %v", err)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Panicf("failed to open %s: %v", os.Args[1], err)
	}

	log.Printf("using file: %s", os.Args[1])

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
		log.Panicf("failed to prepare insert statement: %v", err)
	}

	scanner := bufio.NewScanner(f)

	// ignorar primeira linha
	scanner.Scan()

	log.Println("started processing data")

	succeeded := 0
	failed := 0
	total := 0

	batch := &pgx.Batch{}
	batchSize := 50

	for scanner.Scan() {
		total++

		cols := strings.Fields(strings.ToUpper(scanner.Text()))
		args, err := util.SanitizeColumns(cols)
		if err != nil {
			log.Printf("failed to sanitize customer data (%v): %v", cols, err)
			continue
		}

		_ = batch.Queue("stmt-insert-customer", args...)

		if total%batchSize == 0 {
			n, err := sendBatch(db, batch)
			if err != nil {
				log.Printf("failed to bulk insert: %v", err)
				failed += n
			} else {
				succeeded += n
			}
			batch = &pgx.Batch{}
		}
	}

	n, err := sendBatch(db, batch)
	if err != nil {
		log.Printf("failed to bulk insert: %v", err)
		failed += n
	} else {
		succeeded += n
	}

	log.Printf("finished processing data: %d succeeded, %d failed, %d total", succeeded, failed, total)
}

func sendBatch(conn *pgx.Conn, batch *pgx.Batch) (int, error) {
	res := conn.SendBatch(context.Background(), batch)

	for i := 0; i < batch.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			return i, err
		}
	}

	return batch.Len(), res.Close()
}
