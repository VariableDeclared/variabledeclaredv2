package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
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

type BlogPost struct {
	Content string `json:"content"`
	Title   string `json:"title"`
	Date    string `json:"date"`
}
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

var fetcher = BlogFetcher{}

var GetList = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, userId := fetcher.OpenConnection()

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

	db, userId := fetcher.OpenConnection()
	defer db.Close()

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

	db, userId := fetcher.OpenConnection()
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

var EditTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	number, err := strconv.Atoi((vars["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	sqlStatement := `UPDATE tasks SET task = $2 WHERE id = $1 AND user_uuid = $1 RETURNING id, task, status;`

	var newTask Item

	err = json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	db, userId := fetcher.OpenConnection()
	defer db.Close()

	var updatedTask Item
	err = db.QueryRow(sqlStatement, number, newTask.Task, userId).Scan(&updatedTask.TaskNum, &updatedTask.Task, &updatedTask.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(updatedTask)
})

var DoneTask = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	var currStatus bool

	var updatedTask Item

	sqlStatement1 := `SELECT status FROM tasks WHERE id = $1 AND user_uuid = $2;`
	sqlStatement2 := `UPDATE tasks set status = $2 WHERE id = $1 AND user_uuid = $3 RETURNING id, task, status;`

	db, userId := fetcher.OpenConnection()
	defer db.Close()

	err = db.QueryRow(sqlStatement1, number, userId).Scan(&currStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	err = db.QueryRow(sqlStatement2, number, !currStatus, userId).Scan(&updatedTask.TaskNum, &updatedTask.Task, &updatedTask.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(updatedTask)
})

var GetBlogPosts = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var queryArticles = "SELECT title, content, date FROM posts LIMIT 10;"
	var posts []BlogPost

	db, _ := fetcher.OpenConnection()
	defer db.Close()
	rows, err := db.Query(queryArticles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var post BlogPost
		if err := rows.Scan(&post.Title, &post.Content, &post.Date); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			panic(err)
		}
		posts = append(posts, post)

	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(posts)
})
