package pojo

import (
	"gorm.io/gorm"
	"time"
)

type Merchant struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(255);not null;comment:商家名称"`
	Description string    `gorm:"type:text;comment:商家描述"`
	CategoryID  uint      `gorm:"not null;comment:商家类别ID"`
	Phone       string    `gorm:"type:varchar(20);comment:商家联系电话"`
	Email       string    `gorm:"type:varchar(255);comment:商家联系邮箱"`
	Address     string    `gorm:"type:varchar(255);comment:商家地址"`
	City        string    `gorm:"type:varchar(100);comment:所在城市"`
	State       string    `gorm:"type:varchar(100);comment:所在省份/州"`
	PostalCode  string    `gorm:"type:varchar(20);comment:邮政编码"`
	Latitude    float64   `gorm:"type:decimal(10,7);comment:地理位置纬度"`
	Longitude   float64   `gorm:"type:decimal(10,7);comment:地理位置经度"`
	Rating      float64   `gorm:"type:decimal(3,2);comment:商家评分"`
	OpeningTime string    `gorm:"type:varchar(255);comment:开始营业时间"`
	ClosingTime string    `gorm:"type:varchar(255);comment:结束营业时间"`
	Status      int       `gorm:"type:tinyint;default:1;comment:商家状态（1: 正常, 0: 停用）"`
	IsOpen      bool      `gorm:"type:boolean;comment:是否在营业"`
	Products    []Product `gorm:"foreignKey:MerchantID"` // 商品关联
}
type Category struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null;comment:类别名称"`
	Description string `gorm:"type:text;comment:类别描述"`
	Status      int    `gorm:"type:tinyint;default:1;comment:状态（1: 正常, 0: 停用）"`
}

// Product 表示商品表
type Product struct {
	gorm.Model
	MerchantID     uint     `gorm:"not null;index;comment:商家ID，外键"`
	Name           string   `gorm:"type:varchar(255);not null;comment:商品名称"`
	Description    string   `gorm:"type:text;comment:商品描述"`
	Price          float64  `gorm:"type:decimal(10,2);not null;comment:商品价格"`
	Stock          uint     `gorm:"not null;comment:商品库存"`
	ImageURL       []string `gorm:"type:text;comment:商品图片URL"`
	OnSale         bool     `gorm:"default:false;comment:商品是否在售"`
	SaleCount      uint     `gorm:"default:0;comment:商品销量"`
	Views          uint     `gorm:"default:0;comment:商品浏览次数"`
	CategoryID     uint     `gorm:"not null;comment:商品类别ID"`
	Specifications []string `gorm:"type:json;comment:商品规格"`
}

// 购物车表
type Cart struct {
	gorm.Model
	UserID    uint       `json:"user_id" comment:"用户ID"`
	CartItems []CartItem `gorm:"foreignKey:CartID" json:"cart_items" comment:"购物车项"`
}

// 购物车项表
type CartItem struct {
	gorm.Model
	CartID    uint    `json:"cart_id" comment:"购物车ID"`
	ProductID uint    `json:"product_id" comment:"商品ID"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product" comment:"商品"`
	Quantity  uint    `gorm:"not null" json:"quantity" comment:"商品数量"`
}

// 订单表
type Order struct {
	gorm.Model
	UserID     uint        `json:"user_id" comment:"用户ID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items" comment:"订单项"`
	TotalPrice float64     `gorm:"type:decimal(10,2);not null" json:"total_price" comment:"总价"`
	Status     string      `gorm:"type:varchar(50);not null" json:"status" comment:"订单状态"`
	OrderDate  time.Time   `gorm:"not null" json:"order_date" comment:"订单日期"`
}

// 订单项表
type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id" comment:"订单ID"`
	ProductID uint    `json:"product_id" comment:"商品ID"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product" comment:"商品"`
	Quantity  uint    `gorm:"not null" json:"quantity" comment:"商品数量"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price" comment:"商品价格"`
}

// 评价表
type Review struct {
	gorm.Model
	UserID    uint    `json:"user_id" comment:"用户ID"`
	ProductID uint    `json:"product_id" comment:"商品ID"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product" comment:"商品"`
	Rating    uint    `gorm:"type:int;not null" json:"rating" comment:"评分"`
	Comment   string  `gorm:"type:text" json:"comment" comment:"评价内容"`
}
