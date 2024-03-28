package service

type Service struct {
}

type ServiceInterface interface {
	CreateWordPress() (string, int, error)
}

func NewService() ServiceInterface {
	svc := &Service{}
	return svc
}
