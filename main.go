package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
	basePath := os.Args[1]

	log.Println("trying to connect to database", os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Panicf("failed to connect to postgresql: %v", err)
	}

	defer conn.Close(context.Background())
	log.Println("connection succeeded")

	err = setupDatabase(conn)
	if err != nil {
		log.Panic(err)
	}

	f, err := os.Open(basePath)
	if err != nil {
		log.Panicf("failed to open %s: %v", basePath, err)
	}
	defer f.Close()

	log.Printf("using file: %s", basePath)
	processLines(conn, f)
}

func processLines(conn *pgx.Conn, f io.Reader) {
	// Scanner é capaz de ler de forma buferizada e já separa os resultados por linha por padrão
	scanner := bufio.NewScanner(f)

	// ignorar primeira linha
	scanner.Scan()

	log.Println("started processing data")

	succeeded := 0
	failed := 0
	total := 0

	batch := &pgx.Batch{}
	batchSize := 50

	// alocando uma única vez para evitar onerar o GC
	args := make([]any, 8)

	for scanner.Scan() {
		total++

		// lê uma linha e quebra em um slice de strings
		line := strings.ToUpper(scanner.Text())
		cols := strings.Fields(line)

		err := util.SanitizeColumns(cols, args)
		if err != nil {
			log.Printf("failed to sanitize customer data (%v): %v", cols, err)
			failed += 1
			continue
		}

		_ = batch.Queue("stmt-insert-customer", args...)

		// caso o batch atinja o tamanho definido, envia todas as queries e reseta o batch
		if batch.Len() == batchSize {
			n, err := sendBatch(conn, batch)
			if err != nil {
				log.Printf("failed to bulk insert: %v", err)
				failed += n
			} else {
				succeeded += n
			}
			batch = &pgx.Batch{}
		}
	}

	// envia o que pode ter sobrado no batch
	if batch.Len() > 0 {
		n, err := sendBatch(conn, batch)
		if err != nil {
			log.Printf("failed to bulk insert: %v", err)
			failed += n
		} else {
			succeeded += n
		}
	}

	log.Printf("finished processing data: %d succeeded, %d failed, %d total", succeeded, failed, total)
}

// cria a tabela customer se não existir e prepara o statement de inserção.
// não inicializo a conexão aqui, porque prefiro fazer o defer conn.Close na main.
func setupDatabase(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
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
		return fmt.Errorf("failed to create table customer: %v", err)
	}

	_, err = conn.Prepare(context.Background(), "stmt-insert-customer",
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
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}

	return nil
}

// envia as queries que foram enfileiradas e as executa.
// o batch não deve ser reutilizado.
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
