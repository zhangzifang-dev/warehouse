package repository

import (
	"time"
	"context"

	"github.com/uptrace/bun"
	"warehouse/internal/model"
)

type CustomerQueryFilter struct {
	Code     string
	Name     string
	Phone    string
	Status   *int
	Page     int
	PageSize int
}

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
		Where("deleted_at IS NULL").
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
		Where("deleted_at IS NULL").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *CustomerRepository) List(ctx context.Context, filter *CustomerQueryFilter) ([]model.Customer, int, error) {
	var customers []model.Customer
	q := r.db.NewSelect().
		Model(&customers).
		Where("deleted_at IS NULL")

	if filter.Code != "" {
		q = q.Where("code LIKE ?", "%"+filter.Code+"%")
	}
	if filter.Name != "" {
		q = q.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.Phone != "" {
		q = q.Where("phone LIKE ?", "%"+filter.Phone+"%")
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}

	total, err := q.
		Order("id DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
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
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}

func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Customer)(nil)).
		Set("deleted_at = ?", time.Now()).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Exec(ctx)
	return err
}
