package repository

import (
	"context"
	"strings"

	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"gorm.io/gorm"
)

type witCompanyRepository struct {
	db *gorm.DB
}

func NewWitCompanyRepository(db *gorm.DB) interfaces.WitCompanyRepository {
	return &witCompanyRepository{db: db}
}

func (r *witCompanyRepository) Create(ctx context.Context, company *types.WitCompany) error {
	return r.db.WithContext(ctx).Create(company).Error
}

func (r *witCompanyRepository) GetByID(ctx context.Context, id int64) (*types.WitCompany, error) {
	var company types.WitCompany
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *witCompanyRepository) GetByCompanyCode(ctx context.Context, code string) (*types.WitCompany, error) {
	var company types.WitCompany
	if err := r.db.WithContext(ctx).Where("company_code = ?", code).First(&company).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *witCompanyRepository) List(ctx context.Context, page *types.Pagination, keyword string) ([]*types.WitCompany, int64, error) {
	if page == nil {
		page = &types.Pagination{}
	}
	keyword = strings.TrimSpace(keyword)

	var total int64
	baseQuery := r.db.WithContext(ctx).Model(&types.WitCompany{})
	if keyword != "" {
		escaped := escapeLikeKeyword(keyword)
		baseQuery = baseQuery.Where("company_name LIKE ? OR company_code LIKE ?", "%"+escaped+"%", "%"+escaped+"%")
	}
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	dataQuery := r.db.WithContext(ctx)
	if keyword != "" {
		escaped := escapeLikeKeyword(keyword)
		dataQuery = dataQuery.Where("company_name LIKE ? OR company_code LIKE ?", "%"+escaped+"%", "%"+escaped+"%")
	}

	var companies []*types.WitCompany
	if err := dataQuery.
		Order("id DESC").
		Offset(page.Offset()).
		Limit(page.Limit()).
		Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

func (r *witCompanyRepository) Update(ctx context.Context, company *types.WitCompany) error {
	return r.db.WithContext(ctx).Save(company).Error
}

func (r *witCompanyRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&types.WitCompany{}).Error
}
