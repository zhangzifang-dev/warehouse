package model

import "testing"

func TestInventory_TableName(t *testing.T) {
	inv := &Inventory{}
	if inv.TableName() != "inventories" {
		t.Errorf("expected table name 'inventories', got '%s'", inv.TableName())
	}
}
