package repository

import (
	"github.com/core-api/internal/model"
	"github.com/sirupsen/logrus"
)

type WordpressRepoInterface interface {
	CreateWordPress(data *model.Wordpress) (*model.Wordpress, error)
	CountUser() (int64, error)
}

func (repo *Repo) CreateWordPress(data *model.Wordpress) (*model.Wordpress, error) {
	logrus.Info("data => ", data)
	err := repo.db.Model(&model.Wordpress{}).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *Repo) CountUser() (int64, error) {
	var count int64
	err := repo.db.Model(&model.Wordpress{}).Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
}
