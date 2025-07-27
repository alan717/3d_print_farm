package handlers

import (
	"net/http"

	"github.com/3d-printing-farm/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MaterialHandler struct {
	db *gorm.DB
}

func NewMaterialHandler(db *gorm.DB) *MaterialHandler {
	return &MaterialHandler{db: db}
}

// 获取所有材料
func (h *MaterialHandler) GetMaterials(c *gin.Context) {
	var materials []models.Material
	
	// 只获取可用的材料
	if err := h.db.Where("available = ?", true).Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取材料列表失败"})
		return
	}

	c.JSON(http.StatusOK, materials)
}