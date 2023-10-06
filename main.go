package main

import (
	"blogapi/database"
	"blogapi/models"
	"encoding/json"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var post models.Post
	err = db.QueryRow("SELECT id, title, content FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var posts []models.Post
	rows, err := db.Query("SELECT id, title, content FROM posts")
	if err == sql.ErrNoRows {
		json.NewEncoder(w).Encode([]models.Post{})
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		posts = append(posts, post)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO posts (title, content) VALUES ($1, $2)", post.Title, post.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.String())
        next.ServeHTTP(w, r)
    })
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/posts/{id}", getPost).Methods("GET")
	r.HandleFunc("/posts", getAllPosts).Methods("GET", "POST")
	http.Handle("/", r)

	r.Use(loggingMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}

