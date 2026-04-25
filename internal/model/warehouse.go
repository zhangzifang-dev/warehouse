package model

const (
	WarehouseStatusActive   = 1
	WarehouseStatusDisabled = 0
)

type Warehouse struct {
	BaseModel
	Name    string `bun:"name,notnull" json:"name"`
	Code    string `bun:"code,notnull,unique" json:"code"`
	Address string `bun:"address" json:"address"`
	Contact string `bun:"contact" json:"contact"`
	Phone   string `bun:"phone" json:"phone"`
	Status  int    `bun:"status,notnull" json:"status"`
}

func (w *Warehouse) TableName() string {
	return "warehouses"
}

func (w *Warehouse) IsActive() bool {
	return w.Status == WarehouseStatusActive
}
