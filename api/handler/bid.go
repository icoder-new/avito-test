package handler

import (
	"net/http"
	"strconv"

	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/internal/dto"
	customerrors "git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) GetBidsByUsername(c *gin.Context) {
	const fn = "hadnler.GetBidsByUsername"

	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	username := c.Query("username")

	limit, err := strconv.Atoi(limitStr)
	if limitStr == "" || limitStr == " " {
		limit = 5
		err = nil
	}
	if err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid limit parameter, must be integer",
		})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if offsetStr == "" || offsetStr == " " {
		offset = 0
		err = nil
	}

	if err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"reason": "Invalid offset parameter, must be integer",
		})
		return
	}

	response, err := h.svc.Bid().GetBidsByFiltering(limit, offset, username)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateBid(c *gin.Context) {
	const fn = "handler.CreateBid"

	var bidRequest dto.BidRequest
	if err := c.ShouldBindJSON(&bidRequest); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": "Invalid input"})
		return
	}

	bidResponse, err := h.svc.Bid().CreateBid(bidRequest)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusCreated, bidResponse)
}

func (h *Handler) GetBidListByTenderId(c *gin.Context) {
	const fn = "handler.GetBidListByTenderId"
	tenderId := c.Param("id")
	username := c.Query("username")

	bids, err := h.svc.Bid().GetBidListByTenderId(tenderId, username)
	if err != nil {
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, bids)
}

func (h *Handler) GetBidStatus(c *gin.Context) {
	const fn = "handler.GetBidStatus"

	bidId := c.Param("id")
	if bidId == "" {
		h.log.Error(fn, "bidId is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": "bidId is required"})
		return
	}

	username := c.Query("username")
	if username == "" {
		h.log.Error(fn, "username is required")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": "username is required"})
		return
	}

	status, err := h.svc.Bid().GetBidStatus(bidId, username)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (h *Handler) ChangeBidStatus(c *gin.Context) {
	const fn = "handler.ChangeBidStatus"

	bidID := c.Param("id")
	status := c.Query("status")
	username := c.Query("username")

	bidStatusRequest := dto.BidStatusRequest{
		BidID:    bidID,
		Status:   dto.Status(status),
		Username: username,
	}

	validate := validator.New()
	if err := validate.Struct(bidStatusRequest); err != nil {
		h.log.Error(fn, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	bidResponse, err := h.svc.Bid().ChangeBidStatus(&bidStatusRequest)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, bidResponse)
}

func (h *Handler) ChangeBidById(c *gin.Context) {
	const fn = "handler.ChangeBidById"

	bidID := c.Param("id")
	username := c.Query("username")

	var bidRequest dto.ChangeBidRequest
	if err := c.ShouldBindJSON(&bidRequest); err != nil {
		h.log.Error(fn, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": "Invalid input"})
		return
	}

	bidRequest.BidID = bidID
	bidRequest.Username = username

	validate := validator.New()
	if err := validate.Struct(bidRequest); err != nil {
		h.log.Error(fn, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	bidResponse, err := h.svc.Bid().ChangeBidById(&bidRequest)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, bidResponse)
}

func (h *Handler) SendDesicionByBid(c *gin.Context) {
	const fn = "handler.SendDesicionByBid"

	bidID := c.Param("id")
	decision := c.Query("decision")
	username := c.Query("username")

	req := dto.SendBidDesicionRequest{
		BidID:    bidID,
		Username: username,
		Decision: dto.Decision(decision),
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Error(fn, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	bidResponse, err := h.svc.Bid().SendDesicionByBid(&req)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, bidResponse)
}

func (h *Handler) SendFeedbackByBid(c *gin.Context) {
	const fn = "handler.SendFeedbackByBid"

	bidID := c.Param("id")
	bidFeedback := c.Query("bidFeedback")
	username := c.Query("username")

	req := dto.SendBidFeedbackRequest{
		BidID:       bidID,
		Username:    username,
		BidFeedback: bidFeedback,
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		h.log.Error(fn, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	bidResponse, err := h.svc.Bid().SendFeedbackByBid(&req)
	if err != nil {
		h.log.Error(fn, err)
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, bidResponse)
}

func (h *Handler) GetFeedbacksByTenderAndAuthor(c *gin.Context) {
	const fn = "handler.GetFeedbacksByTenderAndAuthor"

	tenderID := c.Param("id")
	authorUsername := c.Query("authorUsername")
	requesterUsername := c.Query("requesterUsername")

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		h.log.Error(fn, "invalid limit parameter")
		customerrors.HandleError(c, err, h.log)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		h.log.Error(fn, "invalid offset parameter")
		customerrors.HandleError(c, err, h.log)
		return
	}

	feedbacks, err := h.svc.Bid().GetFeedbacksByTenderAndAuthor(tenderID, authorUsername, requesterUsername, limit, offset)
	if err != nil {
		h.log.Error(fn, err.Error())
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, feedbacks)
}

func (h *Handler) RollbackBid(c *gin.Context) {
	bidId := c.Param("id")
	versionStr := c.Param("version")
	username := c.Query("username")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Invalid version format"})
		return
	}

	bid, err := h.svc.Bid().RollbackBidById(bidId, username, version)
	if err != nil {
		customerrors.HandleError(c, err, h.log)
		return
	}

	c.JSON(http.StatusOK, bid)
}
