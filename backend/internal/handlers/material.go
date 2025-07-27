package handlers

import (
	"encoding/json"
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

// 材料响应结构
type MaterialResponse struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	PricePerGram float64                `json:"price_per_gram"`
	Properties   map[string]interface{} `json:"properties"`
	Available    bool                   `json:"available"`
	CreatedAt    string                 `json:"created_at"`
	UpdatedAt    string                 `json:"updated_at"`
}

// 获取所有材料
func (h *MaterialHandler) GetMaterials(c *gin.Context) {
	var materials []models.Material
	
	// 只获取可用的材料
	if err := h.db.Where("available = ?", true).Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取材料列表失败"})
		return
	}

	// 转换为响应格式
	var response []MaterialResponse
	for _, material := range materials {
		var properties map[string]interface{}
		if err := json.Unmarshal([]byte(material.Properties), &properties); err != nil {
			// 如果解析失败，创建空对象
			properties = make(map[string]interface{})
		}
		
		response = append(response, MaterialResponse{
			ID:           material.ID,
			Name:         material.Name,
			Description:  material.Description,
			PricePerGram: material.PricePerGram,
			Properties:   properties,
			Available:    material.Available,
			CreatedAt:    material.CreatedAt.Format("2006-01-02T15:04:05.999999999-07:00"),
			UpdatedAt:    material.UpdatedAt.Format("2006-01-02T15:04:05.999999999-07:00"),
		})
	}

	c.JSON(http.StatusOK, response)
}