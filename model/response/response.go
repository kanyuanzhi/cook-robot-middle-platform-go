package response

import (
	"github.com/gin-gonic/gin"
	"github.com/kanyuanzhi/middle-platform/global"
	"log/slog"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Result(c *gin.Context, code int, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		message,
	})
}

func SuccessMessage(c *gin.Context, message string) {
	Result(c, global.FXConfig.System.SuccessCode, nil, message)
}

func SuccessData(c *gin.Context, data interface{}) {
	Result(c, global.FXConfig.System.SuccessCode, data, global.FXConfig.System.SuccessMessage)
}

func SuccessMessageData(c *gin.Context, data interface{}, message string) {
	Result(c, global.FXConfig.System.SuccessCode, data, message)
}

func ErrorMessage(c *gin.Context, message string) {
	slog.Warn(message)
	Result(c, global.FXConfig.System.ErrorCode, nil, message)
}

func ErrorData(c *gin.Context, data interface{}) {
	Result(c, global.FXConfig.System.ErrorCode, data, global.FXConfig.System.ErrorMessage)
}

func ErrorMessageData(c *gin.Context, data interface{}, message string) {
	Result(c, global.FXConfig.System.ErrorCode, data, message)
}
