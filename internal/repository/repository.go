package repository

import (
	"github.com/core-api/internal/database"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo() *Repo {
	return &Repo{
		db: database.InitializeDB(),
	}
}

type RepoInterface interface {
	WordpressRepoInterface
}
