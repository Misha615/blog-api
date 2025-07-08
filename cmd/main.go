package main

import (
	"blog-api/internal/delivery/http"
	"blog-api/internal/repositories"
	"blog-api/internal/usecases"
	"blog-api/package/database"
	"log"
	nethttp "net/http"
	"time"

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

	r.Use(loggingMiddleware)
	r.Use(panicRecoveryMiddleware)

	r.HandleFunc("/articles", handler.GetArticles).Methods("GET")
	r.HandleFunc("/articles/{id}", handler.GetArticle).Methods("GET")
	r.HandleFunc("/articles", handler.CreateArticle).Methods("POST")
	r.HandleFunc("/articles/{id}", handler.UpdateArticle).Methods("PUT")
	r.HandleFunc("/articles/{id}", handler.DeleteArticle).Methods("DELETE")

	server := &nethttp.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server started on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func loggingMiddleware(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func panicRecoveryMiddleware(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				http.RespondWithError(w, nethttp.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
