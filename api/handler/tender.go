package handler

import (
	"net/http"
	"strconv"

	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/dto"
	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func (h *Handler) GetTenders(c *gin.Context) {
	const fn = "handler.GetTenders"

	limitStr := c.DefaultQuery("limit", "5")
	offsetStr := c.DefaultQuery("offset", "0")
	serviceType := c.Query("service_type")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid limit parameter, must be integer",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid offset parameter, must be integer",
		})
		return
	}

	req := dto.TenderGetRequest{
		Limit:       limit,
		Offset:      offset,
		ServiceType: dto.TenderServiceType(serviceType),
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	response, err := h.svc.Tender().GetAllTenders(req)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, response)

}

func (h *Handler) CreateTender(c *gin.Context) {
	const fn = "handler.CreateTender"

	var tender dto.TenderRequest

	if err := c.ShouldBindJSON(&tender); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	validate := validator.New()

	if err := validate.Struct(tender); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	response, err := h.svc.Tender().CreateTender(tender)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetMyTenders(c *gin.Context) {
	const fn = "handler.GetMyTenders"

	limitStr := c.DefaultQuery("limit", "5")
	offsetStr := c.DefaultQuery("offset", "0")
	username := c.Query("username")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid limit parameter, limit must be integer",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid offset parameter, offset must be integer",
		})
		return
	}

	req := dto.TenderMyGetRequest{
		Limit:    limit,
		Offset:   offset,
		Username: username,
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	tenders, err := h.svc.Tender().GetMyTenders(req)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, tenders)
}

func (h *Handler) GetTenderStatusById(c *gin.Context) {
	const fn = "handler.GetTenderStatusById"

	req := dto.TenderStatusRequest{
		TenderId: uuid.MustParse(c.Param("id")),
		Username: c.Query("username"),
	}

	response, err := h.svc.Tender().GetTenderStatusById(req)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.String(http.StatusOK, response)
}

func (h *Handler) ChangeTenderStatusById(c *gin.Context) {
	const fn = "handler.ChangeTenderStatusById"

	tenderStatus := dto.Status(c.Query("status"))

	req := dto.TenderStatusRequest{
		TenderId: uuid.MustParse(c.Param("id")),
		Status:   tenderStatus,
		Username: c.Query("username"),
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	response, err := h.svc.Tender().ChangeTenderStatusById(req)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) EditTenderById(c *gin.Context) {
	const fn = "handler.EditTenderById"

	var tEdit dto.TenderEditRequestJSON
	if err := c.ShouldBindJSON(&tEdit); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	req := dto.TenderEditRequest{
		TenderId:    uuid.MustParse(c.Param("id")),
		Username:    c.Query("username"),
		Name:        tEdit.Name,
		Description: tEdit.Description,
		ServiceType: tEdit.ServiceType,
	}

	validate := validator.New()

	if err := validate.Struct(req); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if req.TenderId == uuid.Nil {
		h.log.Error(fn, "Invalid tenderId parameter, must be uuid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid tenderId parameter, must be uuid or nut null",
		})
		return
	}

	response, err := h.svc.Tender().EditTenderById(req)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) RollbackTenderById(c *gin.Context) {
	const fn = "handler.RollbackTenderById"

	validate := validator.New()

	version, err := strconv.Atoi(c.Param("version"))
	if err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid version parameter, must be integer",
		})
		return
	}

	req := dto.TenderRollbackRequest{
		TenderId: uuid.MustParse(c.Param("id")),
		Version:  version,
		Username: c.Query("username"),
	}

	if err := validate.Struct(req); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": err.Error(),
		})
		return
	}

	if req.TenderId == uuid.Nil {
		h.log.Error(fn, "Invalid tenderId parameter, must be uuid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid tenderId parameter, must be uuid or nut null",
		})
		return
	}

	response, err := h.svc.Tender().RollbackTenderById(req)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, response)
}
