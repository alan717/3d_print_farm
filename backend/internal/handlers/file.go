package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/3d-printing-farm/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileHandler struct {
	db        *gorm.DB
	uploadDir string
}

func NewFileHandler(db *gorm.DB) *FileHandler {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	
	// 确保上传目录存在
	os.MkdirAll(uploadDir, 0755)
	
	return &FileHandler{
		db:        db,
		uploadDir: uploadDir,
	}
}

// 文件上传
func (h *FileHandler) Upload(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}
	defer file.Close()

	// 验证文件类型
	filename := header.Filename
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".stl" && ext != ".obj" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持STL和OBJ格式文件"})
		return
	}

	// 验证文件大小（50MB）
	if header.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件大小不能超过50MB"})
		return
	}

	// 生成唯一文件名
	fileID := uuid.New()
	savedFilename := fmt.Sprintf("%s%s", fileID.String(), ext)
	savedPath := filepath.Join(h.uploadDir, savedFilename)

	// 保存文件
	dst, err := os.Create(savedPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 保存文件记录到数据库
	userUUID, _ := uuid.Parse(userID)
	fileRecord := models.File{
		ID:           fileID,
		UserID:       userUUID,
		Filename:     savedFilename,
		OriginalName: filename,
		Size:         header.Size,
		MimeType:     header.Header.Get("Content-Type"),
		Path:         savedPath,
	}

	if err := h.db.Create(&fileRecord).Error; err != nil {
		// 删除已保存的文件
		os.Remove(savedPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件记录失败"})
		return
	}

	// 返回文件信息
	response := map[string]interface{}{
		"id":            fileRecord.ID,
		"filename":      fileRecord.Filename,
		"original_name": fileRecord.OriginalName,
		"size":          fileRecord.Size,
		"url":           fmt.Sprintf("/api/files/%s", fileRecord.ID),
		"created_at":    fileRecord.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// 获取文件
func (h *FileHandler) GetFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
		return
	}

	var fileRecord models.File
	if err := h.db.Where("id = ?", fileID).First(&fileRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(fileRecord.Path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileRecord.OriginalName))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", strconv.FormatInt(fileRecord.Size, 10))

	c.File(fileRecord.Path)
}

// 删除文件
func (h *FileHandler) DeleteFile(c *gin.Context) {
	userID := c.GetString("user_id")
	fileID := c.Param("id")
	
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件ID不能为空"})
		return
	}

	var fileRecord models.File
	if err := h.db.Where("id = ? AND user_id = ?", fileID, userID).First(&fileRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在或无权限"})
		return
	}

	// 删除物理文件
	if err := os.Remove(fileRecord.Path); err != nil {
		// 记录错误但继续删除数据库记录
		fmt.Printf("删除物理文件失败: %v\n", err)
	}

	// 删除数据库记录
	if err := h.db.Delete(&fileRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件删除成功"})
}