package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/3d-printing-farm/backend/internal/models"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	db *gorm.DB
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// AdminMiddleware 创建管理员中间件
func AdminMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
			c.Abort()
			return
		}

		var user models.User
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		if user.Role != models.UserRoleAdmin && user.Role != models.UserRoleSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		c.Set("admin_user", user)
		c.Next()
	}
}

// logAdminAction 记录管理员操作
func (h *AdminHandler) logAdminAction(c *gin.Context, action, resourceType, resourceID string, details interface{}) {
	admin, _ := c.Get("admin_user")
	adminUser := admin.(models.User)

	log := models.AdminLog{
		AdminID:      &adminUser.ID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.GetHeader("User-Agent"),
	}

	if details != nil {
		log.Details = fmt.Sprintf("%+v", details)
	}

	h.db.Create(&log)
}

// GetDashboardStats 获取仪表板统计数据
func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	var stats struct {
		TotalUsers       int64   `json:"total_users"`
		TotalOrders      int64   `json:"total_orders"`
		TotalRevenue     float64 `json:"total_revenue"`
		TotalFiles       int64   `json:"total_files"`
		ActiveOrders     int64   `json:"active_orders"`
		CompletedOrders  int64   `json:"completed_orders"`
		TodayOrders      int64   `json:"today_orders"`
		TodayRevenue     float64 `json:"today_revenue"`
		RecentOrders     []models.Order `json:"recent_orders"`
		TopMaterials     []struct {
			Name  string `json:"name"`
			Count int64  `json:"count"`
		} `json:"top_materials"`
	}

	// Get basic counts
	h.db.Model(&models.User{}).Count(&stats.TotalUsers)
	h.db.Model(&models.Order{}).Count(&stats.TotalOrders)
	h.db.Model(&models.File{}).Count(&stats.TotalFiles)
	h.db.Model(&models.Order{}).Where("status IN ?", []string{"pending", "printing"}).Count(&stats.ActiveOrders)
	h.db.Model(&models.Order{}).Where("status = ?", "completed").Count(&stats.CompletedOrders)

	// Get total revenue
	h.db.Model(&models.Order{}).Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&stats.TotalRevenue)

	// Get today's stats
	today := time.Now().Format("2006-01-02")
	h.db.Model(&models.Order{}).Where("DATE(created_at) = ?", today).Count(&stats.TodayOrders)
	h.db.Model(&models.Order{}).Where("DATE(created_at) = ?", today).Select("COALESCE(SUM(total_amount), 0)").Row().Scan(&stats.TodayRevenue)

	// Get recent orders
	h.db.Preload("User").Order("created_at DESC").Limit(10).Find(&stats.RecentOrders)

	// Get top materials
	h.db.Raw(`
		SELECT m.name, COUNT(oi.id) as count 
		FROM materials m 
		JOIN order_items oi ON m.id = oi.material_id 
		GROUP BY m.name 
		ORDER BY count DESC 
		LIMIT 5
	`).Scan(&stats.TopMaterials)

	h.logAdminAction(c, "VIEW", "dashboard", "", nil)
	c.JSON(http.StatusOK, stats)
}

// GetUsers 获取用户列表
func (h *AdminHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	offset := (page - 1) * limit

	var users []models.User
	var total int64

	query := h.db.Model(&models.User{})
	if search != "" {
		query = query.Where("username ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users)

	// Clear passwords
	for i := range users {
		users[i].Password = ""
	}

	h.logAdminAction(c, "VIEW", "users", "", map[string]interface{}{
		"page": page, "limit": limit, "search": search,
	})

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateUserStatus 更新用户状态
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")
	
	var req struct {
		Status models.UserStatus `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	oldStatus := user.Status
	user.Status = req.Status

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	h.logAdminAction(c, "UPDATE", "user", userID, map[string]interface{}{
		"old_status": oldStatus,
		"new_status": req.Status,
	})

	c.JSON(http.StatusOK, gin.H{"message": "用户状态更新成功"})
}

// GetOrders 获取订单列表
func (h *AdminHandler) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	offset := (page - 1) * limit

	var orders []models.Order
	var total int64

	query := h.db.Model(&models.Order{}).Preload("User")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders)

	h.logAdminAction(c, "VIEW", "orders", "", map[string]interface{}{
		"page": page, "limit": limit, "status": status,
	})

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// UpdateOrderStatus 更新订单状态
func (h *AdminHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	
	var req struct {
		Status models.OrderStatus `json:"status" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := h.db.First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	oldStatus := order.Status
	order.Status = req.Status

	if err := h.db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	h.logAdminAction(c, "UPDATE", "order", orderID, map[string]interface{}{
		"old_status": oldStatus,
		"new_status": req.Status,
	})

	c.JSON(http.StatusOK, gin.H{"message": "订单状态更新成功"})
}

// CreateMaterial 创建材料
func (h *AdminHandler) CreateMaterial(c *gin.Context) {
	var material models.Material
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	h.logAdminAction(c, "CREATE", "material", strconv.Itoa(int(material.ID)), material)
	c.JSON(http.StatusCreated, material)
}

// UpdateMaterial 更新材料
func (h *AdminHandler) UpdateMaterial(c *gin.Context) {
	materialID := c.Param("id")
	
	var material models.Material
	if err := h.db.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "材料不存在"})
		return
	}

	oldMaterial := material
	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	h.logAdminAction(c, "UPDATE", "material", materialID, map[string]interface{}{
		"old": oldMaterial,
		"new": material,
	})

	c.JSON(http.StatusOK, material)
}

// DeleteMaterial 删除材料
func (h *AdminHandler) DeleteMaterial(c *gin.Context) {
	materialID := c.Param("id")
	
	var material models.Material
	if err := h.db.First(&material, materialID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "材料不存在"})
		return
	}

	if err := h.db.Delete(&material).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	h.logAdminAction(c, "DELETE", "material", materialID, material)
	c.JSON(http.StatusOK, gin.H{"message": "材料删除成功"})
}

// GetSystemSettings 获取系统设置
func (h *AdminHandler) GetSystemSettings(c *gin.Context) {
	var settings []models.SystemSetting
	h.db.Find(&settings)

	h.logAdminAction(c, "VIEW", "settings", "", nil)
	c.JSON(http.StatusOK, settings)
}

// UpdateSystemSetting 更新系统设置
func (h *AdminHandler) UpdateSystemSetting(c *gin.Context) {
	key := c.Param("key")
	
	var req struct {
		Value string `json:"value" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var setting models.SystemSetting
	if err := h.db.Where("setting_key = ?", key).First(&setting).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设置不存在"})
		return
	}

	oldValue := setting.SettingValue
	setting.SettingValue = req.Value

	if err := h.db.Save(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	h.logAdminAction(c, "UPDATE", "setting", key, map[string]interface{}{
		"old_value": oldValue,
		"new_value": req.Value,
	})

	c.JSON(http.StatusOK, gin.H{"message": "设置更新成功"})
}

// GetAdminLogs 获取管理员操作日志
func (h *AdminHandler) GetAdminLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	action := c.Query("action")
	resourceType := c.Query("resource_type")

	offset := (page - 1) * limit

	var logs []models.AdminLog
	var total int64

	query := h.db.Model(&models.AdminLog{}).Preload("Admin")
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	query.Count(&total)
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}