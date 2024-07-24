package service

import (
	"go.gin.order/src/internal/repository"
	"strconv"
)

type SystemService struct {
	acc *repository.Account
}

func NewSystemService() *SystemService {
	return &SystemService{
		acc: repository.NewAccount(),
	}
}
func (s *SystemService) RoleList(l, p string) (*[]repository.RoleRes, error) {
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	return s.acc.RoleList(limit, page)
}
func (s *SystemService) AuthorizedAccount(l, p string) (*[]repository.AllAccountRes, int64, error) {
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	//s.acc.AuthorzedAccount(limit, page)
	return s.acc.AuthorzedAccount(limit, page)
}
func (s *SystemService) UnAuthorizedAccount(l, p string) (*[]repository.AllAccountRes, int64, error) {
	page, _ := strconv.Atoi(p)
	limit, _ := strconv.Atoi(l)
	return s.acc.UnAuthorzedAccount(limit, page)
}
