package main

import (
	"database/sql"
	"filemon/config"
	"filemon/internal/filemon"
	"filemon/internal/repository"
	"log"

	_ "modernc.org/sqlite"
)

func init() {

}

func main() {
	var cfg = config.MustReadConfig()

	dbh, err := sql.Open("sqlite", cfg.DB_DSN)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := dbh.Exec(`
	CREATE TABLE IF NOT EXISTS fileinfo(file_name varchar(128) constraint pk primary key, file_size int);
	`); err != nil {
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
