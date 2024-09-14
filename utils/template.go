package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	FAILED  int = -1
	SUCCESS int = 0
)

func Success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": SUCCESS,
		"msg":  msg,
		"data": data,
	})
}

func Fail(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": FAILED,
		"msg":  msg,
		"data": data,
	})
}
