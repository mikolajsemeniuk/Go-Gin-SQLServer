package controllers

import (
	"go-gin-sqlserver/program/inputs"
	"go-gin-sqlserver/program/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Post struct{}

func AddPost(context *gin.Context) {
	userId, error := strconv.ParseInt(context.Param("userid"), 10, 64)
	if error != nil {
		context.JSON(http.StatusBadRequest, "cannot convert userid to int64")
		return
	}

	var input inputs.Post
	if error := context.ShouldBindJSON(&input); error != nil {
		context.JSON(http.StatusBadRequest, "error binding data")
		return
	}

	if error := services.AddPost(userId, input); error != nil {
		context.JSON(http.StatusBadRequest, error.Error())
		return
	}
	context.JSON(http.StatusOK, "post added")
}

func UpdatePost(context *gin.Context) {
	postId, error := strconv.ParseInt(context.Param("postid"), 10, 64)
	if error != nil {
		context.JSON(http.StatusBadRequest, "cannot convert postid to int64")
		return
	}

	var input inputs.Post
	if error := context.ShouldBindJSON(&input); error != nil {
		context.JSON(http.StatusBadRequest, "error binding data")
		return
	}

	if error := services.UpdatePost(postId, input); error != nil {
		context.JSON(http.StatusBadRequest, error.Error())
		return
	}

	context.JSON(http.StatusOK, "post updated")
}

func RemovePost(context *gin.Context) {
	postId, error := strconv.ParseInt(context.Param("postid"), 10, 64)
	if error != nil {
		context.JSON(http.StatusBadRequest, "cannot convert postid to int64")
		return
	}

	if error := services.RemovePost(postId); error != nil {
		context.JSON(http.StatusBadRequest, error.Error())
		return
	}

	context.JSON(http.StatusOK, "post removed")
}
