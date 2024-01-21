package main

import (
	"os"

	"github.com/Vodka479/go-shop-tutorial/config"
	"github.com/Vodka479/go-shop-tutorial/modules/servers"
	"github.com/Vodka479/go-shop-tutorial/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 { /*Args คือ argument*/
		return ".env"
	} else {
		return os.Args[1] /*1 คือ Args ตัวที่ 2*/
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.Newserver(cfg, db).Start()
}
