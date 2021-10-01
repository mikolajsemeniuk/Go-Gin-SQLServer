package application

import "go-gin-sqlserver/program/controllers"

func Route() {
	// user
	router.GET("/user", controllers.GetUsers)
	router.GET("/user/:userid", controllers.GetUser)
	router.POST("/user", controllers.AddUser)
	router.PATCH("/user/:userid", controllers.UpdateUser)
	router.DELETE("/user/:userid", controllers.RemoveUser)
}
