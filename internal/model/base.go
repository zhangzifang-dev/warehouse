package model

import (
	"context"
	"time"
)

type BaseModel struct {
	ID        int64      `bun:"id,pk,autoincrement" json:"id"`
	CreatedAt time.Time  `bun:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at" json:"updated_at"`
	CreatedBy int64      `bun:"created_by" json:"created_by"`
	UpdatedBy int64      `bun:"updated_by" json:"updated_by"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete" json:"deleted_at,omitempty"`
}

func (m *BaseModel) BeforeCreate(ctx context.Context) error {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

func (m *BaseModel) BeforeUpdate(ctx context.Context) error {
	m.UpdatedAt = time.Now()
	return nil
}

func (m *BaseModel) IsSoftDeleted() bool {
	return m.DeletedAt != nil
}
