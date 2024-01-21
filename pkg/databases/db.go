package databases

import (
	"log"

	"github.com/Vodka479/go-shop-tutorial/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// Connect db
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("connect to db failed: %v\n", err)
	}
	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}
