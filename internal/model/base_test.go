package model

import (
	"testing"
	"time"
)

func TestBaseModelFields(t *testing.T) {
	now := time.Now()
	model := BaseModel{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: 100,
		UpdatedBy: 100,
	}

	if model.ID != 1 {
		t.Errorf("ID = %d, want 1", model.ID)
	}
	if !model.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", model.CreatedAt, now)
	}
	if !model.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", model.UpdatedAt, now)
	}
	if model.CreatedBy != 100 {
		t.Errorf("CreatedBy = %d, want 100", model.CreatedBy)
	}
	if model.UpdatedBy != 100 {
		t.Errorf("UpdatedBy = %d, want 100", model.UpdatedBy)
	}
}

func TestBaseModelSoftDelete(t *testing.T) {
	now := time.Now()
	model := BaseModel{
		ID:        1,
		DeletedAt: &now,
	}

	if model.DeletedAt == nil {
		t.Error("DeletedAt should not be nil")
	}
}

func TestBaseModelZeroValues(t *testing.T) {
	model := BaseModel{}

	if model.ID != 0 {
		t.Errorf("ID = %d, want 0", model.ID)
	}
	if !model.CreatedAt.IsZero() {
		t.Errorf("CreatedAt should be zero")
	}
	if !model.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should be zero")
	}
	if model.CreatedBy != 0 {
		t.Errorf("CreatedBy = %d, want 0", model.CreatedBy)
	}
	if model.UpdatedBy != 0 {
		t.Errorf("UpdatedBy = %d, want 0", model.UpdatedBy)
	}
	if model.DeletedAt != nil {
		t.Errorf("DeletedAt should be nil")
	}
}

func TestBaseModelBeforeCreate(t *testing.T) {
	model := BaseModel{}

	err := model.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate returned error: %v", err)
	}

	if model.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set by BeforeCreate")
	}
	if model.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set by BeforeCreate")
	}
}

func TestBaseModelBeforeUpdate(t *testing.T) {
	model := BaseModel{}

	err := model.BeforeUpdate(nil)
	if err != nil {
		t.Errorf("BeforeUpdate returned error: %v", err)
	}

	if model.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set by BeforeUpdate")
	}
}

func TestBaseModelIsSoftDeleted(t *testing.T) {
	model := BaseModel{}
	if model.IsSoftDeleted() {
		t.Error("Empty model should not be soft deleted")
	}

	now := time.Now()
	model.DeletedAt = &now
	if !model.IsSoftDeleted() {
		t.Error("Model with DeletedAt should be soft deleted")
	}
}
