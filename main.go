package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {

	for {
		err := run()
		fmt.Println("run() finished: ", err)
		time.Sleep(10 * time.Second)
	}
}

func run() error {
	ctx := context.Background()

	db, err := sqlx.Open("pgx", "host=localhost port=6432 user=postgres password=secret database=postgres sslmode=disable")
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(10)

	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(ctx, db, wg)
	}
	wg.Wait()

	err = db.Close()
	if err != nil {
		panic(err)
	}

	return nil
}

func worker(ctx context.Context, db *sqlx.DB, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		handle(ctx, db)
	}
}

func handle(ctx context.Context, db *sqlx.DB) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	_, _ = db.ExecContext(ctx, "select pg_sleep(10)")
}
