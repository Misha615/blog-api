package usecases

import (
	"blog-api/internal/entities"
	"blog-api/internal/repositories"
	"fmt"
)

type ArticleUseCase struct {
	repo *repositories.ArticleRepository
}

func NewArticleUseCase(repo *repositories.ArticleRepository) *ArticleUseCase {
	return &ArticleUseCase{repo: repo}
}

func (uc *ArticleUseCase) GetArticles() ([]entities.Article, error) {
	return uc.repo.GetAll()
}

func (uc *ArticleUseCase) GetArticle(id int64) (*entities.Article, error) {
	return uc.repo.GetByID(id)
}

func (uc *ArticleUseCase) CreateArticle(title, content string) (int64, error) {
	return uc.repo.Create(title, content)
}

func (uc *ArticleUseCase) UpdateArticle(id int64, updateArticle *entities.Article) error {
	if updateArticle == nil {
		return fmt.Errorf("article cannot be nil")
	}
	return uc.repo.Update(id, updateArticle)
}

func (uc *ArticleUseCase) DeleteArticle(id int64) error {
	return uc.repo.DeleteByID(id)
}
