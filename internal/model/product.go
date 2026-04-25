package model

const (
	ProductStatusActive   = 1
	ProductStatusInactive = 0
)

type Product struct {
	BaseModel
	SKU           string  `bun:"sku,notnull,unique" json:"sku"`
	Name          string  `bun:"name,notnull" json:"name"`
	CategoryID    int64   `bun:"category_id" json:"category_id"`
	Specification string  `bun:"specification" json:"specification"`
	Unit          string  `bun:"unit" json:"unit"`
	Barcode       string  `bun:"barcode" json:"barcode"`
	Price         float64 `bun:"price" json:"price"`
	Description   string  `bun:"description" json:"description"`
	Status        int     `bun:"status,notnull" json:"status"`
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) IsActive() bool {
	return p.Status == ProductStatusActive
}
