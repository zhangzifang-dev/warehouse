package model

const (
	SupplierStatusActive   = 1
	SupplierStatusInactive = 0
)

type Supplier struct {
	BaseModel
	Name    string `bun:"name,notnull" json:"name"`
	Code    string `bun:"code,unique" json:"code"`
	Contact string `bun:"contact" json:"contact"`
	Phone   string `bun:"phone" json:"phone"`
	Email   string `bun:"email" json:"email"`
	Address string `bun:"address" json:"address"`
	Status  int    `bun:"status,notnull" json:"status"`
}

func (s *Supplier) TableName() string {
	return "suppliers"
}

func (s *Supplier) IsActive() bool {
	return s.Status == SupplierStatusActive
}
