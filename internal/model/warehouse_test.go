package model

import "testing"

func TestWarehouse_TableName(t *testing.T) {
	w := Warehouse{}
	if w.TableName() != "warehouses" {
		t.Errorf("Warehouse.TableName() = %s, want warehouses", w.TableName())
	}
}

func TestWarehouse_IsActive(t *testing.T) {
	w := Warehouse{Status: WarehouseStatusActive}
	if !w.IsActive() {
		t.Error("IsActive() = false for active warehouse")
	}

	w.Status = WarehouseStatusDisabled
	if w.IsActive() {
		t.Error("IsActive() = true for disabled warehouse")
	}
}
