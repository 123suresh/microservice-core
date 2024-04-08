package repository

import (
	"github.com/core-api/internal/model"
)

type WordpressRepoInterface interface {
	CreateWordPress(data *model.Wordpress) (*model.Wordpress, error)
	CountUser() (int64, error)
	GetWordPress() ([]model.Wordpress, error)
	DeleteWordPress(req *model.DelWordpress) error
}

func (repo *Repo) CreateWordPress(data *model.Wordpress) (*model.Wordpress, error) {
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

func (repo *Repo) GetWordPress() ([]model.Wordpress, error) {
	details := []model.Wordpress{}
	err := repo.db.Model(&model.Wordpress{}).Find(&details).Error
	if err != nil {
		return nil, err
	}
	return details, nil
}

func (repo *Repo) DeleteWordPress(req *model.DelWordpress) error {
	data := &model.Wordpress{}
	err := repo.db.Model(&model.Wordpress{}).Where("namespace=?", req.Namespace).Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}
