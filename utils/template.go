package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(ctx *gin.Context, code int32, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
	})
}

func Fail(ctx *gin.Context, code int32, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": data,
	})
}
