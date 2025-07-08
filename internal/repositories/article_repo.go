package repositories

import (
	"blog-api/internal/entities"
	"database/sql"
	"fmt"
	"log"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepo(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) GetAll() ([]entities.Article, error) {
	row, err := r.db.Query("SELECT id, title, content, created_at FROM articles")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %w", err)
	}
	defer func() {
		if err := row.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()

	var articles []entities.Article
	for row.Next() {
		var a entities.Article
		if err := row.Scan(&a.ID, &a.Title, &a.Content, &a.CreatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	return articles, nil
}

func (r *ArticleRepository) GetByID(id int64) (*entities.Article, error) {
	row := r.db.QueryRow("SELECT id, title, content, created_at FROM articles WHERE id = ?", id)

	var article entities.Article
	if err := row.Scan(&article.ID, &article.Title, &article.Content, &article.CreatedAt); err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *ArticleRepository) Create(title, content string) (int64, error) {
	if title == "" || content == "" {
		return -1, fmt.Errorf("title and content cannot be empty")
	}

	row, err := r.db.Exec("INSERT INTO articles (title, content) VALUES(?, ?)",
		title,
		content,
	)
	if err != nil {
		return -1, fmt.Errorf("failed to create article: %w", err)
	}

	return row.LastInsertId()
}

func (r *ArticleRepository) Update(id int64, updateArticle *entities.Article) error {
	_, err := r.db.Exec(
		`UPDATE articles
		SET title = ?, content = ?, created_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		updateArticle.Title,
		updateArticle.Content,
		id,
	)
	return err
}

func (r *ArticleRepository) DeleteByID(id int64) error {
	_, err := r.db.Exec("DELETE FROM articles WHERE id = ?", id)
	/*if err != nil {
		log.Fatal(err)
	}*/
	return err
}
