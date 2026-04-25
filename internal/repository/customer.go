package repository

import (
	"context"
	"time"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type CustomerRepository struct {
	db *bun.DB
}

func NewCustomerRepository(db *bun.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *model.Customer) error {
	_, err := r.db.NewInsert().Model(customer).Exec(ctx)
	return err
}

func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	customer := new(model.Customer)
	err := r.db.NewSelect().
		Model(customer).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *CustomerRepository) GetByCode(ctx context.Context, code string) (*model.Customer, error) {
	customer := new(model.Customer)
	err := r.db.NewSelect().
		Model(customer).
		Where("code = ?", code).
		Where("deleted_at = ?", timeZero).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *CustomerRepository) List(ctx context.Context, page, pageSize int, keyword string) ([]model.Customer, int, error) {
	var customers []model.Customer
	query := r.db.NewSelect().
		Model(&customers).
		Where("deleted_at = ?", timeZero)

	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	total, err := query.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}
	return customers, total, nil
}

func (r *CustomerRepository) Update(ctx context.Context, customer *model.Customer) error {
	_, err := r.db.NewUpdate().
		Model(customer).
		WherePK().
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}

func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Customer)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at = ?", timeZero).
		Exec(ctx)
	return err
}
