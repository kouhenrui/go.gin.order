package repository

import (
	mysql "go.gin.order/src/config/database"
	"go.gin.order/src/internal/pojo"
	"gorm.io/gorm"
)

type Mertant struct {
	merchant pojo.Merchant
	category pojo.Category
	product  pojo.Product
	db       *gorm.DB
}

func NewMerchat() *Mertant {
	return &Mertant{db: mysql.MysqlClient}
}
func (m *Mertant) MerctantList(page, limit int, name string) (*[]MerchantListRes, error) { //([]dto.MerchantListRes, error)
	//var me []pojo.Merchant
	var result []MerchantListRes
	query := m.db.Model(&m.merchant).Select("merchant.Name as name," +
		"merchant.Phone as phone," +
		"merchant.Address as address," +
		"merchant.City as city," +
		"merchant.State as state," +
		"merchant.Status as status," +
		"merchant.Description as description," + "category.name as category").Joins("left join category on category.id = merchant.category_iD").Scan(&result)
	if len(name) > 0 {
		query.Where("name like ?", "%"+name+"%")
	}

	//log.Println(query, "000000000000000000000000")
	//var mertantres []dto.MerchantListRes
	//query := db.Preload("Category").Model(&m.merchant)
	//if len(name) > 0 {
	//	query.Where("name like ?", "%"+name+"%")
	//}
	if err := query.Limit(limit).Offset(page).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
func (m *Mertant) MerchantProducts(id uint) (*[]ProductRes, error) {
	var products []ProductRes
	err := m.db.Table("product as p").
		Joins("left join merchant m on m.id=p.merchant_id").
		Joins("left join category c on c.id = p.category_id").
		Select("p.name as name,p.description as description,p.price as price,p.stock as stock,"+
			"p.image_url as image_url,p.on_sale as on_sale,p.sale_count as sale_count,p.views as views,"+
			"p.specifications as specifications,c.name as category,m.name as merchant_name,m.id as merchant_id").Where("m.id = ?", id).
		Scan(&products).
		Error
	if err != nil {
		return nil, err
	}
	return &products, nil

}
