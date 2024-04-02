package model

import "gorm.io/gorm"

type Wordpress struct {
	gorm.Model
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type WordPressRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type WordPressResponse struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Message   string `json:"message"`
}

func RequestWordPress(WordPressreq *WordPressRequest) *Wordpress {
	return &Wordpress{
		Name:      WordPressreq.Name,
		Namespace: WordPressreq.Namespace,
	}
}

func (w *Wordpress) WordPressResponse() *WordPressResponse {
	return &WordPressResponse{
		Name:      w.Name,
		Namespace: w.Namespace,
	}
}
