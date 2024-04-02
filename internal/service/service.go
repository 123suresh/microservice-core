package service

import (
	"github.com/core-api/internal/model"
	"github.com/core-api/internal/repository"
)

type Service struct {
	repo repository.RepoInterface
}

type ServiceInterface interface {
	CreateWordPress(req *model.WordPressRequest) (*model.WordPressResponse, int, error)
}

func NewService(repo repository.RepoInterface) ServiceInterface {
	svc := &Service{}
	svc.repo = repo
	// svc.repo = repository.NewRepo()
	return svc
}
