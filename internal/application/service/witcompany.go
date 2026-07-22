package service

import (
	"context"
	"errors"
	"strings"

	werrors "github.com/Tencent/WeKnora/internal/errors"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	"gorm.io/gorm"
)

type witCompanyService struct {
	repo interfaces.WitCompanyRepository
}

func NewWitCompanyService(repo interfaces.WitCompanyRepository) (interfaces.WitCompanyService, error) {
	return &witCompanyService{repo: repo}, nil
}

func (s *witCompanyService) CreateWitCompany(ctx context.Context, companyName, companyCode, address, contactPerson, createdBy string) (*types.WitCompany, error) {
	companyName = strings.TrimSpace(companyName)
	companyCode = strings.TrimSpace(companyCode)
	address = strings.TrimSpace(address)
	contactPerson = strings.TrimSpace(contactPerson)
	createdBy = strings.TrimSpace(createdBy)

	if companyName == "" {
		return nil, werrors.NewBadRequestError("企业名称不能为空")
	}
	if companyCode == "" {
		return nil, werrors.NewBadRequestError("企业编码不能为空")
	}

	existing, err := s.repo.GetByCompanyCode(ctx, companyCode)
	if err == nil && existing != nil {
		return nil, werrors.NewConflictError("企业编码已存在")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	company := &types.WitCompany{
		CompanyName:   companyName,
		CompanyCode:   companyCode,
		Address:       address,
		ContactPerson: contactPerson,
		CreatedBy:     createdBy,
	}
	if err := s.repo.Create(ctx, company); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"company_code": companyCode})
		return nil, err
	}
	return company, nil
}

func (s *witCompanyService) GetWitCompanyByID(ctx context.Context, id int64) (*types.WitCompany, error) {
	if id <= 0 {
		return nil, werrors.NewBadRequestError("企业ID无效")
	}
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, werrors.NewNotFoundError("企业不存在")
		}
		return nil, err
	}
	return company, nil
}

func (s *witCompanyService) ListWitCompanies(ctx context.Context, page *types.Pagination, keyword string) (*types.PageResult, error) {
	if page == nil {
		page = &types.Pagination{}
	}
	keyword = strings.TrimSpace(keyword)

	companies, total, err := s.repo.List(ctx, page, keyword)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		return nil, err
	}
	return types.NewPageResult(total, page, companies), nil
}

func (s *witCompanyService) UpdateWitCompany(ctx context.Context, id int64, companyName, companyCode, address, contactPerson *string) (*types.WitCompany, error) {
	if id <= 0 {
		return nil, werrors.NewBadRequestError("企业ID无效")
	}
	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, werrors.NewNotFoundError("企业不存在")
		}
		return nil, err
	}

	if companyName != nil {
		name := strings.TrimSpace(*companyName)
		if name == "" {
			return nil, werrors.NewBadRequestError("企业名称不能为空")
		}
		company.CompanyName = name
	}
	if companyCode != nil {
		code := strings.TrimSpace(*companyCode)
		if code == "" {
			return nil, werrors.NewBadRequestError("企业编码不能为空")
		}
		if code != company.CompanyCode {
			existing, err := s.repo.GetByCompanyCode(ctx, code)
			if err == nil && existing != nil {
				return nil, werrors.NewConflictError("企业编码已存在")
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			company.CompanyCode = code
		}
	}
	if address != nil {
		company.Address = strings.TrimSpace(*address)
	}
	if contactPerson != nil {
		company.ContactPerson = strings.TrimSpace(*contactPerson)
	}

	if err := s.repo.Update(ctx, company); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"id": id})
		return nil, err
	}
	return company, nil
}

func (s *witCompanyService) DeleteWitCompany(ctx context.Context, id int64) error {
	if id <= 0 {
		return werrors.NewBadRequestError("企业ID无效")
	}
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return werrors.NewNotFoundError("企业不存在")
		}
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"id": id})
		return err
	}
	return nil
}
