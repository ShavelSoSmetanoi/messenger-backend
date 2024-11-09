package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) InitFileRoutes(r *gin.RouterGroup) {
	files := r.Group("/files")
	{
		files.POST("/upload", h.UploadFileHandler)             // Загрузка файла
		files.GET("/download/:file_id", h.DownloadFileHandler) // Скачивание файла по ID
		files.DELETE("/:file_id", h.DeleteFileHandler)         // Удаление файла по ID
		files.GET("/:file_id/info", h.GetFileInfoHandler)      // Получение информации о файле по ID
	}
}

// UploadFileHandler handles file upload
func (h *Handler) UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	fileID, err := h.services.File.UploadFile(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_id": fileID})
}

// DownloadFileHandler handles file download by ID
// DownloadFileHandler handles file download by ID
func (h *Handler) DownloadFileHandler(c *gin.Context) {
	fileID := c.Param("file_id")

	// Получаем файл и его метаданные
	file, err := h.services.File.DownloadFile(c, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Получаем метаданные файла для определения типа контента
	fileInfo, err := h.services.File.GetFileInfo(c, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error retrieving file metadata"})
		return
	}

	// Отправляем файл клиенту
	c.DataFromReader(http.StatusOK, fileInfo.Size, fileInfo.FileType, file.Content, nil)
}

// DeleteFileHandler handles file deletion by ID
func (h *Handler) DeleteFileHandler(c *gin.Context) {
	fileID := c.Param("file_id")

	err := h.services.File.DeleteFile(c, fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// GetFileInfoHandler returns information about a file by ID
func (h *Handler) GetFileInfoHandler(c *gin.Context) {
	fileID := c.Param("file_id")

	fileInfo, err := h.services.File.GetFileInfo(c, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_info": fileInfo})
}
