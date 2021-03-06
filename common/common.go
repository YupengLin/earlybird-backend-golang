package common

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	DB     *sql.DB
	Config config
)

type config struct {
	JwtSecret []byte
}

func init() {
	Config = config{
		JwtSecret: []byte(os.Getenv("JWT_SECRET")),
	}

	DB = getDb()
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123!@#456"
	dbname   = "earlybird"
)

func getDb() *sql.DB {
	log.Println("connecting to DB......")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panicln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("connected!")
	return db
}
