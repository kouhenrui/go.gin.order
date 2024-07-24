package repository

import (
	"errors"
	"go.gin.order/pkg/msg"
	"go.gin.order/src/config/database"
	"go.gin.order/src/internal/pojo"
	"gorm.io/gorm"
)

type Account struct {
	db      *gorm.DB
	account pojo.Account
	role    pojo.Role
}

func NewAccount() *Account {
	return &Account{db: database.MysqlClient}
}

func (a *Account) FindAccountByName(username string) (*AccountRes, error) {
	var acc AccountRes
	err := a.db.Table("account as a").Joins("left join role as r on r.id =a.role").
		Select("a.id as id,a.name as name,a.phone as phone,a.email as email,"+
			"a.password as password,a.salt as salt,a.role as role,a.status as status,"+
			"a.type as type,a.access_token as access_token,r.name as role_name").Where("a.name = ?", username).Scan(&acc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(msg.AccountNotFoundError)
		}
		return nil, err
	}
	return &acc, nil
}
func (a *Account) AddAccount(acc pojo.Account) error {
	err := a.db.Create(&acc).Error
	if err != nil {
		return err
	}
	return nil
}
func (a *Account) FindEmail(email string) (*AccountRes, error) {
	var acc AccountRes
	err := a.db.Table("account as a").
		Joins("left join role as r on r.id =a.role").
		Joins("left join account_profile as ap on ap.account_id=a.id").
		Select("a.id as id,a.name as name,a.phone as phone,a.email as email,"+
			"a.password as password,a.salt as salt,a.role as role,a.status as status,"+
			"a.type as type,a.access_token as access_token,r.name as role_name,ap.gender as gender").Where("a.email = ?", email).Scan(&acc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(msg.AccountNotFoundError)
		}
		return nil, err
	}
	return &acc, nil
}
func (a *Account) SaveAccessToken(id uint, accessToken string) error {
	a.account.ID = id
	a.account.AccessToken = accessToken
	return a.db.Updates(&a.account).Error

}
func (a *Account) InfoById(id uint) (*InfoRes, error) {
	var info InfoRes
	err := a.db.
		Table("account as a").
		Joins("left join role as r on r.id =a.role").
		Joins("left join account_profile as ap on ap.account_id=a.id").
		Joins("left join department as d on d.id=ap.department_id").
		Select("a.id as id,a.name as name,a.phone as phone,a.email as email,a.role as role,a.status as status,a.type as type,r.name as role_name,"+
			"ap.nick_name as nick_name,d.name as department_name,ap.gender as gender,a.created_at as created_at").Where("a.id = ?", id).Scan(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}
func (a *Account) AllAccount(page, limit int) (*[]AllAccountRes, error) {
	var all []AllAccountRes
	err := a.db.Table("account as a").Joins("left join role as r on r.id = a.role ").
		Select("a.id as id,a.name as name,a.email as email,a.phone as phone,a.role as role,a.type as type,a.status as status,r.name as role_name").
		Limit(limit).Offset(page).Scan(&all).Error
	return &all, err
}

func (a *Account) RoleList(limit, page int) (*[]RoleRes, error) {
	var role []RoleRes
	err := a.db.Find(&a.role).Limit(limit).Offset(page).Scan(&role).Error
	return &role, err
}
func (a *Account) AuthorzedAccount(limit, page int) (*[]AllAccountRes, int64, error) { //
	var auth []AllAccountRes
	var total int64
	err := a.db.
		Table("account as a").
		Joins("left join account_profile ap on ap.account_id=a.id").
		Joins("left join role  r on r.id = a.role ").
		Where("a.role != ?", "").Count(&total).Error
	err = a.db.Table("account as a").
		Joins("left join account_profile ap on ap.account_id=a.id").
		Joins("left join role  r on r.id = a.role ").
		Select("a.id as id,a.name as name,a.email as email,a.phone as phone,a.role as role,a.type as type,a.status as status,r.name as role_name,ap.nick_name as nick_name").
		Where("a.role != ?", "").Limit(limit).Offset(page).Scan(&auth).Error //.Count(&total).Error
	return &auth, total, err
}
func (a *Account) UnAuthorzedAccount(limit, page int) (*[]AllAccountRes, int64, error) {
	var auth []AllAccountRes
	var total int64
	err := a.db.
		Table("account as a").
		Joins("left join account_profile ap on ap.account_id=a.id").
		Joins("left join role  r on r.id = a.role ").
		Where("a.role != ?", "").Count(&total).Error
	err = a.db.Table("account as a").
		Joins("left join account_profile ap on ap.account_id=a.id").
		Joins("left join role  r on r.id = a.role ").
		Select("a.id as id,a.name as name,a.email as email,a.phone as phone,a.role as role,a.type as type,a.status as status,r.name as role_name,ap.nick_name as nick_name").
		Where("a.role = ?", "").Limit(limit).Offset(page).Scan(&auth).Error //.Count(&total).Error
	return &auth, total, err
}
