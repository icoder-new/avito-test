package api

import (
	"net/http"

	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/api/handler"
	"git.codenrock.com/cnrprod1725373421-user-90073/AvitoTestTask/pkg/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetUpRoutes(h *handler.Handler, log logger.ILogger) *gin.Engine {
	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(
		requestid.New(
			requestid.WithGenerator(func() string {
				return "avito-test-task-" + uuid.New().String()
			}),
			requestid.WithCustomHeaderStrKey("X-Request-ID"),
		),
	)
	router.Use(gin.Recovery())
	router.Use(handler.NewMWLogger(log))

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"reason": "Method not allowed",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"reason": "Route not found",
		})
	})

	route := router.Group("/api")
	route.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	route.GET("/tenders", h.GetTenders)
	route.POST("/tenders/new", h.CreateTender)
	route.GET("/tenders/my", h.GetMyTenders)
	route.GET("/tenders/:id/status", h.GetTenderStatusById)
	route.PUT("/tenders/:id/status", h.ChangeTenderStatusById)
	route.PATCH("/tenders/:id/edit", h.EditTenderById)
	route.PUT("/tenders/:id/rollback/:version", h.RollbackTenderById)

	route.GET("/bids/my", h.GetBidsByUsername)
	route.POST("/bids/new", h.CreateBid)
	route.GET("/bids/:id/list", h.GetBidListByTenderId)
	route.GET("/bids/:id/status", h.GetBidStatus)
	route.PUT("/bids/:id/status", h.ChangeBidStatus)
	route.PATCH("/bids/:id/edit", h.ChangeBidById)
	route.PUT("/bids/:id/submit_decision", h.SendDesicionByBid)
	route.PUT("/bids/:id/feedback", h.SendFeedbackByBid)
	route.PUT("/bids/:id/rollback/:version", h.RollbackBid)
	route.GET("/bids/:id/reviews", h.GetFeedbacksByTenderAndAuthor)

	return router
}
