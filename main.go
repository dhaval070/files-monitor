package main

import (
	"database/sql"
	"filemon/config"
	"filemon/internal/filemon"
	"filemon/internal/repository"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var cfg = config.MustReadConfig()

	dbh, err := sql.Open("mysql", cfg.DB_DSN)
	if err != nil {
		log.Fatal(err)
	}
	repo := &repository.Repository{
		Dbh: dbh,
	}

	monitor, err := filemon.NewFileMon(cfg.DataDir)
	if err != nil {
		log.Fatal(err)
	}

	ch := monitor.Watch()

	for {
		select {
		case event := <-ch:
			if err := repo.Insert(event.Name, event.Size); err != nil {
				log.Println(err)
			}
			log.Printf("%+v\n", event)
		}
	}
}
