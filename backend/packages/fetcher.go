package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type BlogDBFetcher interface {
	OpenConnection() (db *sql.DB, userId string)
	GetEmail() string
}

type BlogFetcher struct{}

func (fetcher *BlogFetcher) OpenConnection() (*sql.DB, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user, ok := os.LookupEnv("USER")
	if !ok {
		log.Fatal("Error loading env variables")
	}
	password, ok := os.LookupEnv("PASSWORD")
	if !ok {
		log.Fatal("Error loading env variables")
	}
	dbname, ok := os.LookupEnv("DB_NAME")
	if !ok {
		log.Fatal("Error loading env variables")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// lets implement auth later...
	// email := GetEmail()
	// addEmail := `INSERT INTO users (email) VALUES ($1) ON CONFLICT (email) DO NOTHING;`
	// _, err = db.Exec(addEmail, email)
	// if err != nil {
	// 	panic(err)
	// }

	// var userId string
	// getUser := `SELECT user_id FROM users WHERE email = $1`
	// err = db.QueryRow(getUser, email).Scan(&userId)
	// if err != nil {
	// 	panic(err)
	// }

	return db, ""
}

func (fetcher *BlogFetcher) GetEmail() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key, ok := os.LookupEnv("NAMESPACE_DOMAIN")
	if !ok {
		log.Fatal("Error loading env variables (namespace domain)")
	}

	_, token := Middleware()

	email := token[key].(string)

	return email
}
