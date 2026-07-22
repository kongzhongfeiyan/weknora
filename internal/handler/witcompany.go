package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Tencent/WeKnora/internal/errors"
	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	"github.com/Tencent/WeKnora/internal/types/interfaces"
	secutils "github.com/Tencent/WeKnora/internal/utils"
)

type WitCompanyHandler struct {
	witCompanyService interfaces.WitCompanyService
}

func NewWitCompanyHandler(witCompanyService interfaces.WitCompanyService) *WitCompanyHandler {
	return &WitCompanyHandler{witCompanyService: witCompanyService}
}

type createWitCompanyRequest struct {
	CompanyName   string `json:"company_name"   binding:"required"`
	CompanyCode   string `json:"company_code"   binding:"required"`
	Address       string `json:"address"`
	ContactPerson string `json:"contact_person"`
	CreatedBy     string `json:"created_by"`
}

type updateWitCompanyRequest struct {
	CompanyName   *string `json:"company_name"`
	CompanyCode   *string `json:"company_code"`
	Address       *string `json:"address"`
	ContactPerson *string `json:"contact_person"`
}

// CreateWitCompany godoc
// @Summary      创建企业
// @Description  创建新的企业信息
// @Tags         企业管理
// @Accept       json
// @Produce      json
// @Param        request  body      createWitCompanyRequest  true  "企业信息"
// @Success      200      {object}  map[string]interface{}   "创建的企业"
// @Failure      400      {object}  errors.AppError          "请求参数错误"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /witcompanies [post]
func (h *WitCompanyHandler) CreateWitCompany(c *gin.Context) {
	ctx := c.Request.Context()

	var req createWitCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind create witcompany payload", err)
		c.Error(errors.NewBadRequestError("请求参数不合法").WithDetails(err.Error()))
		return
	}

	result, err := h.witCompanyService.CreateWitCompany(ctx,
		secutils.SanitizeForLog(req.CompanyName),
		secutils.SanitizeForLog(req.CompanyCode),
		secutils.SanitizeForLog(req.Address),
		secutils.SanitizeForLog(req.ContactPerson),
		secutils.SanitizeForLog(req.CreatedBy),
	)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetWitCompanyByID godoc
// @Summary      获取企业详情
// @Description  根据ID获取企业详细信息
// @Tags         企业管理
// @Accept       json
// @Produce      json
// @Param        id   path      int                    true  "企业ID"
// @Success      200  {object}  map[string]interface{} "企业详情"
// @Failure      400  {object}  errors.AppError        "请求参数错误"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /witcompanies/{id} [get]
func (h *WitCompanyHandler) GetWitCompanyByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(ctx, "Invalid witcompany ID", err)
		c.Error(errors.NewBadRequestError("企业ID无效"))
		return
	}

	result, err := h.witCompanyService.GetWitCompanyByID(ctx, id)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"id": id})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// ListWitCompanies godoc
// @Summary      获取企业列表
// @Description  分页获取企业列表，支持关键词搜索
// @Tags         企业管理
// @Accept       json
// @Produce      json
// @Param        page       query     int                    false "页码"
// @Param        page_size  query     int                    false "每页数量"
// @Param        keyword    query     string                 false "关键词搜索"
// @Success      200        {object}  map[string]interface{} "企业列表"
// @Failure      400        {object}  errors.AppError        "请求参数错误"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /witcompanies [get]
func (h *WitCompanyHandler) ListWitCompanies(c *gin.Context) {
	ctx := c.Request.Context()

	var page types.Pagination
	if err := c.ShouldBindQuery(&page); err != nil {
		logger.Error(ctx, "Failed to bind pagination query", err)
		c.Error(errors.NewBadRequestError("分页参数不合法").WithDetails(err.Error()))
		return
	}

	keyword := secutils.SanitizeForLog(c.Query("keyword"))

	result, err := h.witCompanyService.ListWitCompanies(ctx, &page, keyword)
	if err != nil {
		logger.ErrorWithFields(ctx, err, nil)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// UpdateWitCompany godoc
// @Summary      更新企业
// @Description  更新企业信息（部分更新）
// @Tags         企业管理
// @Accept       json
// @Produce      json
// @Param        id       path      int                     true  "企业ID"
// @Param        request  body      updateWitCompanyRequest true  "更新信息"
// @Success      200      {object}  map[string]interface{}  "更新后的企业"
// @Failure      400      {object}  errors.AppError         "请求参数错误"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /witcompanies/{id} [put]
func (h *WitCompanyHandler) UpdateWitCompany(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(ctx, "Invalid witcompany ID", err)
		c.Error(errors.NewBadRequestError("企业ID无效"))
		return
	}

	var req updateWitCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error(ctx, "Failed to bind update witcompany payload", err)
		c.Error(errors.NewBadRequestError("请求参数不合法").WithDetails(err.Error()))
		return
	}

	result, err := h.witCompanyService.UpdateWitCompany(ctx, id,
		req.CompanyName, req.CompanyCode, req.Address, req.ContactPerson)
	if err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"id": id})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// DeleteWitCompany godoc
// @Summary      删除企业
// @Description  根据ID删除企业
// @Tags         企业管理
// @Accept       json
// @Produce      json
// @Param        id   path      int                    true  "企业ID"
// @Success      200  {object}  map[string]interface{} "删除成功"
// @Failure      400  {object}  errors.AppError        "请求参数错误"
// @Security     Bearer
// @Security     ApiKeyAuth
// @Router       /witcompanies/{id} [delete]
func (h *WitCompanyHandler) DeleteWitCompany(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logger.Error(ctx, "Invalid witcompany ID", err)
		c.Error(errors.NewBadRequestError("企业ID无效"))
		return
	}

	if err := h.witCompanyService.DeleteWitCompany(ctx, id); err != nil {
		logger.ErrorWithFields(ctx, err, map[string]interface{}{"id": id})
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
