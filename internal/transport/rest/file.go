package rest

import (
	"encoding/json"
	"github.com/ShavelSoSmetanoi/messenger-backend/internal/transport/Websocket"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) InitFileRoutes(r *gin.RouterGroup) {
	// Define a new route group under "/files"
	files := r.Group("/files")
	{
		files.POST("/upload/:chat_id/", h.UploadFileHandler)
		files.GET("/download/:file_id", h.DownloadFileHandler)
		files.DELETE("/:file_id", h.DeleteFileHandler)
		files.GET("/:file_id/info", h.GetFileInfoHandler)
	}
}

// UploadFileHandler handles file upload
func (h *Handler) UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	chatID := c.Param("chat_id")

	chatIDInt, err := strconv.Atoi(chatID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	fileID, err := h.services.File.UploadFile(c, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	message, participants, err := h.services.Message.SendMessage(chatIDInt, userID.(string), fileID, "file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	for _, participant := range participants {
		if participant.UserID != userID {
			sendMessage := SendMessageResponse{
				Status:  "send",
				Message: *message,
			}
			jsonData, err := json.Marshal(sendMessage)
			if err != nil {
				log.Printf("Failed to serialize notification: %v", err)
				continue
			}

			Websocket.NotifyUser(participant.UserID, string(jsonData))
		}
	}

	c.JSON(http.StatusOK, gin.H{"file_id": fileID})
}

// DownloadFileHandler handles file download by ID
func (h *Handler) DownloadFileHandler(c *gin.Context) {
	fileID := c.Param("file_id")

	file, err := h.services.File.DownloadFile(c, fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.DataFromReader(http.StatusOK, file.Size, file.FileType, file.Content, nil)
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
