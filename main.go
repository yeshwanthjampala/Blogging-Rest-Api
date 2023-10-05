package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func main() {
	// Connect to the MySQL database
	// db, err = gorm.Open("mysql", "username:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	db, err = gorm.Open("mysql", "root:Yeshwanth@1234@tcp(localhost:3306)/blogging?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println("Unable to connect to database")
	} else {
		log.Println("Connection Successful")
	}

	// AutoMigrate ensures the table "posts" exists and has the appropriate schema
	db.AutoMigrate(&Post{})

	r := mux.NewRouter()
	r.HandleFunc("/posts", GetPosts).Methods("GET")
	r.HandleFunc("/posts/{id}", GetPost).Methods("GET")
	r.HandleFunc("/posts", CreatePost).Methods("POST")
	r.HandleFunc("/posts/{id}", UpdatePost).Methods("PUT")
	r.HandleFunc("/posts/{id}", DeletePost).Methods("DELETE")

	http.Handle("/", r)

	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	db.Find(&posts)
	respondWithJSON(w, http.StatusOK, posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var post Post
	db.First(&post, params["id"])
	if post.ID == 0 {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}
	respondWithJSON(w, http.StatusOK, post)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&post); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	db.Create(&post)
	respondWithJSON(w, http.StatusCreated, post)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var post Post
	db.First(&post, params["id"])
	if post.ID == 0 {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	var updatedPost Post
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updatedPost); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	db.Model(&post).Updates(updatedPost)
	respondWithJSON(w, http.StatusOK, post)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var post Post
	db.First(&post, params["id"])
	if post.ID == 0 {
		respondWithError(w, http.StatusNotFound, "Post not found")
		return
	}

	db.Delete(&post)
	respondWithJSON(w, http.StatusNoContent, nil)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
