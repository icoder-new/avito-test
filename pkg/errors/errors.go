package errors

import (
	"errors"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrTenderNotFound           = errors.New("tender not found")
	ErrBidNotFound              = errors.New("bid not found")
	ErrOrganizationNotFound     = errors.New("organization not found")
	ErrEmployeeNotFound         = errors.New("employee not found")
	ErrUserNotIsOwnerOrNotExist = errors.New("user is not owner or user/organization does not exist")
	ErrNotEnoughRights          = errors.New("not enough rights")
)

var ErrorStatusMapping = map[error]int{
	ErrTenderNotFound:       http.StatusNotFound,
	ErrOrganizationNotFound: http.StatusNotFound,
	ErrEmployeeNotFound:     http.StatusUnauthorized,
	ErrNotEnoughRights:      http.StatusForbidden,
}

func HandleError(c *gin.Context, err error, log logger.ILogger) {
	log.Error("service error", err)

	for customErr, statusCode := range ErrorStatusMapping {
		if errors.Is(err, customErr) {
			c.AbortWithStatusJSON(statusCode, gin.H{
				"reason": customErr.Error(),
			})
			return
		}
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"reason": err.Error(),
	})
}
