package services

import (
	"context"
	"database/sql"
	"errors"
	"go-gin-sqlserver/program/database"
	"go-gin-sqlserver/program/payloads"
	"time"
)

const (
	GET_FOLLOWERS_QUERY     = "SELECT username FROM [db].[dbo].[users] JOIN [db].[dbo].[user_likes] ON [db].[dbo].[users].[user_id] = [db].[dbo].[user_likes].[follower_id] WHERE [following_id] = @UserId"
	GET_FOLLOWING_QUERY     = "SELECT username FROM [db].[dbo].[users] JOIN [db].[dbo].[user_likes] ON [db].[dbo].[users].[user_id] = [db].[dbo].[user_likes].[following_id] WHERE [follower_id] = @UserId"
	EXISTS_USERLIKE_QUERY   = "SELECT COUNT(*) FROM [db].[dbo].[user_likes] WHERE [following_id] = @FollowingId AND [follower_id] = @FollowerId;"
	INSERT_USERLIKE_COMMAND = "INSERT INTO [db].[dbo].[user_likes] (following_id, follower_id) VALUES (@FollowingId, @FollowerId);"
	REMOVE_USERLIKE_COMMAND = "DELETE FROM [db].[dbo].[user_likes] WHERE following_id = @FollowingId AND follower_id = @FollowerId;"
)

func GetUserLikes(userId int64, query string) ([]payloads.UserLike, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var likes []payloads.UserLike
	rows, error := database.Client.QueryContext(context, query, sql.Named("UserId", userId))
	if error != nil {
		return likes, error
	}
	defer rows.Close()

	for rows.Next() {
		var like payloads.UserLike
		if error := rows.Scan(&like.Username); error != nil {
			return likes, error
		}
		likes = append(likes, like)
	}

	if error = rows.Err(); error != nil {
		return likes, error
	}

	if len(likes) == 0 {
		likes = []payloads.UserLike{}
	}

	return likes, nil
}

func UserLikeExists(followingId int64, followerId int64) (bool, error) {
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exists := 0

	if _, error := GetUser(followingId); error != nil {
		return false, error
	}

	if _, error := GetUser(followerId); error != nil {
		return false, error
	}

	if error := database.Client.
		QueryRowContext(context, EXISTS_USERLIKE_QUERY, sql.Named("FollowingId", followingId), sql.Named("FollowerId", followerId)).
		Scan(&exists); error != nil {
		return false, error
	}

	return exists != 0, nil
}

func SetUserLike(followingId int64, followerId int64) error {
	if followingId == followerId {
		return errors.New("you cannot like yourself")
	}

	exists, error := UserLikeExists(followingId, followerId)
	if error != nil {
		return error
	}

	query := INSERT_USERLIKE_COMMAND
	if exists {
		query = REMOVE_USERLIKE_COMMAND
	}

	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement, error := database.Client.PrepareContext(context, query)
	if error != nil {
		return error
	}
	defer statement.Close()

	result, error := statement.ExecContext(context, sql.Named("FollowingId", followingId), sql.Named("FollowerId", followerId))
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
