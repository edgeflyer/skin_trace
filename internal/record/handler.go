package record

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *gin.Context) {
	// 限制整个请求体为 20 MiB，包含照片和表单内容。
	const maxRequestSize = 20 << 20
	c.Request.Body = http.MaxBytesReader(
		c.Writer,
		c.Request.Body,
		maxRequestSize,
	)

	header, err := c.FormFile("photo")
	if err != nil {
		var sizeErr *http.MaxBytesError
		if errors.As(err, &sizeErr) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error": "上传内容不能超过 20 MiB",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请使用 multipart/form-data 上传 photo 文件",
		})
		return
	}

	photo, err := header.Open()
	if err != nil {
		log.Printf("open uploaded photo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "读取上传照片失败",
		})
		return
	}
	defer photo.Close()

	rec, err := h.service.Create(
		c.Request.Context(),
		photo,
		c.PostForm("note"),
	)
	if err != nil {
		if errors.Is(err, ErrInvalidPhoto) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "请上传 JPEG 或 PNG 格式的照片",
			})
			return
		}

		log.Printf("create record: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存记录失败",
		})
		return
	}

	c.JSON(http.StatusCreated, toResponse(*rec))
}

func (h *Handler) List(c *gin.Context) {
	records, err := h.service.List(c.Request.Context())
	if err != nil {
		log.Printf("list records: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "查询记录失败",
		})
		return
	}

	response := make([]Response, 0, len(records))
	for _, rec := range records {
		response = append(response, toResponse(rec))
	}

	c.JSON(http.StatusOK, response)
}
