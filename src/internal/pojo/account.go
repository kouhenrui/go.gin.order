package pojo

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	gorm.Model
	Name        string `gorm:"unique;not null" json:"name,omitempty"`
	Phone       string `gorm:"type:varchar(255)" json:"phone,omitempty"`
	Email       string `gorm:"unique;not null" validate:"required,email" json:"email,omitempty"`
	Password    string `gorm:"not null" json:"password,omitempty"`
	Salt        string `gorm:"not null" json:"salt,omitempty"`
	Role        int32  `gorm:"not null" json:"role,omitempty"`
	Status      int    `gorm:"type:tinyint;default:1" json:"status,omitempty"`
	Type        string `gorm:"not null" json:"type,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
}
type AccountProfile struct {
	gorm.Model
	AccountID    uint      `gorm:"not null;comment:'账号ID，关联account表'" json:"account_id,omitempty"`
	NickName     string    `gorm:"type:varchar(255);comment:'昵称'" json:"nick_name,omitempty"`
	DateOfBirth  time.Time `gorm:"type:date;comment:'出生日期'" json:"date_of_birth,omitempty"`
	Age          int       `gorm:"type:int;comment:'年龄'" json:"age,omitempty"`
	Gender       string    `gorm:"type:varchar(6);default:'other';comment:'性别（male, female, other）'" json:"gender,omitempty"`
	Identity     string    `gorm:"type:varchar(255);comment:'身份'" json:"identity,omitempty"`
	Address      string    `gorm:"type:varchar(255);comment:'地址'" json:"address,omitempty"`
	DepartmentID string    `gorm:"type:varchar(100);comment:'部门'" json:"department_id,omitempty"`
}

type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null;comment:名称" json:"name,omitempty"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态（1: 正常, 0: 停用）" json:"status,omitempty"`
	Description string `gorm:"type:varchar(255);not null;comment:描述" json:"description,omitempty"`
}
type Department struct {
	gorm.Model
	ParentID    uint   `gorm:"type:bigint;comment:'父级ID'" json:"parent_id,omitempty"`
	AncestorIDs string `gorm:"type:varchar(255);comment:'祖级ID列表'" json:"ancestor_ids,omitempty"`
	Status      bool   `gorm:"type:bool;default:true;comment:'部门状态（true, false）'" json:"status,omitempty"`
	Name        string `gorm:"type:varchar(100);not null;unique;comment:'部门名称'" json:"name,omitempty"`
	Order       int    `gorm:"type:int;comment:'显示顺序'" json:"order,omitempty"`
	Leader      string `gorm:"type:varchar(100);comment:'负责人'" json:"leader,omitempty"`
	Phone       string `gorm:"type:varchar(20);comment:'联系电话'" json:"phone,omitempty"`
	Email       string `gorm:"type:varchar(100);comment:'邮箱'" json:"email,omitempty"`
}
