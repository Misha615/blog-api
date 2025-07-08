package main

import (
	"blog-api/internal/delivery/http"
	"blog-api/internal/repositories"
	"blog-api/internal/usecases"
	"blog-api/package/database"
	"log"
	nethttp "net/http"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repositories.NewArticleRepo(db)
	uc := usecases.NewArticleUseCase(repo)
	handler := http.NewArticleHandler(uc)

	r := mux.NewRouter()
	r.HandleFunc("/articles", handler.GetArticles).Methods("GET")
	r.HandleFunc("/articles/{id}", handler.GetArticle).Methods("GET")
	r.HandleFunc("/articles", handler.CreateArticle).Methods("POST")
	r.HandleFunc("/articles/{id}", handler.UpdateArticle).Methods("PUT")
	r.HandleFunc("/articles/{id}", handler.DeleteArticle).Methods("DELETE")

	log.Println("Server started on :8080")
	log.Fatal(nethttp.ListenAndServe(":8080", r))
}
