package controllers

import (
	"go-gin-sqlserver/program/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetPostLike(context *gin.Context) {
	userId, error := strconv.ParseInt(context.Param("userid"), 10, 64)
	if error != nil {
		context.JSON(http.StatusBadRequest, "cannot convert userid to int64")
		return
	}

	postId, error := strconv.ParseInt(context.Param("postid"), 10, 64)
	if error != nil {
		context.JSON(http.StatusBadRequest, "cannot convert postid to int64")
		return
	}

	if error := services.SetPostLike(userId, postId); error != nil {
		context.JSON(http.StatusBadRequest, error.Error())
		return
	}

	context.JSON(http.StatusOK, "postlike set")
}
