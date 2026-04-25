package model

type Inventory struct {
	BaseModel
	WarehouseID int64   `bun:"warehouse_id,notnull" json:"warehouse_id"`
	ProductID   int64   `bun:"product_id,notnull" json:"product_id"`
	LocationID  int64   `bun:"location_id" json:"location_id"`
	Quantity    float64 `bun:"quantity,notnull" json:"quantity"`
	BatchNo     string  `bun:"batch_no" json:"batch_no"`
}

func (i *Inventory) TableName() string {
	return "inventories"
}
