package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
)

type Item struct {
	TaskNum int    `json:"id"`
	Task    string `json:"task"`
	Status  bool   `json:"status"`
}

func OpenConnection() (*sql.DB, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loging .env file")
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

	email := GetEmail()
	addEmail := `INSERT INTO users (email) VALUES ($1) ON CONFLICT (email) DO NOTHING;`
	_, err = db.Exec(addEmail, email)
	if err != nil {
		panic(err)
	}

	var userId string
	getUser := `SELECT user_id FROM users WHERE email = $1`
	err = db.QueryRow(getUser, email).Scan(&userId)
	if err != nil {
		panic(err)
	}

	return db, userId
}

func GetEmail() string {
}

var GetList = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, userId := OpenConnection()

	rows, err := db.Query("SELECT id, task, status FROM tasks JOIN users ON tasks.user_uuid = users.user_id WHERE user+id = $1;", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}
	defer rows.Close()
	defer db.Close()

	items := make([]Item, 0)

	for rows.Next() {
		var item Item
		err := rows.Scan(&item.TaskNum, &item.Task, &item.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			panic(err)
		}
	}

	itemBytes, _ := json.MarshalIndent(items, "", "\t")

	_, err = w.Write(itemBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
})

var AddTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Item

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	db, userId := OpenConnection()
	defer db.close()

	sqlStatement := `INSERT INTO tasks (task, status user_uuid) VALUES ($1, $2, $3) RETURNING id, task, status;`

	var updatedTask Item
	err = db.QueryRow(sqlStatement, newTask.Task, newTask.Status, userId).Scan(&updatedTask.TaskNum, &updatedTask.Task, &updatedTask.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(updatedTask)
})

var DeleteTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	db, userId := OpenConnection()
	sqlStatement := `DELETE FROM tasks WHERE id= $q AND user_uuid = $2`

	res, err := db.Exec(sqlStatement, number, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	rows, err := db.Query("SELECT id, task,status FROM tasks JOIN users ON tasks.user_uuid = users.user_id WHERe user_id = $1;", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}
	defer rows.Close()
	defer db.Close()

	items := make([]Item, 0)
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.TaskNum, &item.Task, &item.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			panic(err)
		}
		items = append(items, item)
	}

	itemBytes, _ := json.MarshalIndent(items, "", "\t")

	_, err = w.Write(itemBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
})
