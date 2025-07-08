package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"blog-api/internal/entities"
	"blog-api/internal/usecases"

	"github.com/gorilla/mux"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("JSON encoding error: %v", err)
	}
}

type ArticleHandler struct {
	uc *usecases.ArticleUseCase
}

func NewArticleHandler(uc *usecases.ArticleUseCase) *ArticleHandler {
	return &ArticleHandler{uc: uc}
}

func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	articles, err := h.uc.GetArticles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		RespondWithError(w, http.StatusMethodNotAllowed, "Invalid article ID")
		return
	}
	article, err := h.uc.GetArticle(id)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Article not found")
		} else {
			RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, article)
}

func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var article entities.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if article.Title == " " || article.Content == "" {
		http.Error(w, `{"error": "Title and content are required"}`, http.StatusBadRequest)
		return
	}

	createdArticle, err := h.uc.CreateArticle(article.Title, article.Content)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdArticle)
}

func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var article entities.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	h.uc.UpdateArticle(id, &article)
	fmt.Fprintf(w, "Статья с id = %d обновлена", id)
}

func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.uc.DeleteArticle(id)

	fmt.Fprintf(w, "Article with id = %d deleted", id)
}
