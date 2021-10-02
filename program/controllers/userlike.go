package controllers

import (
	"go-gin-sqlserver/program/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetUserLike(context *gin.Context) {
	followingId, error := strconv.ParseInt(context.Param("followingid"), 10, 64)
	if error != nil {
		context.JSON(http.StatusBadRequest, "cannot convert followingid to int64")
		return
	}

	followerId, error := strconv.ParseInt(context.Param("followerid"), 10, 64)
	if error != nil {
		context.JSON(http.StatusBadRequest, "cannot convert followerid to int64")
		return
	}

	if error := services.SetUserLike(followingId, followerId); error != nil {
		context.JSON(http.StatusBadRequest, error.Error())
		return
	}

	context.JSON(http.StatusOK, "like set")
}
