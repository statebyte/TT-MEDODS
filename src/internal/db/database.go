package db

import (
	"backend/src/config"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type DatabaseInstance struct {
	info config.DatabaseInfo
	Conn *sqlx.DB
}

func (db *DatabaseInstance) Connect(info config.DatabaseInfo) error {
	db.info = info
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db.info.USER, db.info.PASS, db.info.HOST, db.info.PORT, db.info.NAME)

	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Println(err)
		return err
	}

	db.Conn = conn

	return db.UpMigrations()
}

func (db *DatabaseInstance) UpMigrations() error {
	if err := goose.Up(db.Conn.DB, "./migrations"); err != nil {
		log.Fatalf("Error applying migrations: %v", err)
		return err
	}

	return nil
}
