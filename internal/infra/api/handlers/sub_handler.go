package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/elaurentium/exilium-blog-backend/internal/domain/services"
)

type SubHandler struct {
	subService *services.SubService
}

func NewSubHandler(subService *services.SubService) *SubHandler {
	return &SubHandler{subService: subService}
}

type CreateSubRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Rules       []string `json:"rules" binding:"required"`
	IsPrivate   bool     `json:"is_private"`
}

type UpdateSubRequest struct {
	Description string   `json:"description"`
	Rules       []string `json:"rules"`
	IsPrivate   bool     `json:"is_private"`
	BannerURL   string   `json:"banner_url"`
	IconURL     string   `json:"icon_url"`
}

func (h *SubHandler) createSub(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, errorResponse("unauthorized"))
		return
	}

	var createReq CreateSubRequest
	if err := ctx.ShouldBindJSON(&createReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	sub, err := h.subService.CreateSub(
		ctx.Request.Context(),
		createReq.Name,
		createReq.Description,
		createReq.Rules,
		userID.(uuid.UUID),
		createReq.IsPrivate,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, sub)
}

func errorResponse(message string) gin.H {
	return gin.H{"error": message}
}

func (h *SubHandler) UpdateSub(c *gin.Context) {
	userID, userExists := c.Get("user_id")
	if !userExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subID, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sub ID"})
		return
	}

	var updateReq UpdateSubRequest
	if bindErr := c.ShouldBindJSON(&updateReq); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	updatedSub, updateErr := h.subService.UpdateSub(
		c.Request.Context(),
		subID,
		userID.(uuid.UUID),
		updateReq.Description,
		updateReq.Rules,
		updateReq.IsPrivate,
		updateReq.BannerURL,
		updateReq.IconURL,
	)
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedSub)
}

func (h *SubHandler) DeleteSub(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	subID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sub ID"})
		return
	}

	if err := h.subService.DeleteSub(c.Request.Context(), subID, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *SubHandler) GetSub(c *gin.Context) {
	subID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sub ID"})
		return
	}

	sub, err := h.subService.GetSub(c.Request.Context(), subID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubHandler) GetSubByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sub name is required"})
		return
	}

	sub, err := h.subService.GetSubByName(c.Request.Context(), name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubHandler) ListSubs(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
		return
	}

	subs, err := h.subService.ListSubs(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *SubHandler) GetTrendingSubreddits(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	subreddits, err := h.subService.GetTrendingSub(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subreddits)
}


func (h *SubHandler) CreateSub(c *gin.Context) {
	userID, userExists := c.Get("user_id")
	if !userExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var createReq CreateSubRequest
	if bindErr := c.ShouldBindJSON(&createReq); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	sub, err := h.subService.CreateSub(
	c.Request.Context(),
	createReq.Name,
	createReq.Description,
	createReq.Rules,
	userID.(uuid.UUID),
	createReq.IsPrivate,
)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}
