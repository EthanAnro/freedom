//Package po generated by 'freedom new-po'
package po

import (
	"gorm.io/gorm"
	"time"
)

// Goods .
type Goods struct {
	changes map[string]interface{}
	ID      int       `gorm:"primaryKey;column:id"`
	Name    string    `gorm:"column:name"`    // 商品名称
	Price   int       `gorm:"column:price"`   // 价格
	Stock   int       `gorm:"column:stock"`   // 库存
	Version int       `gorm:"column:version"` // 乐观锁版本号
	Created time.Time `gorm:"column:created"`
	Updated time.Time `gorm:"column:updated"`
}

// TableName .
func (obj *Goods) TableName() string {
	return "goods"
}

// Location .
func (obj *Goods) Location() map[string]interface{} {
	return map[string]interface{}{"id": obj.ID}
}

// GetChanges .
func (obj *Goods) GetChanges() map[string]interface{} {
	if obj.changes == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range obj.changes {
		result[k] = v
	}
	obj.changes = nil
	return result
}

// Update .
func (obj *Goods) Update(name string, value interface{}) {
	if obj.changes == nil {
		obj.changes = make(map[string]interface{})
	}
	obj.changes[name] = value
}

// SetName .
func (obj *Goods) SetName(name string) {
	obj.Name = name
	obj.Update("name", name)
}

// SetPrice .
func (obj *Goods) SetPrice(price int) {
	obj.Price = price
	obj.Update("price", price)
}

// SetStock .
func (obj *Goods) SetStock(stock int) {
	obj.Stock = stock
	obj.Update("stock", stock)
}

// SetVersion .
func (obj *Goods) SetVersion(version int) {
	obj.Version = version
	obj.Update("version", version)
}

// SetCreated .
func (obj *Goods) SetCreated(created time.Time) {
	obj.Created = created
	obj.Update("created", created)
}

// SetUpdated .
func (obj *Goods) SetUpdated(updated time.Time) {
	obj.Updated = updated
	obj.Update("updated", updated)
}

// AddPrice .
func (obj *Goods) AddPrice(price int) {
	obj.Price += price
	obj.Update("price", gorm.Expr("price + ?", price))
}

// AddStock .
func (obj *Goods) AddStock(stock int) {
	obj.Stock += stock
	obj.Update("stock", gorm.Expr("stock + ?", stock))
}

// AddVersion .
func (obj *Goods) AddVersion(version int) {
	obj.Version += version
	obj.Update("version", gorm.Expr("version + ?", version))
}
