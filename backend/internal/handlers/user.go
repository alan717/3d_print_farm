package handlers

import (
	"net/http"

	"github.com/3d-printing-farm/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

type UpdateProfileRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

// 获取用户档案
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var user models.User
	if err := h.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 清除密码字段
	user.Password = ""
	c.JSON(http.StatusOK, user)
}

// 更新用户档案
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否被其他用户使用
	var existingUser models.User
	if err := h.db.Where("username = ? AND id != ?", req.Username, userID).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "用户名已被使用"})
		return
	}

	// 更新用户信息
	userUUID, _ := uuid.Parse(userID)
	updates := models.User{
		Username: req.Username,
		Phone:    req.Phone,
		Address:  req.Address,
	}

	if err := h.db.Model(&models.User{}).Where("id = ?", userUUID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	// 返回更新后的用户信息
	var user models.User
	h.db.Where("id = ?", userUUID).First(&user)
	user.Password = ""

	c.JSON(http.StatusOK, user)
}