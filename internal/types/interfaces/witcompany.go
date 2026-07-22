package interfaces

import (
	"context"

	"github.com/Tencent/WeKnora/internal/types"
)

type WitCompanyService interface {
	CreateWitCompany(ctx context.Context, companyName, companyCode, address, contactPerson, createdBy string) (*types.WitCompany, error)
	GetWitCompanyByID(ctx context.Context, id int64) (*types.WitCompany, error)
	ListWitCompanies(ctx context.Context, page *types.Pagination, keyword string) (*types.PageResult, error)
	UpdateWitCompany(ctx context.Context, id int64, companyName, companyCode, address, contactPerson *string) (*types.WitCompany, error)
	DeleteWitCompany(ctx context.Context, id int64) error
}

type WitCompanyRepository interface {
	Create(ctx context.Context, company *types.WitCompany) error
	GetByID(ctx context.Context, id int64) (*types.WitCompany, error)
	GetByCompanyCode(ctx context.Context, code string) (*types.WitCompany, error)
	List(ctx context.Context, page *types.Pagination, keyword string) ([]*types.WitCompany, int64, error)
	Update(ctx context.Context, company *types.WitCompany) error
	Delete(ctx context.Context, id int64) error
}
