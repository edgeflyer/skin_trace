package agent

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	agent *Agent
}

func NewHandler(agent *Agent) *Handler {
	return &Handler{agent: agent}
}

func (h *Handler) Chat(c *gin.Context) {
	log.Println("接收到调用请求")
	var req struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求体必须是合法的 JSON，message 必须是字符串",
		})
		return
	}

	question := strings.TrimSpace(req.Message)
	if question == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请输入问题",
		})
		return
	}

	answer, err := h.agent.Run(c.Request.Context(), question)
	if err != nil {
		log.Printf("agent run: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "对话处理失败，请查看服务端日志",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"answer": answer,
	})
}
