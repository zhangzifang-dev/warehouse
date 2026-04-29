package service

import (
	"context"
	"encoding/json"

	"warehouse/internal/model"
	apperrors "warehouse/internal/pkg/errors"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *model.Customer) error
	GetByID(ctx context.Context, id int64) (*model.Customer, error)
	GetByCode(ctx context.Context, code string) (*model.Customer, error)
	List(ctx context.Context, page, pageSize int, keyword string) ([]model.Customer, int, error)
	Update(ctx context.Context, customer *model.Customer) error
	Delete(ctx context.Context, id int64) error
}

type CustomerService struct {
	customerRepo CustomerRepository
	auditLogger  AuditLogger
}

func NewCustomerService(customerRepo CustomerRepository, auditLogger AuditLogger) *CustomerService {
	return &CustomerService{
		customerRepo: customerRepo,
		auditLogger:  auditLogger,
	}
}

type CreateCustomerInput struct {
	Name    string
	Code    string
	Contact string
	Phone   string
	Email   string
	Address string
	Status  int
}

type UpdateCustomerInput struct {
	Name    *string
	Code    *string
	Contact *string
	Phone   *string
	Email   *string
	Address *string
	Status  *int
}

type ListCustomersResult struct {
	Customers []model.Customer
	Total     int
}

func (s *CustomerService) Create(ctx context.Context, input *CreateCustomerInput) (*model.Customer, error) {
	if input.Name == "" {
		return nil, apperrors.NewAppError(apperrors.CodeBadRequest, "customer name is required")
	}

	if input.Code != "" {
		existing, err := s.customerRepo.GetByCode(ctx, input.Code)
		if err == nil && existing != nil {
			return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "customer code already exists")
		}
	}

	customer := &model.Customer{
		Name:    input.Name,
		Code:    input.Code,
		Contact: input.Contact,
		Phone:   input.Phone,
		Email:   input.Email,
		Address: input.Address,
		Status:  input.Status,
	}

	if customer.Status == 0 {
		customer.Status = model.CustomerStatusActive
	}

	err := s.customerRepo.Create(ctx, customer)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to create customer")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(customer)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "customers",
			RecordID:   customer.ID,
			Action:     "create",
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: customer.CreatedBy,
		})
	}

	return customer, nil
}

func (s *CustomerService) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "customer not found")
	}
	return customer, nil
}

func (s *CustomerService) List(ctx context.Context, page, pageSize int, keyword string) (*ListCustomersResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	customers, total, err := s.customerRepo.List(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to list customers")
	}

	return &ListCustomersResult{
		Customers: customers,
		Total:     total,
	}, nil
}

func (s *CustomerService) Update(ctx context.Context, id int64, input *UpdateCustomerInput) (*model.Customer, error) {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeRecordNotFound, "customer not found")
	}

	var oldValue []byte
	if s.auditLogger != nil {
		oldValue, _ = json.Marshal(customer)
	}

	if input.Code != nil && *input.Code != customer.Code {
		if *input.Code != "" {
			existing, err := s.customerRepo.GetByCode(ctx, *input.Code)
			if err == nil && existing != nil {
				return nil, apperrors.NewAppError(apperrors.CodeDuplicateEntry, "customer code already exists")
			}
		}
		customer.Code = *input.Code
	}

	if input.Name != nil {
		customer.Name = *input.Name
	}
	if input.Contact != nil {
		customer.Contact = *input.Contact
	}
	if input.Phone != nil {
		customer.Phone = *input.Phone
	}
	if input.Email != nil {
		customer.Email = *input.Email
	}
	if input.Address != nil {
		customer.Address = *input.Address
	}
	if input.Status != nil {
		customer.Status = *input.Status
	}

	err = s.customerRepo.Update(ctx, customer)
	if err != nil {
		return nil, apperrors.NewAppError(apperrors.CodeInternalError, "failed to update customer")
	}

	if s.auditLogger != nil {
		newValue, _ := json.Marshal(customer)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "customers",
			RecordID:   customer.ID,
			Action:     "update",
			OldValue:   map[string]any{"data": string(oldValue)},
			NewValue:   map[string]any{"data": string(newValue)},
			OperatedBy: customer.UpdatedBy,
		})
	}

	return customer, nil
}

func (s *CustomerService) Delete(ctx context.Context, id int64) error {
	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeRecordNotFound, "customer not found")
	}

	err = s.customerRepo.Delete(ctx, id)
	if err != nil {
		return apperrors.NewAppError(apperrors.CodeInternalError, "failed to delete customer")
	}

	if s.auditLogger != nil {
		oldValue, _ := json.Marshal(customer)
		s.auditLogger.Log(ctx, &CreateAuditLogInput{
			TableName:  "customers",
			RecordID:   customer.ID,
			Action:     "delete",
			OldValue:   map[string]any{"data": string(oldValue)},
			OperatedBy: customer.UpdatedBy,
		})
	}

	return nil
}
