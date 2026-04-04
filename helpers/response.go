package helpers

import "github.com/gin-gonic/gin"

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SendResponseHTTP(c *gin.Context, code int, msg string, data any) {
	c.JSON(code, Response{
		Message: msg,
		Data:    data,
	})
}
