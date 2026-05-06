package model

import "time"

type InboundOrderQueryFilter struct {
	OrderNo        string
	SupplierID     *int64
	WarehouseID    *int64
	QuantityMin    *float64
	QuantityMax    *float64
	CreatedAtStart *time.Time
	CreatedAtEnd   *time.Time
	Page           int
	PageSize       int
}

type InboundOrder struct {
	BaseModel
	OrderNo       string  `bun:"order_no,notnull,unique" json:"order_no"`
	SupplierID    int64   `bun:"supplier_id" json:"supplier_id"`
	WarehouseID   int64   `bun:"warehouse_id,notnull" json:"warehouse_id"`
	TotalQuantity float64 `bun:"total_quantity,notnull" json:"total_quantity"`
	Status        int     `bun:"status,notnull" json:"status"`
	Remark        string  `bun:"remark" json:"remark"`
}

func (o *InboundOrder) TableName() string {
	return "inbound_orders"
}

type InboundItem struct {
	BaseModel
	OrderID    int64   `bun:"order_id,notnull" json:"order_id"`
	ProductID  int64   `bun:"product_id,notnull" json:"product_id"`
	LocationID int64   `bun:"location_id" json:"location_id"`
	Quantity   float64 `bun:"quantity,notnull" json:"quantity"`
	BatchNo    string  `bun:"batch_no" json:"batch_no"`
}

func (i *InboundItem) TableName() string {
	return "inbound_items"
}

type OutboundOrder struct {
	BaseModel
	OrderNo       string  `bun:"order_no,notnull,unique" json:"order_no"`
	CustomerID    int64   `bun:"customer_id" json:"customer_id"`
	WarehouseID   int64   `bun:"warehouse_id,notnull" json:"warehouse_id"`
	TotalQuantity float64 `bun:"total_quantity,notnull" json:"total_quantity"`
	Status        int     `bun:"status,notnull" json:"status"`
	Remark        string  `bun:"remark" json:"remark"`
}

func (o *OutboundOrder) TableName() string {
	return "outbound_orders"
}

type OutboundItem struct {
	BaseModel
	OrderID    int64   `bun:"order_id,notnull" json:"order_id"`
	ProductID  int64   `bun:"product_id,notnull" json:"product_id"`
	LocationID int64   `bun:"location_id" json:"location_id"`
	Quantity   float64 `bun:"quantity,notnull" json:"quantity"`
	BatchNo    string  `bun:"batch_no" json:"batch_no"`
}

func (i *OutboundItem) TableName() string {
	return "outbound_items"
}

type StockTransfer struct {
	BaseModel
	OrderNo           string `bun:"order_no,notnull,unique" json:"order_no"`
	SourceWarehouseID int64  `bun:"source_warehouse_id,notnull" json:"source_warehouse_id"`
	TargetWarehouseID int64  `bun:"target_warehouse_id,notnull" json:"target_warehouse_id"`
	Status            int    `bun:"status,notnull" json:"status"`
}

func (s *StockTransfer) TableName() string {
	return "stock_transfers"
}

type StockTransferItem struct {
	BaseModel
	TransferID int64   `bun:"transfer_id,notnull" json:"transfer_id"`
	ProductID  int64   `bun:"product_id,notnull" json:"product_id"`
	LocationID int64   `bun:"location_id" json:"location_id"`
	Quantity   float64 `bun:"quantity,notnull" json:"quantity"`
	BatchNo    string  `bun:"batch_no" json:"batch_no"`
}

func (i *StockTransferItem) TableName() string {
	return "stock_transfer_items"
}
