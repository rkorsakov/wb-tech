package handlers

import (
	"L0/internal/cache"
	"L0/internal/db/postgres"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderHandler struct {
	storage    *postgres.Storage
	orderCache *cache.OrderCache
}

func NewOrderHandler(storage *postgres.Storage, orderCache *cache.OrderCache) *OrderHandler {
	return &OrderHandler{storage: storage, orderCache: orderCache}
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	ctx := c.Request.Context()
	if cachedOrder, ok := h.orderCache.Get(orderID); ok {
		c.JSON(http.StatusOK, cachedOrder)
		return
	}
	order, err := h.storage.GetOrder(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get order: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Просмотр заказов",
	})
}
