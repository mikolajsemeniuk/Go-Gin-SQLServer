package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-gin-sqlserver/program/database"
	"go-gin-sqlserver/program/inputs"
	"go-gin-sqlserver/program/payloads"
	"time"
)

const (
	GET_USERS_QUERY     = "SELECT [user_id], [username] FROM [db].[dbo].[users];"
	GET_USER_QUERY      = "SELECT [user_id], [username] FROM [db].[dbo].[users] WHERE [user_id] = @UserId;"
	ADD_USER_COMMAND    = "INSERT INTO [db].[dbo].[users] ([username]) VALUES (@Username);"
	UPDATE_USER_COMMAND = "UPDATE [db].[dbo].[users] SET [username] = @Username WHERE [user_id] = @UserId;"
	REMOVE_USER_COMMAND = "DELETE FROM [db].[dbo].[users] WHERE [user_id] = @UserId;"
)

func GetUsers() ([]payloads.User, error) {
	var users []payloads.User

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, error := database.Client.QueryContext(context, GET_USERS_QUERY)

	if error != nil {
		return users, error
	}
	defer rows.Close()

	for rows.Next() {
		var user payloads.User
		if error := rows.Scan(&user.UserId, &user.Username); error != nil {
			return users, error
		}

		users = append(users, user)
	}

	if error = rows.Err(); error != nil {
		return users, error
	}

	if len(users) == 0 {
		users = []payloads.User{}
	}

	return users, nil
}

func AddUser(input inputs.User) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement, error := database.Client.PrepareContext(context, ADD_USER_COMMAND)
	if error != nil {
		return error
	}
	defer statement.Close()

	result, error := statement.ExecContext(context, sql.Named("Username", input.Username))
	if error != nil {
		return error
	}

	rows, error := result.RowsAffected()
	if error != nil {
		return error
	}

	if rows == 0 {
		return errors.New("something went wrong")
	}

	return nil
}

func GetUser(userId int64) (payloads.User, error) {
	var user payloads.User

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if error := database.Client.
		QueryRowContext(context, GET_USER_QUERY, sql.Named("UserId", userId)).
		Scan(&user.UserId, &user.Username); error != nil {
		if error == sql.ErrNoRows {
			return payloads.User{}, fmt.Errorf("no user with id of: %d", userId)
		} else {
			return payloads.User{}, error
		}
	}

	return user, nil
}

func UpdateUser(userId int64, input inputs.User) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, error := GetUser(userId); error != nil {
		return error
	}

	statement, error := database.Client.PrepareContext(context, UPDATE_USER_COMMAND)
	if error != nil {
		return error
	}

	result, error := statement.ExecContext(context, sql.Named("Username", input.Username), sql.Named("UserId", userId))
	if error != nil {
		return error
	}

	rows, error := result.RowsAffected()
	if error != nil {
		return error
	}

	if rows == 0 {
		return errors.New("something went wrong")
	}

	return nil
}

func RemoveUser(userId int64) error {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, error := GetUser(userId); error != nil {
		return error
	}

	statement, error := database.Client.PrepareContext(context, REMOVE_USER_COMMAND)
	if error != nil {
		return error
	}

	result, error := statement.ExecContext(context, sql.Named("UserId", userId))
	if error != nil {
		return error
	}

	rows, error := result.RowsAffected()
	if error != nil {
		return error
	}

	if rows == 0 {
		return errors.New("something went wrong")
	}

	return nil
}
