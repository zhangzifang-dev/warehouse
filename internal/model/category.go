package model

const (
	CategoryStatusActive   = 1
	CategoryStatusInactive = 0
)

type Category struct {
	BaseModel
	Name      string `bun:"name,notnull" json:"name"`
	ParentID  int64  `bun:"parent_id" json:"parent_id"`
	SortOrder int    `bun:"sort_order" json:"sort_order"`
	Status    int    `bun:"status,notnull" json:"status"`
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) IsActive() bool {
	return c.Status == CategoryStatusActive
}
