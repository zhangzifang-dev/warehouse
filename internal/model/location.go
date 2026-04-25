package model

const (
	LocationStatusActive   = 1
	LocationStatusInactive = 0
)

type Location struct {
	BaseModel
	WarehouseID int64  `bun:"warehouse_id,notnull" json:"warehouse_id"`
	Zone        string `bun:"zone,notnull" json:"zone"`
	Shelf       string `bun:"shelf,notnull" json:"shelf"`
	Level       string `bun:"level,notnull" json:"level"`
	Position    string `bun:"position,notnull" json:"position"`
	Code        string `bun:"code,notnull" json:"code"`
	Status      int    `bun:"status,notnull" json:"status"`
}

func (l *Location) TableName() string {
	return "locations"
}

func (l *Location) IsActive() bool {
	return l.Status == LocationStatusActive
}

func (l *Location) GenerateCode() string {
	return l.Zone + "-" + l.Shelf + "-" + l.Level + "-" + l.Position
}
