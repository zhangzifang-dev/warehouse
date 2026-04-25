package model

import "testing"

func TestInboundOrder_TableName(t *testing.T) {
	order := &InboundOrder{}
	if order.TableName() != "inbound_orders" {
		t.Errorf("expected table name 'inbound_orders', got '%s'", order.TableName())
	}
}

func TestInboundItem_TableName(t *testing.T) {
	item := &InboundItem{}
	if item.TableName() != "inbound_items" {
		t.Errorf("expected table name 'inbound_items', got '%s'", item.TableName())
	}
}
