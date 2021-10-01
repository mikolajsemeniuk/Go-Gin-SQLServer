package services

import "go-gin-sqlserver/program/models"

func GetUsers() ([]models.User, error) {
	var users []models.User
	return users, nil
}
