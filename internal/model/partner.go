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

const (
	CustomerStatusActive   = 1
	CustomerStatusInactive = 0
)

type Customer struct {
	BaseModel
	Name    string `bun:"name,notnull" json:"name"`
	Code    string `bun:"code,unique" json:"code"`
	Contact string `bun:"contact" json:"contact"`
	Phone   string `bun:"phone" json:"phone"`
	Email   string `bun:"email" json:"email"`
	Address string `bun:"address" json:"address"`
	Status  int    `bun:"status,notnull" json:"status"`
}

func (c *Customer) TableName() string {
	return "customers"
}

func (c *Customer) IsActive() bool {
	return c.Status == CustomerStatusActive
}
