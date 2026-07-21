package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloWorldHandler handles hello world requests
type HelloWorldHandler struct{}

// NewHelloWorldHandler creates a new hello world handler
func NewHelloWorldHandler() *HelloWorldHandler {
	return &HelloWorldHandler{}
}

// HelloWorld returns a hello world message.
//
// HelloWorld godoc
// @Summary      Hello World
// @Description  返回 hello world 消息
// @Tags         测试
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "hello world"
// @Router       /hello-world [get]
func (h *HelloWorldHandler) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"message": "hello world",
		},
	})
}
