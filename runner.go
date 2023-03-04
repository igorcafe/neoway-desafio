package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/igoracmelo/neoway-desafio/util"
	"github.com/jackc/pgx/v5"
)

type Runner struct {
	conn      *pgx.Conn
	batch     *pgx.Batch
	batchSize int
	file      io.ReadCloser
}

func NewRunner(dbUrl string, basePath string) (*Runner, error) {
	log.Println("trying to connect to database", os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}
	log.Println("connection to postgres succeeded")

	f, err := os.Open(basePath)
	if err != nil {
		return nil, err
	}
	log.Printf("using file: %s", basePath)

	runner := &Runner{
		conn:      conn,
		batch:     &pgx.Batch{},
		batchSize: 1000,
		file:      f,
	}

	return runner, nil
}

// encerra a conexão com o banco e fecha o arquivo
func (r *Runner) Close() error {
	err := r.conn.Close(context.Background())
	if err != nil {
		return err
	}

	err = r.file.Close()
	return err
}

// cria a tabela se não existir e prepara o statement de inserção
func (r *Runner) PrepareDatabase() error {

	// cria a tabela customer se não existir e prepara o statement de inserção.
	// não inicializo a conexão aqui, porque prefiro fazer o defer conn.Close na main.
	_, err := r.conn.Exec(context.Background(), `
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

	// cria um prepared statement no banco
	_, err = r.conn.Prepare(context.Background(), "stmt-insert-customer", `
		INSERT INTO Customer (
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

// lê e processa cada linha do arquivo `file` e insere os dados no banco
func (r *Runner) ProcessLines() {
	// Scanner é capaz de ler de forma buferizada e já separa os resultados por linha por padrão
	scanner := bufio.NewScanner(r.file)

	// ignorar primeira linha
	scanner.Scan()

	log.Println("started processing data")

	succeeded := 0
	failed := 0
	total := 0

	r.batch = &pgx.Batch{}

	for scanner.Scan() {
		total++

		// a otimização de alocar o slice antes do loop e reutilizá-lo não pode ser
		// usada, porque a query guarda a referência para o campo do slice, mas
		// quando a query é de fato executada o valor de args foi alterado, fazendo
		// com que um mesmo dado seja inserido múltiplas vezes a cada batch
		args := make([]any, 8)

		// lê uma linha e quebra em um slice de strings
		line := strings.ToUpper(scanner.Text())
		cols := strings.Fields(line)

		err := util.SanitizeColumns(cols, args)
		if err != nil {
			log.Printf("failed to sanitize customer data (%v): %v", cols, err)
			failed += 1
			continue
		}

		_ = r.batch.Queue("stmt-insert-customer", args...)

		// caso o batch atinja o tamanho definido, envia todas as queries e reseta o batch
		if r.batch.Len() == r.batchSize {
			n, err := r.sendBatch()
			if err != nil {
				log.Printf("failed to bulk insert: %v", err)
				failed += n
			} else {
				succeeded += n
			}
			r.batch = &pgx.Batch{}
		}
	}

	// envia o que pode ter sobrado no batch
	if r.batch.Len() > 0 {
		n, err := r.sendBatch()
		if err != nil {
			log.Printf("failed to bulk insert: %v", err)
			failed += n
		} else {
			succeeded += n
		}
	}

	log.Printf("finished processing data: %d succeeded, %d failed, %d total", succeeded, failed, total)
}

// envia as queries que foram enfileiradas e as executa.
// o batch não deve ser reutilizado.
func (r *Runner) sendBatch() (int, error) {
	res := r.conn.SendBatch(context.Background(), r.batch)

	for i := 0; i < r.batch.Len(); i++ {
		_, err := res.Exec()
		if err != nil {
			return i, err
		}
	}

	return r.batch.Len(), res.Close()
}
