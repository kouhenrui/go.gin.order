package repository

import "time"

type PaginationRes struct {
	data  interface{}
	total int64
}

type AccountRes struct {
	Id          uint   `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	Salt        string `json:"salt,omitempty"`
	Role        int32  `json:"role,omitempty"`
	RoleName    string `json:"role_name,omitempty"`
	Status      int    `json:"status,omitempty"`
	Type        string `json:"type,omitempty"`
	AccessToken string `json:"access_token,omitempty"`
	Gender      string `json:"gender"`
}
type AddAccount struct {
	Name     string
	Password string
	Salt     string
	Phone    string
	Status   int
	Email    string
	Type     string
	Role     int32
}
type InfoRes struct {
	Id             uint      `json:"id,omitempty"`
	NickName       string    `json:"nick_name"`
	Name           string    `json:"name,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	Email          string    `json:"email,omitempty"`
	Role           int32     `json:"role,omitempty"`
	RoleName       string    `json:"role_name,omitempty"`
	Type           string    `json:"type,omitempty"`
	DepartmentName string    `json:"department_name"`
	Gender         string    `json:"gender"`
	CreatedAt      time.Time `json:"created_at"`
}
type ProductRes struct {
	MerchantID     uint     `json:"merchant_id,omitempty"`
	MerchantName   string   `json:"merchant_name"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	Price          float64  `json:"price,omitempty"`
	Stock          uint     `json:"stock,omitempty"`
	ImageURL       []string `json:"image_url,omitempty"`
	OnSale         bool     `json:"on_sale,omitempty"`
	SaleCount      uint     `json:"sale_count,omitempty"`
	Views          uint     `json:"views,omitempty"`
	Category       string   `json:"category,omitempty"`
	Specifications []string `json:"specifications,omitempty"`
}
type MerchantListRes struct {
	Name        string
	Category    string
	Phone       string
	Address     string
	City        string
	State       string
	Status      int
	Description string
}
type AllAccountRes struct {
	Id       uint   `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	NickName string `json:"nick_name"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     int32  `json:"role,omitempty"`
	RoleName string `json:"role_name,omitempty"`
	Status   int    `json:"status,omitempty"`
	Type     string `json:"type,omitempty"`
}
type RoleRes struct {
	Id          uint
	Name        string
	Status      int
	Description string
}
